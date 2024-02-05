#!/bin/bash
PID=$(lsof -i :4884 | grep main | awk '{print $2}')
sudo kill $PID

rm -fr terraform/providers/gitlab/serviec-template/.terraform-cache
rm -fr terraform/providers/gitlab/serviec-template/.terraform.lock.hcl

sudo cp .release/defaults/config.hcl /etc/service-factory/config.hcl

go run main.go agent
