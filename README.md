# Tunnel 

Tunnel is a project to create tunnels to your local host automatically.

## How it works

The `tunnelctl` create host command will provision a new server in DigitalOcean and install [envoy]()
 to proxy the defined ports to a internal port.
 
Running `tunneld` from you local machine you want to expose will maintain and ssh tunnel to the DO host on the defined
ports and proxy all requests through to your local address.

## Project structure

[layout](https://github.com/golang-standards/project-layout)

## Install

```shell script
go get github.com/costap/tunnel
```

## Usage

The below example will proxy ports 443 and 6443 from the external host to local addresses 192.168.0.26:443 and 
192.168.0.26:6443

1. Setup a config file in `~/.tunnelctl.yaml` like `configs/tunnelctl.yaml`.
1. Create ssh keys pair if you don't have one.
```shell script
tunnelctl keys create -p ~/.ssh -n id_rsa_tunnel.
```
1. Create entry host in DO (assuming you want to expose 443 and 6443 in host)
```shell script
tunnelctl hosts create -p ~/.ssh --sshName id_rsa_tunnel --proxy 443:10443 --proxy 6443:16443 --name tunnel-proxy
```
Take note of new host external IP and replace <NEWHOSTIP> next.
1. Start the tunnels
```shell script
tunneld -c ~/.ssh/id_rsa_tunnel --sshServer root@<NEWHOSTIP> --localAddr 192.168.0.26:443 --remoteAddr 0.0.0.0:10443 
tunneld -c ~/.ssh/id_rsa_tunnel --sshServer root@<NEWHOSTIP> --localAddr 192.168.0.26:6443 --remoteAddr 0.0.0.0:16443 
```
