# About 
A go todo app built for demos entirely using ChatGPT and Claude


# Run it as 
```shell
export DATABASE_URL=postgres://USER:PASS@HOST:PORT/DATABSE_NAME?sslmode=disable
# include sslmode=disable if ssl is not enabled on PG server
go run main.go
```


# Docker build 

```shell
docker build -t golang-todo-app .
```

# Docker compose 
Create .env with following keys 
```shell
APP_PORT=8080
POSTGRES_PASSWORD=password
APP_DB_NAME=todo_app
docker compose up -d
```