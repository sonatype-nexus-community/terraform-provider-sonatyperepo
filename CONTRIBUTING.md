# Contributing

Pull requests are welcome - we thank you in advance. All we ask is that your PR has a clear description of its purpose or intent and that the purpose of a PR is a singluar intent - we'd rather accept multiple smaller PRs than fewer big ones.

## Setup

This provider uses the Custom Provider Framework from HashiCorp. A great reference is available from HashiCorp [here](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider).

Development is conducted using Go version defined in `go.mod`.

Documentation and examples should be regenerated and included in Pull Requests - run `cd tools; go generate ./...` before finalising your PR.

## Linting

Run `golangci-lint run`.

## Testing

### Unit Tests

These can be run locally by running `go test -v -cover ./internal/provider/`.

## Acceptance/Integration Tests

These require an active and licensed Sonatype IQ Server. PRs originating from outside this project will fail to pass the automated Integration Tests in CI due to our Repository Secrets not being available for these CI Executions (a GitHub restriction).

To run Integration Tests locally, set the following 3 environment variables and then run `TF_ACC=1 go test -v -cover ./internal/provider/`:
- `NXRM_SERVER_URL`: Full URL to your Sonatype Nexus Repository Manager
- `NXRM_SERVER_USERNAME`: Username to authenticate with
- `NXRM_SERVER_PASSWORD`: Password to authentivate with

It is helpful when submitting Pull Requests to confirm whether you have been able to execute the Integraton Tests locally, but not mandatory.

Some Acceptance Tests cannot be run in parallel against a single Sonatype Nexus Repository Manager. These are isolated and only run when `TF_ACC_SINGLE_HIT=1` is also set.

## Standardised Development Patterns

See [terraform-provider-shared](https://github.com/sonatype-nexus-community/terraform-provider-shared) and it's examples.

## Sign off your commits

Please sign off your commits, to show that you agree to publish your changes under the current terms and licenses of the project, and to indicate agreement with [Developer Certificate of Origin (DCO)](https://developercertificate.org/).

```shell
git commit --signed-off ...
```