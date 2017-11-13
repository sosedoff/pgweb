FROM alpine:3.6
MAINTAINER Dan Sosedoff <dan.sosedoff@gmail.com>

ENV PGWEB_VERSION 0.9.10

RUN \
  apk update && \
  apk add --update ca-certificates openssl && \
  update-ca-certificates && \
  cd /tmp && \
  wget https://github.com/sosedoff/pgweb/releases/download/v$PGWEB_VERSION/pgweb_linux_amd64.zip && \
  unzip pgweb_linux_amd64.zip -d /usr/bin && \
  mv /usr/bin/pgweb_linux_amd64 /usr/bin/pgweb && \
  rm -f pgweb_linux_amd64.zip

EXPOSE 8081
CMD ["/usr/bin/pgweb", "--bind=0.0.0.0", "--listen=8081"]