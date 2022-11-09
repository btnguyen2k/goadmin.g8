## Sample build command:
## docker build --force-rm --squash -t $shortname$:$version$ .

FROM golang:1.17-alpine AS builder
LABEL maintainer="$author$"
RUN apk add build-base git \
    && mkdir /build
COPY . /build
RUN cd /build && go build -o main

FROM alpine:3
LABEL maintainer="$author$"
RUN mkdir /app
COPY --from=builder /build/main /app/main
COPY --from=builder /build/config /app/config
COPY --from=builder /build/views /app/views
COPY --from=builder /build/public /app/public
RUN apk add --no-cache -U tzdata bash ca-certificates \
    && update-ca-certificates \
    && cp /usr/share/zoneinfo/$timezone$ /etc/localtime \
    && chmod 711 /app/main \
    && rm -rf /var/cache/apk/*
WORKDIR /app
CMD ["/app/main"]
#ENTRYPOINT /app/main
