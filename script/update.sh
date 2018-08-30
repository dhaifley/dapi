#!/bin/sh
docker stack rm dapi
sleep 2s
docker pull dhaifley/dapi:latest
sleep 2s
docker system prune --volumes -f
sleep 2s
docker stack deploy -c /home/dhaifley/dapi/docker-compose.yml dapi
