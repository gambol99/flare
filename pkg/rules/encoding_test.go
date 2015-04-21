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

package rules

import (


    "testing"
    "github.com/gambol99/flare/pkg/api"
    "github.com/coreos/etcd/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

const (
    TEST_GROUP = `
    {
  "id": "security_group_authentication",
  "desc": "",
  "active": true,
  "last_updated": "2015-04-19T13:56:00.230670926+01:00",
  "rules": [
    {
      "type": "address",
      "description": "allow ldap authentication",
      "destination": "10.0.100.32",
      "protocol": "tcp",
      "inverse": false,
      "ports": [
        389
      ],
      "action": 0
    },
    {
      "type": "address",
      "description": "allow ldaps authentication",
      "destination": "10.0.100.32",
      "protocol": "tcp",
      "inverse": false,
      "ports": [
        669
      ],
      "action": 0
    },
    {
      "type": "address",
      "description": "allow ldap authentication",
      "destination": "10.0.100.32",
      "protocol": "tcp",
      "inverse": false,
      "ports": [
        389
      ],
      "action": 0
    },
    {
      "type": "address",
      "description": "allow ldaps authentication",
      "destination": "10.0.100.32",
      "protocol": "tcp",
      "inverse": false,
      "ports": [
        669
      ],
      "action": 0
    }
  ]
}`
)

func TestDecode(t *testing.T) {
    agent := createTestRuleService(t).(*RulesStore)
    group := new(api.FlareRuleGroup)
    assert.NoError(t, agent.decode(TEST_GROUP, group))
    assert.Equal(t, "security_group_authentication", group.ID)
    assert.True(t, group.Active)
    if assert.NotNil(t, group.Rules) {
       t.SkipNow()
    }
    assert.Equal(t, 4, len(group.Rules))
    rule := group.Rules[0]
    assert.Equal(t, rule.RuleType, api.TYPE_ADDRESS)
    assert.Equal(t, rule.Description, "allow ldap authentication")
    assert.Equal(t, rule.Action, api.ACTION_ACCEPT)
    assert.NotEmpty(t, rule.Ports)
}

func TestEncode(t *testing.T) {
    agent := createTestRuleService(t).(*RulesStore)
    group := new(api.FlareRuleGroup)
    assert.NoError(t, agent.decode(TEST_GROUP, group))
    content, err := agent.encode(group)
    assert.NoError(t, err)
    if !assert.NotEmpty(t, content) {
        t.FailNow()
    }
}