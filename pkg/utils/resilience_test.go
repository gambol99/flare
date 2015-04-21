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

package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttempt(t *testing.T) {
	assert.Nil(t, Attempt(3, func() error {
		return nil
	}))
	assert.NotNil(t, Attempt(4, func() error {
		return fmt.Errorf("failed to perform method")
	}))
}

func TestTimed(t *testing.T) {
	result, taken, err := Timed(func(arguments ...interface{}) (interface{}, error) {
		return 1, nil
	})
	assert.Nil(t, err)
	assert.NotNil(t, taken)
	//t.Logf("Time taken measured at: %s", taken)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result)
}
