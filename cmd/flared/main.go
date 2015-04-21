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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gambol99/flare/pkg/config"
	"github.com/gambol99/flare/pkg/service"

	log "github.com/Sirupsen/logrus"
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
	flag.StringVar(&cfg.DefaultSecurityGroup, "default_group", "", "the default security group applied regardless")
	flag.StringVar(&cfg.DockerInterface, "interface", "docker0", "the bridge interface the containers are attached to")
	flag.StringVar(&cfg.FlarePrefix, "prefix", "FLARE_", "the prefix used by flare specs in the environment variables")
	flag.StringVar(&cfg.StoreURL, "store", "etcd://127.0.0.1:4001", "the url for the store container the rules")
	flag.Var(cfg.StoreOptions, "so", "add optional argument to the store (e.g key=value)")
	flag.StringVar(&cfg.DiscoveryURL, "discovery", "", "the service discovery url, e.g. consul://host:8500, etcd://host:4001 etc")
	flag.Var(cfg.DiscoveryOptions, "do", "add optional arguments to the discovery provider (KEY=VALUE)")
}

func main() {
	flag.Parse()
	fmt.Println("Starting Flared Service, version: %s, author: %s (%s)",
		service.VERSION, service.AUTHOR, service.EMAIL)
	// step: create the flare service
	if flare, err := service.New(cfg); err != nil {
		log.Fatalf("unable to create the flare service, error: %s", err)
	} else {
		_ = flare
	}

	// step: setup the channel for shutdown signals
	signalChannel := make(chan os.Signal)
	// step: register the signals
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// step: wait on the signal
	<-signalChannel

	log.Info("Exitting the Flare service")
}
