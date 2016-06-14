# ASG Creator

This tool can be used to create Application Security Groups for use in Cloud Foundry.

## Installation

```
go get github.com/cloudfoundry-incubator/asg-creator
```

## Usage

Options

* *private_networks*: Set to `true` to write out all private networks less any blacklisted IPs/networks
* *public_networks*: Set to `true` to write out all public networks less any blacklisted IPs/networks
* *excluded_ips*: An array of IPs to exclude
* *excluded_networks*: An array of CIDRs to exclude

Create a config:

```yaml
private_networks: true

excluded_ips:
- 192.168.100.4

excluded_networks:
- 192.168.1.0/24
```

Use the config to create ASG configuration files:

```
$ asg-creator create --config config.yml
Wrote private-networks.json

OK
$ cat private-networks.json
[
  {
    "protocol": "all",
    "destination": "10.0.0.0-10.255.255.255"
  },
  {
    "protocol": "all",
    "destination": "172.16.0.0-172.31.255.255"
  },
  {
    "protocol": "all",
    "destination": "192.168.0.0-192.168.0.255"
  },
  {
    "protocol": "all",
    "destination": "192.168.2.0-192.168.100.3"
  },
  {
    "protocol": "all",
    "destination": "192.168.100.5-192.168.255.255"
  }
]
```

Use the configuration files to create ASGs:

```
$ cf create-security-group my-security-group private-networks.json
```
