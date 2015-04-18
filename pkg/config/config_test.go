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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	config := DefaultConfig()
	assert.Nil(t, config.IsValid())

	config.FlarePrefix = ""
	assert.Error(t, config.IsValid())
}

func TestNewConfigurationOptions(t *testing.T) {
	options := NewConfigurationOptions()
	assert.NotNil(t, options)
}

func TestConfigurationOptionGet(t *testing.T) {
	options := NewConfigurationOptions()
	options.Add("test1", 12)
	assert.Equal(t, 12, options.Get("test1", nil))
	assert.Equal(t, 12, options.Get("not_there", 12))
}

func TestConfigurationOptionsSet(t *testing.T) {
	options := NewConfigurationOptions()
	err := options.Set("test1=1")
	assert.Nil(t, err)
	assert.Equal(t, "1", options.Get("test1", nil))
	err = options.Set("cacert=/var/log/cert.crt")
	assert.Nil(t, err)
	assert.Equal(t, "/var/log/cert.crt", options.Get("cacert", nil))
	err = options.Set("cacert")
	assert.NotNil(t, err)
}
