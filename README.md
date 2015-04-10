Metadata Service
================

This code creates a binary web service for metadata support within a network containing a CoreOS cluster.

Machines within the network can curl this metadata service for network-specific information.

For example, you can install CoreOS on /dev/sda by simply running (replace metadata-service with the hostname of your installation of the metadata-service):

```
eval "`curl -s metadata-service:8080/metadata/v1/install`"
```

Other endpoints include:

```
/metadata/v1/user-data
/metadata/v1/ip
/metadata/v1/hostname
```

It is important to remember, these endpoints will likely work only within your network (otherwise the ip addresses won't be real). 


Configuration
=============
Customize your `cloud-config.template.yaml` file however you'd like it. This file will be used to generate your user-data endpoint.

Create a `custom.json` file with the variables that match your setup. Use `custom.json.example` to get started. These variables are avilable to be substituted into your `cloud-config.template.yaml`.

Installation
============
```
go build
```

Simply copy the resulting binary to your metadata host and execute it to start your service.


