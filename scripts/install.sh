#!/usr/bin/env bash

RELEASE=0.2.0
OS=linux
ARCH=arm

wget "https://github.com/costap/tunnel/releases/download/${RELEASE}/tunneld_${OS}_${ARCH}" -o tunneld
wget "https://github.com/costap/tunnel/releases/download/${RELEASE}/tunnelctl_${OS}_${ARCH}" -o tunnelctl

install -m 0755 tunnelctl /usr/local/bin
install -m 0755 tunneld /usr/local/bin

rm tunnelctl tunneld
