/*
Copyright 2014 Rohith All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRule(t *testing.T) {
	rule := NewRule()
	assert.NotNil(t, rule)
}

func TestGroupRuleIsValid(t *testing.T) {
	rule := NewRule()
	rule.RuleType = TYPE_GROUP
	rule.ID = "security_group_ldap"
	assert.NoError(t, rule.IsValid())
	rule.ID = ""
	assert.Error(t, rule.IsValid())
}

func TestRuleIsValid(t *testing.T) {
	rule := NewRule()
	rule.RuleType = TYPE_ADDRESS
	assert.Error(t, rule.IsValid())
	rule.Description = "test rule"
	assert.Error(t, rule.IsValid())
	rule.Port(10).TCP()
	assert.Error(t, rule.IsValid())
	rule.Destination = "90.9.9..00"
	assert.Error(t, rule.IsValid())
	rule.Destination = "10.0.0.1"
	assert.NoError(t, rule.IsValid())
}

func TestRuleCompare(t *testing.T) {
	ruleA := new(FlareRule).Address().Comment("a rule").TCP().Port(22).Addr("10.0.0.32/32")
	ruleB := new(FlareRule).Address().Comment("a rule").TCP().Port(22).Addr("10.0.0.32/32")
	assert.True(t, ruleA.Compare(ruleB))
	ruleA.UDP()
	assert.False(t, ruleA.Compare(ruleB))
	ruleA.TCP()
	ruleB.Port(80)
	assert.False(t, ruleA.Compare(ruleB))
}
