---
layout: page
nav_order: 3
---
# Database

Aggregator is configured to use SQLite3 DB by default, but it also supports PostgreSQL.
In CI and QA environments, the configuration is overridden by environment variables to use
PostgreSQL.

To establish connection to the PostgreSQL instance provided by the minimal stack in
`docker-compose.yml` for local setup, the following configuration options need to be changed in
`storage` section of `config.toml`:

```toml
[storage]
db_driver = "postgres"
pg_username = "user"
pg_password = "password"
pg_host = "localhost"
pg_port = 55432
pg_db_name = "aggregator"
pg_params = "sslmode=disable"
```

## Migration mechanism

This service contains an implementation of a simple database migration mechanism that allows
semi-automatic transitions between various database versions as well as building the latest version
of the database from scratch.

The migrations are no longer performed during initialization of the service to make running multiple
instances of the service in parallel against the same database safer. Instead, the database
migration version can now be set using the built-in CLI sub-command `migration` (aliases:
`migrations` and `migrate`).

### Printing information about database migrations

```shell
./insights-results-aggregator migrations
```

### Upgrading the database to the latest available migration

```shell
./insights-results-aggregator migration latest
```

### Downgrading to the base (empty) database migration version

```shell
./insights-results-aggregator migration 0
```

Before using the migration mechanism, it is first necessary to initialize the migration information
table `migration_info`. This can be done using the `migration.InitInfoTable(*sql.DB)` function. Any
attempt to get or set the database version without initializing this table first will result in a
`no such table: migration_info` error from the SQL driver.

New migrations must be added manually into the code, because it was decided that modifying the list
of migrations at runtime is undesirable.

To migrate the database to a certain version, in either direction (both upgrade and downgrade), use
the `migration.SetDBVersion(*sql.DB, migration.Version)` function.

**To upgrade the database to the highest available version, use
`migration.SetDBVersion(db, migration.GetMaxVersion())`.** This will automatically perform all the
necessary steps to migrate the database from its current version to the highest defined version.

See `/migration/migration.go` documentation for an overview of all available DB migration
functionality.
