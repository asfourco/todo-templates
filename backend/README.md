# Go backend

## Usage

To run the backend locally, you can use the following command:

```bash
# format & lint
$ make format && make lint

# Clean and build
$ make clean && make build

# or simply
$ make local-be
```

Note any changes to the code will require re-building the binary and re-starting the backend server.

## migration

First, make sure you have the "migrate" tool installed. You can install it using:

```bash
go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate
```

Then, you can run the following command to create a new migration:

```bash
migrate create -ext sql -dir path_to_your_migrations_folder name_of_your_migration
```


