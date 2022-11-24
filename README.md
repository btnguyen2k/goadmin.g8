# goadmin.g8

[![Actions Status](https://github.com/btnguyen2k/goadmin.g8/workflows/myapp/badge.svg)](https://github.com/btnguyen2k/goadmin.g8/actions)
[![codecov](https://codecov.io/gh/btnguyen2k/goadmin.g8/branch/master/graph/badge.svg?token=HVAP5A0R2Z)](https://codecov.io/gh/btnguyen2k/goadmin.g8)
[![Release](https://img.shields.io/github/release/btnguyen2k/goadmin.g8.svg?style=flat-square)](RELEASE-NOTES.md)

Giter8 template to build `Admin Control Panel` for Go.

Demo: https://demo-goadmin.gpvcloud.com/.

See also: [govueadmin.g8](https://github.com/btnguyen2k/govueadmin.g8) - single-page VueJS-based `Admin Control Panel` application template for Go.

## Features

- [Giter8](https://github.com/btnguyen2k/go-giter8) template.
- Built on [Echo framework v4](https://echo.labstack.com).
- `Landing page` using [Greyscale template](https://startbootstrap.com/theme/grayscale).
- `Admin Control Panel` using [AdminLTE v3 - Bootstrap Admin Dashboard Template](https://adminlte.io):
  - User signin & signout
  - Dashboard
  - Profile page & Change password
  - User & User group management (list, create, update, delete)
  - BO & DAO implementation in SQLite3, MySQL, PostgreSQL and MongoDB
  - Unit tests for BO & DAO
- I18n support.
- Sample `Dockerfile` to package application as Docker image.
- Sample [GitHub Actions](https://docs.github.com/actions) workflow.


## Getting Started

### Install `go-giter8`

This a Giter8 template, so it is meant to be used in conjunction with a giter8 tool.
Since this is a template for Go application, it makes sense to use [go-giter8](https://github.com/btnguyen2k/go-giter8).

See [go-giter8](https://github.com/btnguyen2k/go-giter8) website for installation guide.

### Create new project from template

```
$ g8 new btnguyen2k/goadmin.g8
```

and follow the instructions.

> Note: This template requires `go-giter8` version `0.3.2` or higher.

### Write application code

Directory `src/goadmin` is reserved for `GoAdmin` framework, do _not_ put application source code there.

Source code under directory `src/myapp` is the sample `Admin Control Panel` itself.
It is a good starting point, feel free to reference or modify to build your own application.


## LICENSE & COPYRIGHT

See [LICENSE.md](LICENSE.md) for details.


## Giter8 template

For information on giter8 templates, please see http://www.foundweekends.org/giter8/
