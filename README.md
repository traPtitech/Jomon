# Jomon
This is an accounting support system ,"Jomon".

[This](https://traPtitech.github.io/Jomon/dist/index.html) is a page about Jomon API.

## Environment
### Server test
1. Run following command in the project root.
```shell script
docker-compose -f server-test.yml run --rm jomon-server
```
### Client
1. Run following command in the project root.
```shell script
docker-compose -f mock-for-client.yml up
```
Now you can access to `http://localhost:3000` for Jomon client page.
And you can access to `http://localhost:1323` for Jomon mock server using `swagger.yaml`.

## Staging
1.Enter the server for Jomon staging server and run following comand in the project root.
```shell script
sudo docker build -t jomon_staging .
sudo docker run -p 8080:8080 --env-file .env jomon_staging
```
(At first,you need to clone (or pull) repositry and set .env)
