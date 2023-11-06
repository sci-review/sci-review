# sci-review

sci-review is your streamlined tool for systematic review. Organize, analyze, and collaborate with ease. Elevate your research.

## How to run (dev)

### migrate database

1. Install CLI [golang-migrate](https://github.com/golang-migrate/migrate) download: https://packagecloud.io/golang-migrate/migrate
2. Execute migrate command:
```bash
 migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/sci_review?sslmode=disable" -verbose up
```

Create new migrations:
```bash
 migrate create -ext sql -dir db/migrations -seq <migration_name>
```