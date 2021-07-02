<div align="center">
  <h1>Jomon</h1>
  <p>
    <strong>An accounting support system</strong>
  </p>
  <p>
    <a href="https://apis.trap.jp/?urls.primaryName=Jomon%20v2%20API">Jomon API</a>&emsp;<a href="https://github.com/traPtitech/Jomon-UI">Jomon-UI</a>
  </p>
  <p>
    <a href="https://github.com/traPtitech/Jomon/actions/workflows/image.yml"><img src="https://github.com/traPtitech/Jomon/actions/workflows/image.yml/badge.svg"></a>
    <a href="https://github.com/traPtitech/Jomon/actions/workflows/go.yml"><img src="https://github.com/traPtitech/Jomon/actions/workflows/go.yml/badge.svg"></a>
    <a href="https://codecov.io/gh/traPtitech/Jomon"><img src="https://codecov.io/gh/traPtitech/Jomon/branch/v2/graph/badge.svg?token=kPPuCDXudP"></a>
  </p>
</div>

## Environment

### Testing

1. Make server running.

2. Run following command in the project root.
```shell script
MARIADB_HOSTNAME=localhost go test -v -cover -race ./...
```

### Running

1. Run following command in the project root.

```shell script
make up
```

Now, you can send http requests to `localhost:3000`.

2. Run following command in the project root when making server down.

```shell script
make down
```

## Staging

1.Enter the server for Jomon staging server and run following comand in the project root.

```shell script
sudo docker pull docker.pkg.github.com/traptitech/jomon/jomon:latest
sudo docker run -d -p 1323:1323 --env-file .env docker.pkg.github.com/traptitech/jomon/jomon
```

(At first,you need to set .env file)
