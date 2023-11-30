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