### Database Setup
The application uses PostgreSQL as its database, running inside a Docker container. This setup ensures easy deployment, portability, and reproducibility of the development environment.

The database image and volume were created as following:

```
  docker run \
  --name [container_name] \
  -e POSTGRES_DB=[db_name] \
  -e POSTGRES_PASSWORD=[password] \
  -p 5432:5432 \
  -v [volume_name]:/var/lib/postgresql/data \
  -d postgres:14.2
```

The container management was done using LazyDocker.
