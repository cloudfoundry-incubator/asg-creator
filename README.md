# ASG Creator

This tool can be used to create Application Security Groups for use in Cloud Foundry.

## Installation

```
go get github.com/cloudfoundry-incubator/asg-creator
```

## Usage

Options

* *excluded_ips*: An array of IPs to exclude
* *excluded_networks*: An array of CIDRs to exclude

Create a config:

```yaml
excluded_ips:
- 192.168.100.4

excluded_networks:
- 192.168.1.0/24
```

Use the config to create ASG configuration files:

```
$ asg-creator create --config config.yml
Wrote public-networks.json
Wrote private-networks.json

OK
$ cat private-networks.json | jq '.'
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

$ cat public-networks.json | jq '.'
[
  {
    "protocol": "all",
    "destination": "0.0.0.0-9.255.255.255"
  },
  {
    "protocol": "all",
    "destination": "11.0.0.0-169.254.169.253"
  },
  {
    "protocol": "all",
    "destination": "169.254.169.255-172.15.255.255"
  },
  {
    "protocol": "all",
    "destination": "172.32.0.0-192.167.255.255"
  },
  {
    "protocol": "all",
    "destination": "192.169.0.0-255.255.255.255"
  }
]
```

Use the configuration files to create ASGs:

```
$ cf create-security-group private-networks private-networks.json
$ cf create-security-group public-networks public-networks.json
```
