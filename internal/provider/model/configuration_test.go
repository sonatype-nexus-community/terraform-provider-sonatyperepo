package model_test

import (
      "terraform-provider-sonatyperepo/internal/provider/model"
      "testing"

      "github.com/hashicorp/terraform-plugin-framework/types"
      sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)
func TestIqConnectionModelMapToApi_NexusTrustStoreEnabled(t *testing.T) {
      m := model.IqConnectionModel{
              Enabled:                types.BoolValue(true),
              Url:                    types.StringValue("https://nexus-iq.example.com"),
              AuthenticationMethod:   types.StringValue("USER"),
              Username:               types.StringValue("admin"),
              Password:               types.StringValue("password"),
              NexusTrustStoreEnabled: types.BoolValue(true),
              ConnectionTimeout:      types.Int32Value(300),
              Properties:             types.StringValue(""),
              ShowIQServerLink:       types.BoolValue(false),
              FailOpenModeEnabled:    types.BoolValue(false),
      }

      api := sonatyperepo.NewIqConnectionXoWithDefaults()
      m.MapToApi(api)

      if api.UseTrustStoreForUrl == nil {
              t.Fatal("Expected UseTrustStoreForUrl to be set, got nil")
      }
      if *api.UseTrustStoreForUrl != true {
              t.Errorf("Expected UseTrustStoreForUrl=true, got %v", *api.UseTrustStoreForUrl)
      }
}
