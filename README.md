<div align="center">
  <h1>Jomon</h1>
  <p>
    <strong>An accounting support system</strong>
  </p>
  <p>
    <a href="https://apis.trap.jp/?urls.primaryName=Jomon#/">Jomon API</a>
  <p>
</div>

## Environment

### Testing

1. Run following command in the project root.

```shell script
make test
```

### Running

1. Run following command in the project root.

```shell script
make dev-up
```

Now, you can send http requests to `localhost:3000`.

You can also run following command in the project root.

```shell script
make up
```

The server is up without debugging logs.


## Staging

1.Enter the server for Jomon staging server and run following comand in the project root.

```shell script
sudo docker pull docker.pkg.github.com/traptitech/jomon/jomon:latest
sudo docker run -d -p 1323:1323 --env-file .env docker.pkg.github.com/traptitech/jomon/jomon
```

(At first,you need to set .env file)
