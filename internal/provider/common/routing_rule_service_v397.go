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
	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// routingRuleServiceV397 implements RoutingRuleService against NXRM API client V397 (targets NXRM 3.97.0+).
// V397 uses pluralized method names and different builder method for body, requiring careful bridging.
type routingRuleServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func (s *routingRuleServiceV397) CreateRoutingRule(ctx context.Context, body sonatyperepoV382.RoutingRuleXO) (*http.Response, error) {
	var v397Body sonatyperepoV397.RoutingRuleXO
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RoutingRulesAPI.CreateRoutingRules(ctx).RoutingRuleXO(v397Body).Execute()
}

func (s *routingRuleServiceV397) GetRoutingRule(ctx context.Context, name string) (*sonatyperepoV382.RoutingRuleXO, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RoutingRulesAPI.GetRoutingRules(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.RoutingRuleXO
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *routingRuleServiceV397) UpdateRoutingRule(ctx context.Context, name string, body sonatyperepoV382.RoutingRuleXO) (*http.Response, error) {
	var v397Body sonatyperepoV397.RoutingRuleXO
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RoutingRulesAPI.UpdateRoutingRules(ctx, name).RoutingRuleXO(v397Body).Execute()
}

func (s *routingRuleServiceV397) DeleteRoutingRule(ctx context.Context, name string) (*http.Response, error) {
	return s.client.RoutingRulesAPI.DeleteRoutingRules(ctx, name).Execute()
}

func (s *routingRuleServiceV397) GetRoutingRules(ctx context.Context) ([]sonatyperepoV382.RoutingRuleXO, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RoutingRulesAPI.ListRoutingRules(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.RoutingRuleXO
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}
