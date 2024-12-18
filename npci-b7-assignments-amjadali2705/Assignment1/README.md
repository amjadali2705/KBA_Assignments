# Dockerfile


Command to build the docker image

```
docker build -t hellodockerimage .
```

Command to run the container

```
docker run hellodockerimage
```

Command to list the image

```
docker images
```

Command to list the container

```
docker ps -a
```

Command to see the output mesaage

```
docker run hellodockerimage
```

# Docker compose

Command to use Docker Compose to bring up the containers 

```
docker compose up -d
```

Output from container1

```
docker logs hellodockercontainer1
```

Output from container1

```
docker logs hellodockercontainer2
```