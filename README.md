# ASG Creator

[Application Security
Groups](http://docs.cloudfoundry.org/adminguide/app-sec-groups.html) (ASGs) are
used to whitelist outbound container network access in [Cloud
Foundry](http://cloudfoundry.org). Application containers require ASGs to
enable outbound network access.

Creating ASGs for the `default-staging` and `default-running` ASG sets can be
intimidating as they are defined as whitelists, but you want to specify
specific IPs and networks to be blacklisted. An example of addresses you would
want to blacklist would be those for VMs running Cloud Foundry system
components, or those for Marketplace Service VMs.

The ASG Creator can be used to create baseline public-networks and
private-networks ASGs that allow all public and private networks *except* those
you want to blacklist. Additionally, it will block by default the
169.254.0.0/16 link-local CIDR.

You are encouraged to modify the files created by ASG Creator to suit your
needs.

## Installation

Download the [latest release](https://github.com/cloudfoundry-incubator/asg-creator/releases/latest).

Alternatively, install from source:

```
go get github.com/cloudfoundry-incubator/asg-creator
```

## Usage

Config Options

* *exclude*: An array of IPs, CIDRs, and IP ranges (e.g. `192.168.100.4`, `192.168.0.0/16`, `192.168.1.1-192.168.100.3`) to exclude
* *include*: An array of IPs, CIDRs, and IP ranges to use as the base from which to remove IPs/CIDRs/IP ranges from

### Creating ASG rules based on a provided list of networks

To create ASG rules starting with a specific set of networks and then subtracting IPs from them, create a config, `config.yaml`:

```yaml
include:
- 10.68.192.0/24

exclude:
- 10.68.192.0
- 10.68.192.127
- 10.68.192.128
- 10.68.192.255
- 10.68.192.50-10.68.192.100
```

Use the config, specifying a custom output filename, to create an ASG rules file:

```
$ asg-creator create --config config.yml --output custom.json
Wrote custom.json
OK
$ cat custom.json
[
	{
		"protocol": "all",
		"destination": "10.68.192.1-10.68.192.49"
	},
	{
		"protocol": "all",
		"destination": "10.68.192.101-10.68.192.126"
	},
	{
		"protocol": "all",
		"destination": "10.68.192.129-10.68.192.254"
	}
]
```

Use the [cf cli](https://github.com/cloudfoundry/cli/releases/latest) to create and bind an ASG with the rules file:

```
$ cf create-security-group my-security-group custom.json
$ cf bind-security-group my-security-group my-org my-space

# restart any apps in my-org/my-space for the new ASG to take effect
$ cf restart my-app-running-in-my-space
```

### Creating ASGs for default-staging and default-running

To create `public-networks.json` and `private-networks.json`, where each file contains all public or private networks respectively, except for specific IPs and networks that are configured, create a config, `config.yml`:

```yaml
exclude:
- 192.168.100.4
- 192.168.1.0/24
- 192.168.200.0-192.168.200.50
```

Use the config to create ASG rules files:

```
$ asg-creator create --config config.yml
Wrote public-networks.json
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
		"destination": "192.168.100.5-192.168.199.255"
	},
	{
		"protocol": "all",
		"destination": "192.168.200.51-192.168.255.255"
	}
]

$ cat public-networks.json
[
	{
		"protocol": "all",
		"destination": "0.0.0.0-9.255.255.255"
	},
	{
		"protocol": "all",
		"destination": "11.0.0.0-169.253.255.255"
	},
	{
		"protocol": "all",
		"destination": "169.255.0.0-172.15.255.255"
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

Modify the ASG rules files per your network policy for application containers
running untrusted code.

Use the rules files to create ASGs with the [cf
CLI](https://github.com/cloudfoundry/cli/releases/latest):

```
$ cf create-security-group private-networks private-networks.json
$ cf bind-staging-security-group private-networks
$ cf bind-running-security-group private-networks

$ cf create-security-group public-networks public-networks.json
$ cf bind-staging-security-group public-networks
$ cf bind-running-security-group public-networks
```
