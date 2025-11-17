# Centralized Error Handling

## Overview

This document describes the centralized error handling system implemented in the `terraform-provider-sonatyperepo`. This system provides consistent API error handling across all resources and data sources.

## Background

Previously, API error handling was implemented inline throughout the codebase with repetitive patterns:

```go
if err != nil {
    errorBody, _ := io.ReadAll(httpResponse.Body)
    resp.Diagnostics.AddError(
        "Error creating Resource",
        "Could not create Resource, unexpected error: " + httpResponse.Status + ": " + string(errorBody),
    )
    return
}
```

This led to:
- Code duplication across ~100+ locations
- Inconsistent error message formats
- Manual response body reading
- Mixed error handling patterns

## Implementation

### Core Error Handling Functions (`internal/provider/common/api.go`)

```go
// HandleApiError provides consistent API error handling
func HandleApiError(message string, err *error, httpResponse *http.Response, respDiags *diag.Diagnostics)

// HandleApiWarning provides consistent API warning handling for non-fatal errors
func HandleApiWarning(message string, err *error, httpResponse *http.Response, respDiags *diag.Diagnostics)
```

### Key Features

1. **Network Error Detection**: Automatically detects and categorizes network-related errors (timeouts, DNS failures, connection refused)

2. **Response Body Handling**: Automatically reads and includes HTTP response bodies in error messages

3. **Consistent Formatting**: Standardized error message format across all resources

4. **HTTP Status Awareness**: Different handling for different HTTP status codes (404s use warnings, others use errors)

## Usage

### API Error Handling

**Before:**
```go
if err != nil {
    errorBody, _ := io.ReadAll(httpResponse.Body)
    resp.Diagnostics.AddError(
        "Error creating repository",
        fmt.Sprintf("Error creating repository: %d: %s", httpResponse.StatusCode, string(errorBody)),
    )
    return
}
```

**After:**
```go
if err != nil {
    common.HandleApiError(
        "Error creating repository",
        &err,
        httpResponse,
        &resp.Diagnostics,
    )
    return
}
```

### 404 Warning Handling

**Before:**
```go
if httpResponse.StatusCode == 404 {
    resp.State.RemoveResource(ctx)
    resp.Diagnostics.AddWarning("Resource not found", "...")
}
```

**After:**
```go
if httpResponse.StatusCode == 404 {
    resp.State.RemoveResource(ctx)
    common.HandleApiWarning(
        "Resource to read did not exist",
        &err,
        httpResponse,
        &resp.Diagnostics,
    )
}
```