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

## Run with docker

### Build image
```bash
docker build -t sci-review .
```

### Run container
You can pass environment variables to container (see .env.example file
```bash
docker run --network="host" -p 8080:8080 -e DATABASE_URL=postgresql://postgres:postgres@localhost:5432/sci_review sci-review
```

## Run with docker-compose

### Run containers
```bash
docker-compose up --build
```

### Stop containers
```bash
docker-compose down
```
