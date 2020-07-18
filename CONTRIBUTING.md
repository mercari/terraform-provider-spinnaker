# Contributing

Interested in contributing to Spinnaker? Please help us!

## Setup

### Go

[Install Go](https://golang.org/doc/install). 

### Go modules

Clone the repository to a directory outside of your GOPATH:

```bash
$ git clone https://github.com/mercari/terraform-provider-spinnaker
```

Afterward, use `go build` to build the program. This will automatically fetch dependencies.

```bash
$ go build
```

Upon first build, you may see output while the `go` tool fetches dependencies.
To verify dependencies match checksums under go.sum, run `go mod verify`.
To clean up any old, unused go.mod or go.sum lines, run `go mod tidy`.

## Running Spinnaker Provider

Create a `provider.tf`.

```hcl
provider "spinnaker" {
  gate_endpoint = "https://spinnaker-api.xxx.com"
}
```

Build this provider.

```console
$ go build
```

Then run the Terraform operations.

```console
$ terraform init
```

## Running tests

Test the provider by

```bash
go test -v ./...
```
