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
	"math/rand"
	"time"
	"fmt"
)

func init() {
	rand.Seed(time.Nanosecond.Nanoseconds())
}

var (
	uuidNumbers = []rune("0123456789")
)

// generic validate interface
type Validation interface {
	IsValid() error
}

// Generate a random number string of x amount
//	min:		the length of the string you want generated
func RandomUUID(min int) string {
	id := make([]rune, min)
	for i := range id {
		id[i] = uuidNumbers[rand.Intn(len(uuidNumbers))]
	}
	return string(id)
}

// Generate a random number
//	min:		the minimum value to return
//	max:		the maximum value to return
func RandomInt(min, max int) (int, error) {
	if min >= max {
		return 0, fmt.Errorf("the minimum value cannot be equal to greater than the max")
	}
	return rand.Intn(max - min) + min, nil
}
