FROM golang:1.14-alpine as builder

ARG APP_VERSION
ARG GIT_BRANCH
ARG GIT_EMAIL
ARG GIT_REVISION

RUN apk add --no-cache \
    gcc \
    make \
    musl-dev

COPY . /src
RUN make -C /src vendor \
    && CGO_ENABLED=0 make \
        -C /src install \
        PREFIX=/pkg \
        GO_BUILDFLAGS='-mod vendor' \
        GOLDFLAGS="-w -linkmode external -extldflags -static"

FROM alpine:3.12

LABEL maintainer="Ben Cartwright-Cox <bencartwrightcox@monzo.com>"

COPY --from=builder /pkg/ /usr/

EXPOSE 9481
ENTRYPOINT ["/usr/bin/drbd9_exporter"]
