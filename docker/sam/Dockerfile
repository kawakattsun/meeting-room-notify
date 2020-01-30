FROM python:3.8.0-alpine3.10

ENV PYTHONUSERBASE=/usr/local

COPY packages_requirements.txt ./

RUN set -ex \
  && apk add --no-cache \
    bash \
    gcc \
    make \
    musl-dev \
    jq \
    docker-cli

RUN set -ex \
  && pip install --upgrade pip \
  && pip install --user -r packages_requirements.txt
