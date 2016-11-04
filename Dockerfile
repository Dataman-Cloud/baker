FROM gliderlabs/alpine:3.3
MAINTAINER Fan Bin <bfan@dataman-inc.com>

WORKDIR /opt/baker

RUN mkdir -p /opt/baker /etc/baker
COPY config.yml /etc/baker/
COPY bin/baker /opt/baker

ENV CONFIG_PATH /etc/baker/config.yml
ENTRYPOINT ["/opt/baker"]

EXPOSE 3000

ARG VERSION
LABEL version="$VERSION"
