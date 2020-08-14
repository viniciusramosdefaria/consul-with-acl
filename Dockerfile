FROM alpine:latest

COPY consul-test .

ENTRYPOINT ["/consul-test"]