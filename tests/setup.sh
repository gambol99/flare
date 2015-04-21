#!/bin/bash
#
#   Author: Rohith
#   Date: 2015-04-21 15:24:02 +0100 (Tue, 21 Apr 2015)
#
#  vim:ts=2:sw=2:et
#

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
    eval "$1 >/dev/null" || failed "check"
    echo "[passed]"
  fi
}

perform_setup() {
  annonce "downloading the etcd service for tests"
  check "test -f /tmp/etcd.tar.gz || curl -Lk https://github.com/coreos/etcd/releases/download/v2.0.0/etcd-v2.0.0-linux-amd64.tar.gz -o /tmp/etcd.tar.gz"
  check "test -d /tmp/etcd-v2.0.0-linux-amd64 || tar zvxf /tmp/etcd.tar.gz -C /tmp"
  annonce "starting the etcd service"
  check "pidof etcd >/dev/null || /tmp/etcd-v2.0.0-linux-amd64/etcd > /dev/null 2>&1 &"
  check "sleep 5"
}

perform_setup

