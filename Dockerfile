FROM scratch
LABEL maintainer="Hans-Dieter Stich <info@monkeyguru.dev>"

WORKDIR /app

ADD owm-proxy /app/

EXPOSE 8080

ENTRYPOINT ["/app/owm-proxy"]