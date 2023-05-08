package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"reflect"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/config"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
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

// ProcessRequest func - process the unpaginated request
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

	// create the http request, set the user agent, setup the authentication token, log the request
	request, error, diags := CreateRequest(ctx, m, b, diags, resourceType, operation)
	if error {
		return diags
	}

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

// ProcessPagedRequest func - process the paginated request
func ProcessPagedRequest(ctx context.Context, m interface{}, vars interface{}, data interface{}, query string, resourceType string, operation string, maxPages int) (diags diag.Diagnostics, allthedata []interface{}) {

	tflog.Info(ctx, "client.ProcessPagedRequest called...")
	tflog.Debug(ctx, fmt.Sprintf("Received vars: %T, %s", vars, utils.PrettyPrint(vars)))
	tflog.Debug(ctx, fmt.Sprintf("Received query: %T, %s", query, query))
	tflog.Debug(ctx, fmt.Sprintf("Received resourceType/operation: %s %s", resourceType, operation))
	tflog.Debug(ctx, fmt.Sprintf("Received maxPages: %d", maxPages))

	// get an http client
	client := m.(*config.ProviderConf).HTTPClient

	// encode the request body (graphql query and variables)
	b := new(bytes.Buffer)
	switch op := operation; op {
	case "read":
		err := json.NewEncoder(b).Encode(GraphQLRequest{Query: query, Variables: vars})
		if err != nil {
			return append(diags, diag.FromErr(err)...), nil
		}
		tflog.Debug(ctx, fmt.Sprintf("%s %s request variables: %s", resourceType, operation, utils.PrettyPrint(vars)))
	default:
		return append(diags, diag.FromErr(fmt.Errorf("operation %s not supported for paged operations", operation))...), nil
	}

	// create the http request, set the user agent, setup the authentication token, log the request
	request, error, diags := CreateRequest(ctx, m, b, diags, resourceType, operation)
	if error {
		return diags, nil
	}

	var allData []interface{}
	endCursor := ""
	paginate := true
	currentPage := 0
	// loop through the pages, while there are more pages to process
	// maxPages of 0 fetches all pages, there is an OR grouping for the third sub-condition
	for paginate && maxPages >= 0 && (currentPage < maxPages || maxPages == 0) {
		currentPage++
		tflog.Debug(ctx, fmt.Sprintf("Processing page %d with a maximum of %d pages (maximum of 0 means unlimited)", currentPage, maxPages))
		// make the request using `endCursor` if it's not empty
		if endCursor != "" {
			queryVars, ok := vars.(*internal.QueryVariables)
			if !ok {
				return append(diags, diag.FromErr(fmt.Errorf("unable to cast vars to internal.QueryVariables"))...), nil
			}
			queryVars.After = endCursor
			vars = interface{}(queryVars)
			err := json.NewEncoder(b).Encode(GraphQLRequest{Query: query, Variables: vars})
			if err != nil {
				return append(diags, diag.FromErr(err)...), nil
			}

			// create the http request, set the user agent, setup the authentication token, log the request
			request, error, diags := CreateRequest(ctx, m, b, diags, resourceType, operation)
			if error {
				return diags, nil
			}
			// make the request and handle the response
			error, diags, continuePaging, newEndCursor := RequestDo(ctx, client, request, diags, resourceType, operation, data, &allData)
			if error {
				return diags, nil
			}

			// update `endCursor` and `paginate`
			endCursor = newEndCursor
			paginate = continuePaging
		} else {
			// make the initial request without `endCursor`
			error, diags, continuePaging, newEndCursor := RequestDo(ctx, client, request, diags, resourceType, operation, data, &allData)
			if error {
				return diags, nil
			}
			// update `endCursor` and `paginate`
			endCursor = newEndCursor
			paginate = continuePaging
		}

		if !paginate {
			break // exit loop if there are no more pages to fetch
		}
	}

	return diags, allData
}

// CreateRequest func - create the http request
func CreateRequest(ctx context.Context, m interface{}, b *bytes.Buffer, diags diag.Diagnostics, resourceType string, operation string) (*http.Request, bool, diag.Diagnostics) {
	request, err := http.NewRequest("POST", m.(*config.ProviderConf).Settings.WizURL, b)
	if err != nil {
		return nil, true, append(diags, diag.FromErr(err)...)
	}

	request.Header.Set("User-Agent", m.(*config.ProviderConf).UserAgent)

	authToken := fmt.Sprintf("%s %s", m.(*config.ProviderConf).TokenType, m.(*config.ProviderConf).Token)
	request.Header.Add("Authorization", authToken)
	request.Header.Add("Content-Type", "application/json")

	reqDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, true, append(diags, diag.FromErr(err)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("%s %s request: %s", resourceType, operation, reqDump))
	return request, false, nil
}

// RequestDo func - make the http request and handle the response
func RequestDo(ctx context.Context, client *http.Client, request *http.Request, diags diag.Diagnostics, resourceType string, operation string, data interface{}, alldata *[]interface{}) (error bool, diagnostics diag.Diagnostics, haspages bool, cursor string) {

	// call the api
	resp, err := client.Do(request)
	if err != nil {
		return true, append(diags, diag.FromErr(err)...), false, ""
	}
	defer resp.Body.Close()

	// log the response
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return true, append(diags, diag.FromErr(err)...), false, ""
	}
	tflog.Debug(ctx, fmt.Sprintf("%s %s api response: %s", resourceType, operation, respDump))

	// handle http errors
	if resp.StatusCode != http.StatusOK {
		return true, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("HTTP Response (%d)", resp.StatusCode),
			Detail:   fmt.Sprintf("Response: %s", respDump),
		}), false, ""
	}

	// read the response
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return true, append(diags, diag.FromErr(err)...), false, ""
	}

	// copy the value of data to a new instance
	newData := reflect.New(reflect.TypeOf(data).Elem()).Interface()
	reflect.ValueOf(newData).Elem().Set(reflect.ValueOf(data).Elem())

	// unmarshal the response to the new instance
	responseBody := &MutationPayload{Data: newData}
	err = json.Unmarshal(rbody, &responseBody)
	if err != nil {
		return true, append(diags, diag.FromErr(err)...), false, ""
	}

	// handle errors from the api
	errorCount := len(responseBody.Errors)
	tflog.Debug(ctx, fmt.Sprintf("Error count: %d", errorCount))
	if errorCount > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Errors returned from API (%d)", errorCount))
		return true, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s %s reported errors", resourceType, operation),
			Detail:   fmt.Sprintf("Response: %s", utils.PrettyPrint(responseBody.Errors)),
		}), false, ""
	}

	// append the page of data to the Data slice, and set the data field in the response body to nil to avoid duplication
	paginationDetails, err := ExtractPageInfo(responseBody.Data)
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error extracting pagination details: %s", err))
		return true, append(diags, diag.FromErr(err)...), false, ""
	}

	tflog.Debug(ctx, fmt.Sprintf("Pagination details: %s", utils.PrettyPrint(paginationDetails)))

	// set the data field in the response body to nil to avoid duplication
	responseBody.Data = nil

	// append the new instance to the allData slice
	*alldata = append(*alldata, newData)

	// log the return data
	tflog.Debug(ctx, fmt.Sprintf("Wrote paginated data: %T, %s", data, utils.PrettyPrint(data)))

	if paginationDetails.HasNextPage {
		return false, nil, true, paginationDetails.EndCursor
	}

	return false, nil, false, ""
}

// ExtractPageInfo func - extract the PageInfo struct from a generic response
func ExtractPageInfo(data interface{}) (wiz.PageInfo, error) {
	var pageInfo wiz.PageInfo

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return pageInfo, fmt.Errorf("data %v is not a struct", value.Kind())
	}
	// loop through the fields of the struct
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := field.Type()

		if fieldType.Kind() == reflect.Struct {
			// loop through the fields of the struct
			for j := 0; j < fieldType.NumField(); j++ {
				subfield := fieldType.Field(j)
				subfieldType := subfield.Type
				// search for the PageInfo struct
				if subfieldType.Kind() == reflect.Struct && subfieldType.Name() == "PageInfo" {
					pageInfoValue := field.Field(j)
					pageInfoType := pageInfoValue.Type()

					for k := 0; k < pageInfoType.NumField(); k++ {
						pageInfoField := pageInfoValue.Field(k)
						switch pageInfoType.Field(k).Name {
						case "EndCursor":
							pageInfo.EndCursor = pageInfoField.String()
						case "HasNextPage":
							pageInfo.HasNextPage = pageInfoField.Bool()

						}
					}
				}
			}
		}
	}
	return pageInfo, nil
}
