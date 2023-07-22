## Deploy

````
gcloud app deploy
````

## Docker

````
docker ps -a
docker build -t taplink .
docker run --env-file ./.env taplink
````