# API Service

The main api-service for the application built using:
- Go: language
- Gin: http server
- SQLC: query builder
- PGX: postgre adapter
- Go Migrate: database migrations
- gocraft/work: background jobs and processing
- gomail: mail processing

The `Makefile` contains all the necessary commands for running and setting
up the service. 

All the communication between the services is done using HTTP requests.

For the service to work it requires: 
- A redis or redis compliant(Valkey, DragonflyDB, etc...) instance that will be used as a background worker.
- The embedding service. You can create your own implementation of the embedding service or use the one provided with the `docker-compose` file
- The categorizer service. You can create your own implementation of the categorizer service or use the one provided with the `docker-compose` file

# To Run Locally

## *Prerequisites*:
- Please create a `app.env` file using the same parameters from the given `app.env.example` file. 
- Please install go-migrate to be able to run the migrations. [Install Instructions](https://github.com/golang-migrate/migrate)
- If you make changes to the database and queries, you'll need SQLC [Install Instructions](https://docs.sqlc.dev/en/stable/overview/install.html)

## *Steps to run*:
1. Create a docker network using the **Makefile** - `make network`
2. Create a postgres instance using the **Makefile** - `make postgres`
3. Create the database using the **Makefile** - `make createdb`
4. Run the migrations using the **Makefile** - `make migrateup`
5. Run the application using the **Makefile** - `make server`