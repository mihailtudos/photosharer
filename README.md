# PhotoSharer

## Docker 

To start the postgress instance within a Docker container from the root directory of the project run:

```shell
    docker compose up
```

Note: add the -d flag to run in detached mode

To interact with the container running the postgres instance run 

```shell
    docker compose exec -it SERVICE_NAME psql -U USER_NAME -d DATABASE_NAME
```

To interact with the database a Go application essentially needs two things: 

1. a library to interact with the database using SQL
2. a driver for the database 

Go provides an abstraction layer through the database/sql library which is what is used in this project but better solutions exist out there like the **pgx** library that optimises a lot when working with postgres.

```shell
    go get github.com/jackc/pgx/v4
```