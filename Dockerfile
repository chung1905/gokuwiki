FROM busybox:glibc

RUN mkdir -p /srv/data
WORKDIR /srv

COPY ./gokuwiki .
COPY ./templates templates
COPY ./pub pub

EXPOSE 8080

CMD ["/srv/gokuwiki"]
