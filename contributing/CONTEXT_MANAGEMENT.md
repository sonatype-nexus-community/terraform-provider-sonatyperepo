# Centralized API Context Management

## Overview

This document describes the centralized API context management system implemented in the `terraform-provider-sonatyperepo`. This system provides a standardized way to handle authentication context setup for API calls across all resources and data sources.

## Background

Previously, authentication context was set up manually in each resource and data source method using repetitive code:

```go
ctx = context.WithValue(
    ctx,
    sonatyperepo.ContextBasicAuth,
    r.Auth,  // or d.Auth for data sources
)
```

This pattern was repeated ~50+ times throughout the codebase, leading to:
- Code duplication
- Maintenance overhead
- Potential for inconsistent implementation
- Risk of authentication setup errors

## Implementation

### Core Utilities (`internal/provider/common/context.go`)

The system provides utility functions for context management:

```go
// AuthContext represents authentication information
type AuthContext struct {
    Auth sonatyperepo.BasicAuth
}

// NewAuthContext creates an AuthContext from BasicAuth
func NewAuthContext(auth sonatyperepo.BasicAuth) *AuthContext

// WithAuthContext adds authentication to context for API calls
func WithAuthContext(ctx context.Context, authCtx *AuthContext) context.Context

// WithAuth adds authentication directly from BasicAuth
func WithAuth(ctx context.Context, auth sonatyperepo.BasicAuth) context.Context
```

### Base Interfaces

#### Resources (`internal/provider/common/resource.go`)
```go
type BaseResource struct {
    // ... existing fields
    Auth sonatyperepo.BasicAuth
}

// GetAuthContext returns authenticated context for API calls
func (r *BaseResource) GetAuthContext(ctx context.Context) context.Context {
    return WithAuth(ctx, r.Auth)
}
```

#### Data Sources (`internal/provider/common/data_source.go`)
```go
type BaseDataSource struct {
    // ... existing fields
    Auth sonatyperepo.BasicAuth
}

// GetAuthContext returns authenticated context for API calls
func (d *BaseDataSource) GetAuthContext(ctx context.Context) context.Context {
    return WithAuth(ctx, d.Auth)
}
```

## Usage

### Resources

**Before:**
```go
func (r *myResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // ... validation ...
    ctx = context.WithValue(ctx, sonatyperepo.ContextBasicAuth, r.Auth)
    // API call
}
```

**After:**
```go
func (r *myResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // ... validation ...
    ctx = r.GetAuthContext(ctx)
    // API call
}
```

### Data Sources

**Before:**
```go
func (d *myDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    // ... validation ...
    ctx = context.WithValue(ctx, sonatyperepo.ContextBasicAuth, d.Auth)
    // API call
}
```

**After:**
```go
func (d *myDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    // ... validation ...
    ctx = d.GetAuthContext(ctx)
    // API call
}
```