/*
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"context"
	"net/http"

	sonatyperepoV382 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
)

// RoutingRuleService abstracts routing rule CRUD across NXRM API client generations.
// Every method's request/response shape is expressed in terms of the V382 generated types,
// which the existing model-mapping layer already consumes.
type RoutingRuleService interface {
	CreateRoutingRule(ctx context.Context, body sonatyperepoV382.RoutingRuleXO) (*http.Response, error)
	GetRoutingRule(ctx context.Context, name string) (*sonatyperepoV382.RoutingRuleXO, *http.Response, error)
	UpdateRoutingRule(ctx context.Context, name string, body sonatyperepoV382.RoutingRuleXO) (*http.Response, error)
	DeleteRoutingRule(ctx context.Context, name string) (*http.Response, error)
	GetRoutingRules(ctx context.Context) ([]sonatyperepoV382.RoutingRuleXO, *http.Response, error)
}

// routingRuleServiceV382 implements RoutingRuleService against NXRM API client V382 (targets NXRM < 3.94.0).
type routingRuleServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *routingRuleServiceV382) CreateRoutingRule(ctx context.Context, body sonatyperepoV382.RoutingRuleXO) (*http.Response, error) {
	return s.client.RoutingRulesAPI.CreateRoutingRule(ctx).Body(body).Execute()
}

func (s *routingRuleServiceV382) GetRoutingRule(ctx context.Context, name string) (*sonatyperepoV382.RoutingRuleXO, *http.Response, error) {
	return s.client.RoutingRulesAPI.GetRoutingRule(ctx, name).Execute()
}

func (s *routingRuleServiceV382) UpdateRoutingRule(ctx context.Context, name string, body sonatyperepoV382.RoutingRuleXO) (*http.Response, error) {
	return s.client.RoutingRulesAPI.UpdateRoutingRule(ctx, name).Body(body).Execute()
}

func (s *routingRuleServiceV382) DeleteRoutingRule(ctx context.Context, name string) (*http.Response, error) {
	return s.client.RoutingRulesAPI.DeleteRoutingRule(ctx, name).Execute()
}

func (s *routingRuleServiceV382) GetRoutingRules(ctx context.Context) ([]sonatyperepoV382.RoutingRuleXO, *http.Response, error) {
	return s.client.RoutingRulesAPI.GetRoutingRules(ctx).Execute()
}

// routingRuleServiceV395 implements RoutingRuleService against NXRM API client V395 (targets NXRM 3.94.0+).
// V395 uses pluralized method names and different builder method for body, requiring careful bridging.
type routingRuleServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *routingRuleServiceV395) CreateRoutingRule(ctx context.Context, body sonatyperepoV382.RoutingRuleXO) (*http.Response, error) {
	var v395Body sonatyperepoV395.RoutingRuleXO
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RoutingRulesAPI.CreateRoutingRules(ctx).RoutingRuleXO(v395Body).Execute()
}

func (s *routingRuleServiceV395) GetRoutingRule(ctx context.Context, name string) (*sonatyperepoV382.RoutingRuleXO, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RoutingRulesAPI.GetRoutingRules(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.RoutingRuleXO
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *routingRuleServiceV395) UpdateRoutingRule(ctx context.Context, name string, body sonatyperepoV382.RoutingRuleXO) (*http.Response, error) {
	var v395Body sonatyperepoV395.RoutingRuleXO
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RoutingRulesAPI.UpdateRoutingRules(ctx, name).RoutingRuleXO(v395Body).Execute()
}

func (s *routingRuleServiceV395) DeleteRoutingRule(ctx context.Context, name string) (*http.Response, error) {
	return s.client.RoutingRulesAPI.DeleteRoutingRules(ctx, name).Execute()
}

func (s *routingRuleServiceV395) GetRoutingRules(ctx context.Context) ([]sonatyperepoV382.RoutingRuleXO, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RoutingRulesAPI.ListRoutingRules(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.RoutingRuleXO
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}
