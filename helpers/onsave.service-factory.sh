# !/bin/bash

function auto_push {
  API_VERSION=$(cat API_VERSION)


  first_digit=$(echo ${API_VERSION} | cut -d"." -f1)
  second_digit=$(echo ${API_VERSION} | cut -d"." -f2)
  third_digit=$(echo ${API_VERSION} | cut -d"." -f3)

  let "third_digit++"

  new_API_VERSION="${first_digit}.${second_digit}.${third_digit}"

  echo current api version is: $API_VERSION
  echo $new_API_VERSION | tee  API_VERSION



  git add .
  git commit -m "auto commit"
  git push

}

grep "production = true" GIT_CONTROLS/auto_push && auto_push || echo "Auto push deactivated."

function sync_to_stage {
  go mod tidy
  # rsync -a -P ./* devops@tool:~/dev/service-factory/
  rsync -a -P $SAVED_FILE devops@tool:~/dev/service-factory/$SAVED_FILE
  rsync -a -P ./go.mod devops@tool:~/dev/service-factory/
  rsync -a -P ./go.sum devops@tool:~/dev/service-factory/
  rsync -a -P ./.air.toml devops@tool:~/dev/service-factory/
  rsync -a -P ./.release/ devops@tool:~/dev/service-factory/.release
  # ssh devops@tool "bash dev/service-factory/helpers/restart_server.sh"
  ssh devops@tool "sudo cp /home/devops/dev/service-factory/.release/defaults/config.hcl /etc/service-factory/config.hcl"
  cp .release/defaults/config.hcl /etc/service-factory/config.hcl
}

grep "staging = true" GIT_CONTROLS/auto_push && sync_to_stage || echo "Staging deactivated."
