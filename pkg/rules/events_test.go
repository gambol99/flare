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

    "github.com/stretchr/testify/assert"
)


func TestDeleteListener(t *testing.T) {
    rules := createTestRuleService(t)
    ch := make(RulesEventChannel)
    assert.NoError(t, rules.AddListener(ch))
    assert.Error(t, rules.AddListener(nil))
    listeners := rules.Listeners()
    assert.NotNil(t, listeners)
    assert.NotEmpty(t, listeners)
    assert.Equal(t, ch, listeners[0])

    rules.DeleteListener(ch)
    listeners = rules.Listeners()
    assert.NotNil(t, listeners)
    assert.Empty(t, listeners)
}

func TestAddListener(t *testing.T) {
    rules := createTestRuleService(t)
    ch := make(RulesEventChannel)
    assert.NoError(t, rules.AddListener(ch))
    assert.Error(t, rules.AddListener(nil))
    listeners := rules.Listeners()
    assert.NotNil(t, listeners)
    assert.NotEmpty(t, listeners)
    assert.Equal(t, ch, listeners[0])
    rules.DeleteListener(ch)
}
