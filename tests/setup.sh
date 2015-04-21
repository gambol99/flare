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
    echo -n "check: $1 "
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
  check "curl -Lk https://github.com/coreos/etcd/releases/download/v2.0.0/etcd-v2.0.0-linux-amd64.tar.gz -o /tmp/etcd.tar.gz"
  check "tar zvxf /tmp/etcd.tar.gz -C /tmp"
  annonce "starting the etcd service"
  check "/tmp/etcd*amd64/etcd > /dev/null 2>&1 &"
  check "sleep 3"
}

perform_setup
