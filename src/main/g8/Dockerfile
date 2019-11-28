## Sample build command:
## docker build --force-rm --squash -t demo:0.4.r1 .

FROM golang:1.12-alpine AS builder
MAINTAINER $app_author$
RUN apk add git \
    && mkdir /build
COPY . /build
RUN cd /build && go build -o main

FROM alpine:3.10
RUN mkdir /app
COPY --from=builder /build/main /app/main
COPY --from=builder /build/README.md /app/README.md
COPY --from=builder /build/config /app/config
RUN apk add --no-cache -U tzdata bash ca-certificates \
    && update-ca-certificates \
    && cp /usr/share/zoneinfo/$timezone$ /etc/localtime \
    && chmod 711 /app/main \
    && rm -rf /var/cache/apk/*
WORKDIR /app
CMD ["/app/main"]
#ENTRYPOINT /app/main
