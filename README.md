## Deploy

````
gcloud app deploy
````

## Docker

````
docker ps -a
docker build -t taplink .
docker run -p 8080:8080 --env-file ./.env taplink
````

docker rm -vf $(docker ps -aq)

GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .
