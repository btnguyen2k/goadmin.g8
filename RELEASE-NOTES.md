# goadmin.g8 Release Notes

## 2022-11-xx: template-r3

Frontend:
- `startbootstrap-grayscale v7.0.5`
- `AdminLTE v3.2.0`

## 2020-02-06: template-v0.4.r2

- PostgreSQL implementation of BO & DAO.
- Reduce size of `public` directory.
- (for demo purpose) Changing password and modifying of system admin account are disabled.
- CircleCI config file to build demo site on Heroku.


## 2019-12-30: template-v0.4.r1

- `Echo v4.1.x`.
- `Landing page` using [Greyscale template](https://startbootstrap.com/themes/grayscale/).
- `Admin Control Panel` using [AdminLTE v3 - Bootstrap Admin Dashboard Template](https://adminlte.io):
  - Login page & Logout
  - Dashboard
  - Profile page & Change password
  - User & User group management (list, create, update, delete)
  - BO & DAO implementation using [SQLite3](https://github.com/mattn/go-sqlite3)
- Sample `.gitlab-ci.yaml` & `Dockerfile`.
