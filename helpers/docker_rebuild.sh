# !/bin/bash

SERVICE_PORT="48840"
sudo docker stop sf
sudo docker rm sf
make docker-build
sudo docker run -d -p${SERVICE_PORT}:${SERVICE_PORT} --name sf service-factory
sudo docker logs --follow sf
