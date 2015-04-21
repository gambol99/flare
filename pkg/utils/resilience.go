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

import "time"

func Attempt(attempts int, method func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = method(); err == nil {
			return nil
		}
	}
	return err
}

// Time the duration of a method
//  method:     the method to implement
//  arguments:  a list of arguments to pass to the method
func Timed(method func(...interface{}) (interface{}, error), arguments ...interface{}) (interface{}, time.Duration, error) {
	start_time := time.Now()
	result, err := method(arguments...)
	return result, time.Now().Sub(start_time), err
}

func Try(timeout time.Duration, method func() (interface{}, error)) error {

	return nil
}
