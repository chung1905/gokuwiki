FROM busybox:glibc

RUN mkdir -p /srv/data
WORKDIR /srv

COPY --from=golang:alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./gokuwiki .
COPY ./web web

EXPOSE 8080

CMD ["/srv/gokuwiki"]
