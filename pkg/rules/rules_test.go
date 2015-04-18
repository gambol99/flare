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
	"sync"
	"testing"

	"github.com/gambol99/flare/pkg/api"
	"github.com/gambol99/flare/pkg/config"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	rulesLock  sync.Once
	rules RulesService
)

func createTestRuleService(t *testing.T) RulesService {
	rulesLock.Do(func() {
		var err error
		log.SetLevel(log.PanicLevel)
		cfg := config.DefaultConfig()
		cfg.StoreURL = "etcd://127.0.0.1:4001"
		cfg.Bootstrap = "file://./tests/rules.yaml"
		cfg.DiscoveryURL = "consul://127.0.0.1:8500"
		rules, err = New(cfg)
		if err != nil {
			t.Logf("failed the create a rules service, skipping the tests")
		}
	})
	if rules == nil {
		t.SkipNow()
	}
	return rules
}

func TestNewRuleStore(t *testing.T) {
	rules := createTestRuleService(t)
	assert.NotNil(t, rules)
}

func TestAddGroup(t *testing.T) {
	group := api.NewFlareGroup()
	group.ID = "security_group_authentication"
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(389).Accept().Comment("allow ldap authentication")
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(669).Accept().Comment("allow ldaps authentication")
	group.ID = "security_group_authentication_two"
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(389).Accept().Comment("allow ldap authentication")
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(669).Accept().Comment("allow ldaps authentication")

	err := rules.Add(group)
	assert.Nil(t, err)

	err = rules.Add(group)
	assert.Nil(t, err)
}

func TestIsGroup(t *testing.T) {
	r := createTestRuleService(t)
	group := api.NewFlareGroup()
	group.ID = "security_group_authentication"
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(389).Accept().Comment("allow ldap authentication")

	assert.Nil(t, r.Add(group))
	//assert.True(t, r.IsGroup("security_group_authentication"))
	//assert.False(t, r.IsGroup("security_group_authentication_fake"))
}

func TestGroups(t *testing.T) {
	agent := createTestRuleService(t)
	groups, err := agent.ListGroups()
	assert.Nil(t, err)
	if !assert.NotNil(t, groups) {
		t.FailNow()
	}
	if !assert.NotEmpty(t, groups) {
		t.FailNow()
	}
}

func TestSync(t *testing.T) {
	_ = createTestRuleService(t)
	// delete everything first




}
