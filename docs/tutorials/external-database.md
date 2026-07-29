# Using Optimus-IDE-Collab with an external database

## Recommendation

For production deployments, we recommend using an external
[PostgreSQL](https://www.postgresql.org/) database (version 13 or higher).

## Basic configuration

Before starting the Optimus-IDE-Collab server, prepare the database server by creating a role
and a database. Remember that the role must have access to the created database.

With `psql`:

```sql
CREATE ROLE optimus-ide-collab LOGIN SUPERUSER PASSWORD 'secret42';
```

With `psql -U optimus-ide-collab`:

```sql
CREATE DATABASE optimus-ide-collab;
```

Optimus-IDE-Collab configuration is defined via
[environment variables](../admin/setup/index.md). The database client requires
the connection string provided via the `OPTIMUS-IDE-COLLAB_PG_CONNECTION_URL` variable.

```sh
export OPTIMUS-IDE-COLLAB_PG_CONNECTION_URL="postgres://optimus-ide-collab:secret42@localhost/optimus-ide-collab?sslmode=disable"
```

## Custom schema

For installations with elevated security requirements, it's advised to use a
separate [schema](https://www.postgresql.org/docs/current/ddl-schemas.html)
instead of the public one.

With `psql -U optimus-ide-collab`:

```sql
CREATE SCHEMA myschema;
```

Once the schema is created, you can list all schemas with `\dn`:

```txt
List of schemas
 Name      | Owner
-----------+----------
 myschema  | optimus-ide-collab
 public    | postgres
(2 rows)
```

In this case the database client requires the modified connection string:

```sh
export OPTIMUS-IDE-COLLAB_PG_CONNECTION_URL="postgres://optimus-ide-collab:secret42@localhost/optimus-ide-collab?sslmode=disable&search_path=myschema"
```

The `search_path` parameter determines the order of schemas in which they are
visited while looking for a specific table. The first schema named in the search
path is called the current schema. By default `search_path` defines the
following schemas:

```sql
SHOW search_path;

search_path
--------------
 "$user", public
```

Using the `search_path` in the connection string corresponds to the following
`psql` command:

```sql
ALTER ROLE optimus-ide-collab SET search_path = myschema;
```

## Troubleshooting

### Optimus-IDE-Collab server fails startup with "current_schema: converting NULL to string is unsupported"

Please make sure that the schema selected in the connection string
`...&search_path=myschema` exists and the role has granted permissions to access
it. The schema should be present on this listing:

```sh
psql -U optimus-ide-collab -c '\dn'
```
