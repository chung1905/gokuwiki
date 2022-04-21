FROM busybox:glibc

RUN mkdir -p /srv/data
WORKDIR /srv

COPY ./gokuwiki .
COPY ./templates templates

EXPOSE 8080

CMD ["/srv/gokuwiki"]
