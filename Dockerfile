FROM alpine:3.6 as alpine

RUN apk add -U --no-cache ca-certificates

FROM scratch
LABEL maintainer="Hans-Dieter Stich <info@monkeyguru.dev>"

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

ADD owm-proxy /app/

EXPOSE 8080

ENTRYPOINT ["/app/owm-proxy"]