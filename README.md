# ContainerSSH SMTP authenticator (unofficial)
This project is an auth + config webhook server implementation for [ContainerSSH](https://containerssh.io/) with the following capabilities:
* The authentication is proxied to an SMTP server, which will provide the real authentication. STARTTLS is expected.
* The configuration server part is configured with a static yaml file, which provides docker volume to user mappings.
  * The configuration server part assumes that only authenticated users with at least one volume ownership will make it through the config requests (the auth part also uses the same mapping to determine if the user has access to any volume)
  * The configuration server sets up the SFTP client container to mount the configured docker volume

## Disclaimer
The server is right now is only able to provide plain HTTP service. It is expected to run inside an isolated environment, like a docker network, where only ContainerSSH is available.
Do not expose this to the public internet.

## Use case
Assuming you run a webhosting + email service, you most probably provide SMTP authentication for the users to restrict email sending. This project helps to not implement another authentication service (and database, and a second password to keep in mind by the user, etc) for reaching the webhosting content. Users are able to use SFTP to reach their content with their email credentials.

## Usage:
With plain docker run command:
```
docker run -e LISTEN_ON=0.0.0.0:8090 -e SMTP_EP=<SMTP_SERVER_IP_OR_HOSTNAME>:587 -e SMTP_SERVER_NAME=<SMTP_SERVER_NAME_THAT_MATCHES_ITS_TLS_CERT> -e USER_VOLUME_MAPPING_PATH:=/mapping.yaml -v <MAPPING_FILE_PATH>:/mapping.yaml -p 8090:8090 huszty/containerssh_auth_smtp:v0.2.0
```

With docker-compose, together with the SSH server container:
```yaml
version: '3.2'
services:
  containerssh:
  [...]
  authconfig:
    image: huszty/containerssh_auth_smtp:v0.2.0
    environment:
      LISTEN_ON: "0.0.0.0:8090"
      SMTP_EP: "<SMTP_SERVER_IP_OR_HOSTNAME>:587"
      SMTP_SERVER_NAME: "<SMTP_SERVER_NAME_THAT_MATCHES_ITS_TLS_CERT>:587"
      USER_VOLUME_MAPPING_PATH: "/mapping.yaml"
    volumes:
      - ./mapping.yaml:/mapping.yaml
```

Example `mapping.yaml` file:
```yaml
volumeprefix: sites_
volumes:
- volumename: example-com
  users:
  - admin@example.com

```
This example will mount the sites_example-com docker volume under `/content/example-com` in the SFTP session and puts the user (with the username admin@example.com) under `/content` when successfully logging in. You may define multiple volumes with multiple owners if needed.

## Example ContainerSSH setting snippet
```yaml
[...]
auth:
  url: "http://authconfig:8090/auth"
  timeout: 10s
configserver:
  url: "http://authconfig:8090/config"
  timeout: 10s
[...]
docker:
  connection:
    host: unix:///var/run/docker.sock
  execution:
    container:
      networkdisabled: true
    host:
      readonlyrootfs: true
security:
  maxSessions: 1
  defaultMode: filter
  subsystem:
    allow:
      - sftp
```

## Auth-only mode
If you do not specify the mappind file (with env var USER_VOLUME_MAPPING_PATH), the server starts up without the config server handler. Authentication server will purely SMTP authenticate the user without further checks. In this case you have to provide some other config server for ContainerSSH. 

## Build:
```
docker build -t huszty/containerssh_auth_smtp .
```

## TODO
Planned:
 * Implement hot-reload for the mapping file change (now you have to restart the container to catch up)
 * Add TLS to have HTTPS
 * Add more unit test

Outside of my use-case, but still interested in:
 * Add other backends, like Kubernetes
 * Anything else which does not break the basic idea
