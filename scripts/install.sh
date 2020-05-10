#!/usr/bin/env bash

VERSION=0.2.1
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH_U=$(uname -m)
if [ "$ARCH_U" == 'x86_64' ];
then
  ARCH="amd64"
fi

if [ "$ARCH_U" == 'armv*' ];
then
  ARCH="arm"
fi

if [ "$ARCH_U" == 'x86_32' ];
then
  ARCH="386"
fi

if [ -z "$ARCH" ]; then echo "$ARCH_U is not supported"; exit 1; fi

wget "https://github.com/costap/tunnel/releases/download/${VERSION}/tunneld_${OS}_${ARCH}" -o tunneld
wget "https://github.com/costap/tunnel/releases/download/${VERSION}/tunnelctl_${OS}_${ARCH}" -o tunnelctl

install -m 0755 tunnelctl /usr/local/bin
install -m 0755 tunneld /usr/local/bin

rm tunnelctl tunneld
