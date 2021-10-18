# containerssh_smtp_auth

## Build:
```
docker build -t huszty/containerssh_auth_smtp .
```

## Usage:
```
docker run -e LISTEN_ON=0.0.0.0:8090 -e SMTP_EP=<SMTP_SERVER_IP_OR_HOSTNAME>:587 -e SMTP_SERVER_NAME=<SMTP_SERVER_NAME_THAT_MATCHES_ITS_TLS_CERT> -p 8090:8090 huszty/containerssh_temp_auth
```
