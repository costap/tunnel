# Tunnel 

Tunnel is a project to create tunnels to your local host automatically.

## Project structure

[layout](https://github.com/golang-standards/project-layout)

## Install

1. Create remote host

2. On the remote host (install envoy)[https://www.getenvoy.io/install/envoy/ubuntu/]

`apt-get update`

```
apt-get install -y \
 apt-transport-https \
 ca-certificates \
 curl \
 gnupg-agent \
 software-properties-common
```

`curl -sL 'https://getenvoy.io/gpg' | sudo apt-key add -`

```shell script
add-apt-repository \
  "deb [arch=amd64] https://dl.bintray.com/tetrate/getenvoy-deb \
  $(lsb_release -cs) \
  stable"
```

`sudo apt-get update && sudo apt-get install -y getenvoy-envoy`

`envoy --version`

*OR*

```shell script
curl -L https://getenvoy.io/cli | bash -s -- -b /usr/local/bin
getenvoy run standard:1.14.1 -- --config-path tcp-proxy.yaml
```

3. On the master pi run the tunnel

`ssh -N -R <remoteport>:localhost:<localport> user@server >/dev/null 2>&1 &`