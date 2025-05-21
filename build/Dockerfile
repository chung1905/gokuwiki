FROM busybox:glibc

RUN mkdir -p /srv/data
WORKDIR /srv

COPY --from=golang:alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080

COPY ./web web
COPY ./gokuwiki .

CMD ["/srv/gokuwiki"]
