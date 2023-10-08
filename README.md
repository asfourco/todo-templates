## migration

First, make sure you have the "migrate" tool installed. You can install it using:

```bash
go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate
```

Then, you can run the following command to create a new migration:

```bash
migrate create -ext sql -dir path_to_your_migrations_folder name_of_your_migration
```


