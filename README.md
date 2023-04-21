# Terraform Providerfor Wiz

The Terraform provider for Wiz allows you to manage resources typically managed in the Wiz web interface.

This provider is not yet feature complete and requires development, testing, and polishing.

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) >= 1.0
* [Go](https://golang.org/doc/install) >= 1.18

## Building the Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Getting Started

Binaries are available for tagged releases in this repository.

Once you have the provider installed, follow the instructions in the docs folder to understand what options are available.  The documentation includes examples.

## Using the Provider

See the [provider docs](https://registry.terraform.io/providers/AxtonGrams/wiz/latest/docs)

## Contributing

We welcome your contribution. Please understand that the experimental nature of this repository means that contributing code may be a bit of a moving target. If you have an idea for an enhancement or bug fix, and want to take on the work yourself, please first create an issue so that we can discuss the implementation with you before you proceed with the work.

You can review our [contribution guide](_about/CONTRIBUTING.md) to begin. You can also check out our frequently asked questions.
