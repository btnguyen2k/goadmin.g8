# goadmin.g8

Giter8 template to develop `Admin Control Panel` in Go using Echo framework.

Latest release: [template-v0.4.r2](RELEASE-NOTES.md).

## Features

- Create new project from template with [go-giter8](https://github.com/btnguyen2k/go-giter8).
- Built on [Echo framework v4](https://echo.labstack.com).
- `Landing page` using [Greyscale template](https://startbootstrap.com/themes/grayscale/).
- `Admin Control Panel` using [AdminLTE v3 - Bootstrap Admin Dashboard Template](https://adminlte.io):
  - Login page & Logout
  - Dashboard
  - Profile page & Change password
  - User & User group management (list, create, update, delete)
  - BO & DAO implementation using [SQLite3](https://github.com/mattn/go-sqlite3)
- Sample `.gitlab-ci.yaml` & `Dockerfile` to package application as Docker image.


## Getting Started

### Install `go-giter8`

This a Giter8 template, so it is meant to be used in conjunction with a giter8 tool.
Since this is a template for Go application, it make sense to use [go-giter8](https://github.com/btnguyen2k/go-giter8).

See [go-giter8](https://github.com/btnguyen2k/go-giter8) website for installation guide.

### Create new project from template

```
g8 new btnguyen2k/goadmin.g8
```

and follow the instructions.

> Note: This template requires `go-giter8` version `0.3.2` or higher.

### Write application code

Directory `src/goadmin` is reserved for `GoAdmin` framework, do _not_ put application source code there.

Source code under directory `src/myapp` is the sample `Admin Control Panel` itself.
It is a good starting point, so feel free to reference or modify to build your own application.


## LICENSE & COPYRIGHT

See [LICENSE.md](LICENSE.md) for details.


## Giter8 template

For information on giter8 templates, please see http://www.foundweekends.org/giter8/
