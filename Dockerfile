## Sample build command:
## docker build --force-rm --squash -t goadmin .

FROM golang:1.13-alpine AS builder
MAINTAINER Thanh Nguyen <btnguyen2k@gmail.com>
RUN apk add build-base git \
    && mkdir /build
COPY . /build
#START symlink fixes
RUN rm -rf /build/config/conf.d && mkdir -p /build/config/conf.d \
    && rm -rf /build/config/i18n_myapp && mkdir -p /build/config/i18n_myapp \
    && rm -rf /build/views && mkdir -p /build/views \
    && rm -rf /build/public && mkdir -p /build/public
COPY ./src/main/g8/config/ /build/config/
COPY ./src/main/g8/views/ /build/views/
COPY ./src/main/g8/public/ /build/public/
#END
#START naming
RUN  cd /build \
    && sed -i 's/\$shortname\$/goadmin/g' config/*.conf \
    && sed -i 's/\$name\$/GoAdmin/g' config/*.conf \
    && export BUILD=`date +%Y%m%d%H` && sed -i 's/\$version\$/0.1.1b'$BUILD/g config/*.conf \
    && sed -i 's/\$desc\$/AdminCP Giter8 template for GoLang/g' config/*.conf
#END
RUN cd /build && go build -o main

FROM alpine:3.10
RUN mkdir /app
COPY --from=builder /build/main /app/main
COPY --from=builder /build/README.md /app/README.md
COPY --from=builder /build/config /app/config
COPY --from=builder /build/views /app/views
COPY --from=builder /build/public /app/public
RUN apk add --no-cache -U tzdata bash ca-certificates \
    && update-ca-certificates \
    && cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime \
    && chmod 711 /app/main \
    && rm -rf /var/cache/apk/*
WORKDIR /app
CMD ["/app/main"]
#ENTRYPOINT /app/main
