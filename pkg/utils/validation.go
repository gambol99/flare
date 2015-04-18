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
	"net"
	"net/url"
	"os"
)

// Checks to see of the address is valid
//  address:    the ip address you want to check
func IsValidIpAddress(address string) bool {
	ip_address := net.ParseIP(address)
	if ip_address.To4() == nil && ip_address.To16() == nil {
		return false
	}
	return true
}

// Retrieve a list of all the interfaces on the system

func GetInterfaces() ([]string, error) {
	list := make([]string, 0)
	if interfaces, err := net.Interfaces(); err != nil {
		return list, err
	} else {
		for _, iface := range interfaces {
			list = append(list, iface.Name)
		}
	}
	return list, nil
}

// Check the value is a valid port
func IsPort(port int) bool {
	if port < 0 || port >= 65536 {
		return false
	}
	return true
}

// Checks if a string is empty
//  str:        the string we are checking
func IsEmpty(str string) bool {
	if str == "" {
		return true
	}
	return false
}

// Checks to see if the path points to a unix socket
//  path:       the location of the socket we are checking
func IsUnixSocket(path string) (bool, error) {
	return true, nil
}

// Checks to see if the interface exists
func IsNetworkInterface(iface string) (bool, error) {
	interfaces, err := GetInterfaces()
	if err != nil {
		return false, err
	}
	for _, name := range interfaces {
		if name == iface {
			return true, nil
		}
	}
	return false, nil
}

// Checks to see if the argument is in the allowed list
//  argument:       the value you are checking
//  allowed:        a list of strings which are allowed
func IsValidArgument(argument interface{}, allowed ...interface{}) bool {
	for _, allow := range allowed {
		if argument == allow {
			return true
		}
	}
	return false
}

// Checks to see if the file exists
//  filename:       the path to the file
func FileExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Checks to see if the url location is valid
//  location:      the url you wish to check
func IsValidURL(location string) bool {
	if _, err := url.Parse(location); err != nil {
		return false
	}
	return true
}
