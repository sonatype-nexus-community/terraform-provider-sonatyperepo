# Centralized Error Handling

## Overview

This document describes the centralized error handling system implemented in the `terraform-provider-sonatyperepo`. The provider uses the shared library's error handling functions for consistent API error handling across all resources and data sources.

## Background and Migration

Previously, API error handling was implemented with wrapper functions in `internal/provider/common/api.go`, which added unnecessary abstraction layers. The provider has been refactored to use the `terraform-provider-shared` v0.2.0 library directly, eliminating the wrapper layer and improving consistency across all Sonatype providers.

See the codebase refactoring for details on how direct integration with the shared library improves maintainability and eliminates code duplication.

## Implementation

### Shared Library Error Handling Functions

The provider now uses error handling functions from `terraform-provider-shared/errors`:

```go
import (
    sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
)

// HandleAPIError provides consistent API error handling
func HandleAPIError(message string, err *error, httpResponse *http.Response, respDiags *diag.Diagnostics)

// HandleAPIWarning provides consistent API warning handling for non-fatal errors
func HandleAPIWarning(message string, err *error, httpResponse *http.Response, respDiags *diag.Diagnostics)
```

### Key Features

1. **Network Error Detection**: Automatically detects and categorizes network-related errors (timeouts, DNS failures, connection refused)

2. **Response Body Handling**: Automatically reads and includes HTTP response bodies in error messages

3. **Consistent Formatting**: Standardized error message format across all resources and providers

4. **HTTP Status Awareness**: Different handling for different HTTP status codes (404s use warnings, others use errors)

## Usage

### API Error Handling

All error handling now uses the shared library directly:

```go
import (
    sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
)

func (r *myResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    httpResp, err := client.Create()
    if err != nil {
        sharederr.HandleAPIError(
            "Error creating repository",
            &err,
            httpResp,
            &resp.Diagnostics,
        )
        return
    }
}
```

### 404 Warning Handling

For resource read operations, 404 responses should be handled as warnings and the resource removed from state:

```go
httpResp, err := client.Read()
if httpResp.StatusCode == 404 {
    resp.State.RemoveResource(ctx)
    sharederr.HandleAPIWarning(
        "Resource to read did not exist",
        &err,
        httpResp,
        &resp.Diagnostics,
    )
    return
}
if err != nil {
    sharederr.HandleAPIError("Error reading resource", &err, httpResp, &resp.Diagnostics)
    return
}
```

## Shared Library Features

The shared library provides additional utilities beyond error handling:

- **Standard Error Functions**: `APIError()`, `NotFoundError()`, `ValidationError()`, `ConflictError()`, `UnauthorizedError()`, `TimeoutError()`
- **Diagnostic Helpers**: `AddAPIErrorDiagnostic()`, `AddNotFoundDiagnostic()`, etc.
- **HTTP Status Helpers**: `IsNotFound()`, `IsForbidden()`, `IsUnauthorized()`, etc.
- **Type Conversions**: `StringToPtr()`, `Int64PtrToValue()`, `SafeString()`, `SafeInt32()`, `SafeBool()`, etc.