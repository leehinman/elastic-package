// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package rules

//DetectionRule represents a JSON detection rule
type DetectionRule struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	RiskScore   int    `json:"risk_score"`
	Severity    string `json:"severity"`
	Type        string `json:"type"`
	Query       string `json:"query"`
	Interval    string `json:"interval"`
	From        string `json:"from"`
	RuleID      string `json:"rule_id"`
}
