// Copyright 2020 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"github.com/RedHatInsights/insights-results-aggregator/types"
)

// GetRuleWithContent returns rule with content for provided ruleID and ruleErrorKey
func (storage DBStorage) GetRuleWithContent(ruleID types.RuleID, ruleErrorKey types.ErrorKey) (*types.RuleWithContent, error) {
	var (
		result             types.RuleWithContent
		impact, likelihood int
		tags               string
	)

	err := storage.connection.QueryRow(`
		SELECT
			rule.module,
			rule.name,
			rule.summary,
			rule.reason,
			rule.resolution,
			rule.more_info,
			rek.error_key,
			rek.condition,
			rek.description,
			rek.impact,
			rek.likelihood,
			rek.publish_date,
			rek.active,
			rek.generic,
			rek.tags
		FROM rule
		INNER JOIN rule_error_key rek ON rule.module = rek.rule_module
		WHERE rule.module = $1 AND rek.error_key = $2
	`, ruleID, ruleErrorKey).Scan(
		&result.Module,
		&result.Name,
		&result.Summary,
		&result.Reason,
		&result.Resolution,
		&result.MoreInfo,
		&result.ErrorKey,
		&result.Condition,
		&result.Description,
		&impact,
		&likelihood,
		&result.PublishDate,
		&result.Active,
		&result.Generic,
		&tags,
	)
	err = types.ConvertDBError(err, []interface{}{ruleID, ruleErrorKey})
	if err != nil {
		return nil, err
	}

	result.TotalRisk = calculateTotalRisk(impact, likelihood)
	result.Tags = commaSeparatedStrToTags(tags)

	return &result, nil
}
