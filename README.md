# containerssh_smtp_auth
This project is an auth + config server implementation for ContainerSSH with the following capabilities:
* The authentication is proxied to an SMTP server, which will provide the real authentication. STARTTLS is expected.
* The configuration server part is configured with a static yaml file, which provides docker volume to user mappings.
  * The configuration server part assumes that only authenticated users with at least one volume ownership will make it through the config requests (the auth part also uses the same mapping to determine if the user has access to any volume)

## Build:
```
docker build -t huszty/containerssh_auth_smtp .
```

## Usage:
```
docker run -e LISTEN_ON=0.0.0.0:8090 -e SMTP_EP=<SMTP_SERVER_IP_OR_HOSTNAME>:587 -e SMTP_SERVER_NAME=<SMTP_SERVER_NAME_THAT_MATCHES_ITS_TLS_CERT> -e USER_VOLUME_MAPPING_PATH:=/mapping.yaml -v <MAPPING_FILE_PATH>:/mapping.yaml -p 8090:8090 huszty/containerssh_temp_auth
```

Example mapping file:
```
volumeprefix: sites_
volumes:
- volumename: example-com
  users:
  - admin@example.com

```


## Example ContainerSSH setting snippet
```
auth:
  url: "http://authconfig:8090/auth"
  timeout: 10s
configserver:
  url: "http://authconfig:8090/config"
  timeout: 10s
```
