FROM kong:3.0-alpine

USER root

COPY ./plugins/build/kong-keyless /usr/local/bin/kong-keyless