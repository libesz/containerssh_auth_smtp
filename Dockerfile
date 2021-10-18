FROM golang
COPY . /code
WORKDIR /code
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' .

FROM alpine
COPY --from=0 /code/containerssh_smtp_auth /
ENTRYPOINT ["/containerssh_smtp_auth"]  