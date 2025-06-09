## ingresso.go

An API for practicing Go

This project will be a ticket selling manager for movie theater

#### Running the project

Install dependencies

```bash
go mod tidy
```

Run dev server

```bash
air
```

Build and run compiled

```bash
go build -o go-api /app/api/main.go
```

```
./go-api
```

#### Running with docker

```bash
docker compose up -d
```

#### Endpoints


#### Features

- Send transactional emails using SES
    - Upcoming movie session
    - Empty seats on movie