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

package main

import (
	"flag"

	"github.com/gambol99/flare/pkg/config"

	"fmt"
)

var (
	cfg *config.FlareConfiguration
)

func init() {
	cfg = config.DefaultConfig()
	flag.StringVar(&cfg.DockerSocket, "docker", "/var/run/docker.sock", "the location for the docker socket")
	flag.StringVar(&cfg.LogLevel, "log", "INFO", "the logging level to run as")
	flag.StringVar(&cfg.Bootstrap, "bootstrap", "", "the location of a file used to bootstrap rules into the store")
	flag.StringVar(&cfg.DefaultPolicy, "policy", "DROP", "the default policy for iptables traffic")
	flag.StringVar(&cfg.DockerInterface, "interface", "docker0", "the bridge interface the containers are attached to")
	flag.StringVar(&cfg.StoreURL, "store", "etcd://127.0.0.1:4001", "the url for the store container the rules")
	flag.Var(cfg.StoreOptions, "sopt", "add optional argument to the store (e.g key=value)")
	flag.StringVar(&cfg.DiscoveryURL, "discovery", "", "the service discovery url, e.g. consul://host:8500, etcd://host:4001 etc")
}

func main() {
	flag.Parse()
	fmt.Printf("-- Configuration -- %s", cfg)

}
