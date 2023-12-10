# !/bin/bash

SERVICE_PORT="48840"
HOME_SF_DATA="/opt/sf/sf-data"
DOCKER_SF_DATA="/sf/sf-data"

sudo docker stop sf || true
sudo docker rm sf || true
mkdir -p ./artifacts
cp ${HOME}/dev/artifacts/sf_artifacts/* ./artifacts

sudo make docker-build

if [ $? -eq 1 ]; then
    echo "make docker-build failed with exit code 1"
    exit 1
fi


for c in $(sudo docker ps -a | grep -i "exited" | awk '{print $1}')
do
  sudo docker rm $c
done

for i in $(sudo docker images | grep "none" | awk '{print $3}')
do
  sudo docker rmi $i
done

# set permissions on sf-data
# sudo mkdir -p /opt/sf-data
sudo chown -R 991:991 $HOME_SF_DATA
# sudo find $HOME_SF_DATA -type d -exec chmod 755 {} +
# sudo find $HOME_SF_DATA -type f -exec chmod 766 {} +

sudo docker run -d \
  -p ${SERVICE_PORT}:${SERVICE_PORT} \
  -v ${HOME_SF_DATA}:${DOCKER_SF_DATA} \
  --name sf service-factory

sudo docker logs --follow sf
