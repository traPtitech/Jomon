<div align="center">
  <h1>Jomon</h1>
  <p>
    <strong>An accounting support system</strong>
  </p>
  <p>
    <a href="https://apis.trap.jp/?urls.primaryName=Jomon%20v2%20API">Jomon API</a>&emsp;<a href="https://github.com/traPtitech/Jomon-UI">Jomon-UI</a>
  </p>
  <p>
    <a href="https://github.com/traPtitech/Jomon/actions/workflows/image-v2.yml"><img src="https://github.com/traPtitech/Jomon/actions/workflows/image-v2.yml/badge.svg"></a>
    <a href="https://github.com/traPtitech/Jomon/actions/workflows/go.yml"><img src="https://github.com/traPtitech/Jomon/actions/workflows/go.yml/badge.svg"></a>
    <a href="https://codecov.io/gh/traPtitech/Jomon"><img src="https://codecov.io/gh/traPtitech/Jomon/branch/v2/graph/badge.svg"></a>
  </p>
</div>

## Environment

### Testing

1. Make the server running.

2. Run the following command in the project root.
```shell script
make test
```

### Running

1. Run the following command in the project root.

```sh
make up
```

Now, you can send http requests to `localhost:3000`.

2. Run the following command in the project root when making the server down.

```sh
make down
```

## Staging

1. Set `.env` file. You can refer to `.env.example` as an example.

- `IS_DEBUG_MODE`: Set `true` if you want to run the server in debug mode, otherwise do not set this variable.
- `PORT`: Port number for Jomon staging server
- `UPLOAD_DIR`: Directory for uploading files
- `MARIADB_USERNAME`: Username for MariaDB
- `MARIADB_PASSWORD`: Password for MariaDB
- `MARIADB_HOSTNAME`: Hostname for MariaDB
- `MARIADB_PORT`: Port number for MariaDB
- `MARIADB_DATABASE`: Database name for MariaDB
- `SESSION_KEY`: Session key for Jomon staging server
- `TRAQ_CLIENT_ID`: Client ID for traQ
- `WEBHOOK_SECRET`: Webhook secret for traQ
- `WEBHOOK_ID`: Webhook ID for traQ
- `OS_CONTAINER`: Container name for SwiftStorage
- `OS_USERNAME`: Username for SwiftStorage
- `OS_PASSWORD`: Password for SwiftStorage
- `OS_TENANT_NAME`: Tenant name for SwiftStorage
- `OS_TENANT_ID`: Tenant ID for SwiftStorage
- `OS_AUTH_URL`: Auth URL for SwiftStorage

2. Enter the server for Jomon staging server and run the following command in the project root.

```sh
sudo docker pull ghcr.io/traptitech/jomon-v2:latest
sudo docker run -d -p 1323:1323 --env-file .env ghcr.io/traptitech/jomon-v2
```
