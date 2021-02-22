<div align="center">
  <h1>Jomon</h1>
  <p>
    <strong>An accounting support system</strong>
  </p>
  <p>
    <a href="https://traPtitech.github.io/Jomon/dist/index.html">This</a> is a page about Jomon API.
  <p>
</div>

## Environment

### Server test

1. Run following command in the project root.

```shell script
make server-test
```

### Client

1. Run following command in the project root.

```shell script
make client
```

Now you can access to `http://localhost:3000` for Jomon client page.
And you can access to `http://localhost:1323` for Jomon mock server using `swagger.yaml`.

## Staging

1.Enter the server for Jomon staging server and run following comand in the project root.

```shell script
sudo docker pull docker.pkg.github.com/traptitech/jomon/jomon:latest
sudo docker run -d -p 1323:1323 --env-file .env docker.pkg.github.com/traptitech/jomon/jomon
```

(At first,you need to set .env file)
