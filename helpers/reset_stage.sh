#!/bin/bash
. helpers/show_stage.sh
rm -fr terraform/providers/{shared,gitlab}/service-template/.terragrunt-cache
rm -fr terraform/providers/{shared,gitlab}/service-template/.terraform.lock.hcl

sudo rm -fr sf-data/services/test-service1*
sudo rm -fr terraform/providers/gitlab/services/test-service1*
sudo rm -fr terraform/providers/gitlab/services/test-service1*
sudo dr service-repository

# .terragrunt-cache
# delete test repo
if [ "$1" == "-purge" ]; then
  sudo rm -fr sf-data/repository/gitlab/franksrobins/cookie-cutter-api
fi
