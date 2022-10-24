# $name$

**$desc$** by **$author$ - $organization$**, based on [Go Admin template](https://github.com/btnguyen2k/goadmin.g8).

Copyright (C) by **$organization$**.

Latest release version: `$version$`. See [RELEASE-NOTES.md](RELEASE-NOTES.md).

## Getting Started

### Application Bootstrapper & API Handlers

Implement application bootstrapper by implementing interface `goadmin.IBootstrapper`.

API handler is defined as

```
type IApiHandler func(*ApiContext, *ApiAuth, *ApiParams) *ApiResult`
```

See [samples](src/samples) for example of bootstrapper and API handler.

### Configurations

Default application configuration file is [config/application.conf](config/application.conf), override default value via env `APP_CONFIG`.
Configurations are in [HOCON format](https://github.com/lightbend/config/blob/master/HOCON.md).

Important configurations:

**Application information**

```
app {
  name: "Application's name"
  shortname: "Application's short name, e.g. goadmin"
  version: "Application's version string, e.g. 1.0.0"
  desc: "Application's description, e.g. Caching services with RESTful APIs"
}
```

**API settings**

```
api {
  http {
    # Listen address & port for HTTP/REST API gateway.
    # set listen_port=0 to disable HTTP/REST API gateway.
    listen_addr = "0.0.0.0"
    listen_port = 8080

    # Name of HTTP header that holds "application id" info passed from client.
    header_app_id = "X-App-Id"

    # Name of HTTP header that holds "access token" info passed from client.
    # override this setting with env HTTP_HEADER_ACCESS_TOKEN
    header_access_token = "X-Access-Token"
  }

  grpc {
    # Listen address & port for gRPC API gateway.
    # set listen_port=0 to disable gRPC API gateway.
    listen_addr = "0.0.0.0"
    listen_port = 8090
  }

  # Client cannot send request that exceeds this size
  # - absolute number: size in bytes
  # - or, number+suffix: https://github.com/lightbend/config/blob/master/HOCON.md#size-in-bytes-format
  max_request_size = 64KiB

  # Timeout to parse request data
  # - absolute number: time in milliseconds
  # - or, number+suffix: https://github.com/lightbend/config/blob/master/HOCON.md#duration-format
  # override this setting with env API_REQUEST_TIMEOUT
  request_timeout = 10s
}
```

**REST Endpoints**

```
api {
  http {
    # API HTTP endpoints
    endpoints {
      # format: {url={http-method=handler-name}}
      "/info" {
        get = "info"
      }
      "/echo/:id" {
        get    = "echo"
        post   = "echo"
        put    = "echo"
        delete = "echo"
      }
    }
  }
}
```

## LICENSE & COPYRIGHT

See [LICENSE.md](LICENSE.md).
