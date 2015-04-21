#!/bin/bash

NAME="flare"
ROOT_KEY="/env/testing"

annonce() {
  [ -n "$1" ] && echo "** $@"
}

failed() {
  annonce "[failed] $@"
  exit 1
}

check() {
  if [ -n "$1" ]; then
    echo -n "check: $2 "
    eval "$1 >/dev/null"
    if [ $? -ne 0 ]; then
      echo "[failed]"
      exit 1
    fi
    echo "[passed]"
  fi
}
perform_setup() {
  annonce "downloading the etcd service for tests"
  check "curl -skL https://github.com/coreos/etcd/releases/download/v2.0.0/etcd-v2.0.0-linux-amd64.tar.gz > /tmp/etcd-v2.0.0-linux-amd64.tar.gz "
  check "tar zxf /tmp/etcd-v2.0.0-linux-amd64.tar.gz"
  annonce "starting the etcd service"
  check "nohup /tmp/etcd-v2.0.0-linux-amd64/etcd > etcd.log 2>&1 &"
  check "sleep 3"
}

perform_setup
