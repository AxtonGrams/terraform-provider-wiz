package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/config"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
)

// GraphQLRequest struct
type GraphQLRequest struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

// MutationPayload struct
type MutationPayload struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message    string `json:"message,omitempty"`
		Extensions struct {
			Code      string `json:"code,omitempty"`
			Exception struct {
				Message string   `json:"message,omitempty"`
				Path    []string `json:"path,omitempty"`
			} `json:"exception,omitempty"`
		} `json:"extensions,omitempty"`
	} `json:"errors,omitempty"`
}

// MutationInput struct
type MutationInput struct {
	Input interface{} `json:"input"`
}

// ProcessRequest func
func ProcessRequest(ctx context.Context, m interface{}, vars, data interface{}, query, resourceType, operation string) (diags diag.Diagnostics) {
	tflog.Info(ctx, "client.ProcessRequest called...")
	tflog.Debug(ctx, fmt.Sprintf("Received vars: %T, %s", vars, utils.PrettyPrint(vars)))
	tflog.Debug(ctx, fmt.Sprintf("Received query: %T, %s", query, query))
	tflog.Debug(ctx, fmt.Sprintf("Received resourceType/operation: %s %s", resourceType, operation))

	// get an http client
	client := m.(*config.ProviderConf).HTTPClient

	// encode the request body (graphql query and variables)
	b := new(bytes.Buffer)
	switch op := operation; op {
	case "read":
		err := json.NewEncoder(b).Encode(GraphQLRequest{Query: query, Variables: vars})
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		tflog.Debug(ctx, fmt.Sprintf("%s %s request variables: %s", resourceType, operation, utils.PrettyPrint(vars)))
	default:
		input := &MutationInput{}
		input.Input = vars
		err := json.NewEncoder(b).Encode(GraphQLRequest{Query: query, Variables: input})
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		tflog.Debug(ctx, fmt.Sprintf("%s %s request variables: %s", resourceType, operation, utils.PrettyPrint(input)))
	}

	// create the http request
	request, err := http.NewRequest("POST", m.(*config.ProviderConf).Settings.WizURL, b)

	// set the user agent
	request.Header.Set("User-Agent", m.(*config.ProviderConf).UserAgent)

	// setup the authentication token
	authToken := fmt.Sprintf("%s %s", m.(*config.ProviderConf).TokenType, m.(*config.ProviderConf).Token)
	request.Header.Add("Authorization", authToken)
	request.Header.Add("Content-Type", "application/json")

	// log the request
	reqDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("%s %s request: %s", resourceType, operation, reqDump))

	// call the api
	resp, err := client.Do(request)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	defer resp.Body.Close()

	// log the response
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("%s %s api response: %s", resourceType, operation, respDump))

	// handle http errors
	if resp.StatusCode != http.StatusOK {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("HTTP Response (%d)", resp.StatusCode),
			Detail:   fmt.Sprintf("Response: %s", respDump),
		})
	}

	// read the response
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// unmarshal the response
	responseBody := &MutationPayload{Data: data}
	err = json.Unmarshal(rbody, &responseBody)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// handle errors from the api
	errorCount := len(responseBody.Errors)
	tflog.Debug(ctx, fmt.Sprintf("Error count: %d", errorCount))
	if errorCount > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Errors returned from API (%d)", errorCount))
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s %s reported errors", resourceType, operation),
			Detail:   fmt.Sprintf("Response: %s", utils.PrettyPrint(responseBody.Errors)),
		})
	}

	// log the return data
	tflog.Debug(ctx, fmt.Sprintf("Wrote data: %T, %s", data, utils.PrettyPrint(data)))

	return diags
}
