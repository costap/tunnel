# Tunnel 

Tunnel is a project written in go to create L4 tunnels to your local host automatically.

## How it works

The `tunnelctl create host` command will provision a new server in DigitalOcean and install [envoy](https://www.envoyproxy.io/);
 this will be your public proxy and will be configured to forward requests to an internal port.
 
Once the public host is created, from the local machine you want to expose, running `tunneld` will create and maintain 
a ssh tunnel to the public host to tunnel the requests from envoy to your local address.

## Install

```shell script
curl -L https://raw.githubusercontent.com/costap/tunnel/master/scripts/install.sh | sudo bash -
```

## Usage

The below example will proxy ports 443 and 6443 from the public host to local addresses 192.168.0.26:443 and 
192.168.0.26:6443.

1. setup a config file in `~/.tunnelctl.yaml` like `configs/tunnelctl.yaml`
2. create ssh keys pair if you don't have one
```shell script
tunnelctl keys create -p ~/.ssh -n id_rsa_tunnel
```
3. create public host in DO
```shell script
tunnelctl hosts create -p ~/.ssh --sshName id_rsa_tunnel --proxy 443:10443 --proxy 6443:16443 --name tunnel-proxy
```
_take note of new host external IP and replace <NEWHOSTIP> below_

4. start the tunnels on your local machine
```shell script
nohup tunneld -c ~/.ssh/id_rsa_tunnel \
  --sshServer root@<NEWHOSTIP> \
  --localAddr 192.168.0.26:443 \
  --remoteAddr 0.0.0.0:10443 \
  --adminPort 8080 > /dev/null 2>&1 & 
nohup tunneld -c ~/.ssh/id_rsa_tunnel \
  --sshServer root@<NEWHOSTIP> \
  --localAddr 192.168.0.26:6443 \
  --remoteAddr 0.0.0.0:16443 \
  --adminPort 8081 > /dev/null 2>&1 & 
```

5. Check tunnels are ok
```shell script
$ curl http://localhost:8080/health
{ "started": true }                                                                                                                                                       
$ curl http://localhost:8081/health
{ "started": true }
```

5. Check external endpoints
```shell script
$ curl https://<NEWHOSTIP>/
                                                                                                                                                
$ curl http://<NEWHOSTIP>:6443/
```

6. Enjoy!

## Build

To build the project locally simply run `make`.