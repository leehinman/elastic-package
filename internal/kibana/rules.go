// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package kibana

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

//SearchResultSignal represents a JSON signal object in a search result
type SearchResultSignal struct {
	Rule RuleNameOnly `json:"rule"`
}

//SearchResultSource represents
type SearchResultSource struct {
	Signal SearchResultSignal `json:"signal"`
}

//SearchResultInnerHit represents a Json _source object in a search result
type SearchResultInnerHit struct {
	Source SearchResultSource `json:"_source"`
}

//SearchResultHitsTotal represents a JSON total object in a search result
type SearchResultHitsTotal struct {
	Value int `json:"value"`
}

//SearchResultHits represents a JSON hits object in search result
type SearchResultHits struct {
	Total SearchResultHitsTotal  `json:"total"`
	Hits  []SearchResultInnerHit `json:"hits"`
}

//SearchResult represent a JSON detection_engine/signals/search result body
type SearchResult struct {
	Hits SearchResultHits `json:"hits"`
}

//RuleIDOnly JSON rep on DetectionRule with only rule_id field
type RuleIDOnly struct {
	RuleID string `json:"rule_id"`
}

//RuleNameOnly JSON rep of Detection Rule with only name field
type RuleNameOnly struct {
	Name string `json:"name"`
}

//RulesBulkCreate Creates new detection rules, expectation is that
//rules are a JSON array of rules where each rule contains the
//required fields
func (c *Client) RulesBulkCreate(rules []byte) ([]RuleIDOnly, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	statusCode, respBody, err := c.post("api/detection_engine/index", nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create detection index")
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("could not create detection index; API status code = %d; response body = %s", statusCode, respBody)
	}

	statusCode, respBody, err = c.post("api/detection_engine/rules/_bulk_create", rules)
	if err != nil {
		return nil, errors.Wrap(err, "could not bulk create rules")
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("could not bulk create rules; API status code = %d; response body = %s", statusCode, respBody)
	}
	results := []RuleIDOnly{}
	err = json.Unmarshal(respBody, &results)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal bulk create response")
	}
	return results, nil
}

//RulesBulkDelete deletes detection rules, expectiation is that rules
//are JSON array of id or rule_id fields of the rules you want to
//delete.
func (c *Client) RulesBulkDelete(rules []RuleIDOnly) error {
	if len(rules) == 0 {
		return nil
	}
	postBody, err := json.Marshal(rules)
	if err != nil {
		return errors.Wrap(err, "could not marshal bulk detection rules")
	}
	statusCode, respBody, err := c.delete("api/detection_engine/rules/_bulk_delete", postBody)
	if err != nil {
		return errors.Wrap(err, "could not bulk delete rules")
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("could not bulk delete rules; API status code = %d; response body = %s", statusCode, respBody)
	}
	return nil
}

//RulesGetHits queries to see if any of the detection rules have been hit
func (c *Client) RulesGetHits() (SearchResult, error) {
	results := SearchResult{}
	postBody := []byte("{\"query\":{\"match\":{\"signal.status\":\"open\"}}}")

	statusCode, respBody, err := c.post("api/detection_engine/signals/search", postBody)
	if err != nil {
		return results, errors.Wrap(err, "could not search for hits")
	}
	if statusCode != http.StatusOK {
		return results, fmt.Errorf("could not search for hits; API status code = %d; response body = %s", statusCode, respBody)
	}

	err = json.Unmarshal(respBody, &results)
	if err != nil {
		return results, errors.Wrap(err, "could not unmarshal search results")
	}
	return results, nil
}

//RulesCloseSignals closes any open signals
func (c *Client) RulesCloseSignals() error {
	postBody := []byte("{\"query\":{\"match\":{\"signal.status\":\"open\"}},\"status\":\"closed\"}")
	statusCode, respBody, err := c.post("api/detection_engine/signals/status", postBody)
	if err != nil {
		return errors.Wrap(err, "could not close signals")
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("could not close signals; API status code = %d; response body = %s", statusCode, respBody)
	}
	return nil
}
