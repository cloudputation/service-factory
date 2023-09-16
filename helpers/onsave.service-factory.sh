# #!/bin/bash
# #
# # function auto_push {
# #   API_VERSION=$(cat API_VERSION)
# #
# #
# #   first_digit=$(echo ${API_VERSION} | cut -d"." -f1)
# #   second_digit=$(echo ${API_VERSION} | cut -d"." -f2)
# #   third_digit=$(echo ${API_VERSION} | cut -d"." -f3)
# #
# #   let "third_digit++"
# #
# #   new_API_VERSION="${first_digit}.${second_digit}.${third_digit}"
# #
# #   echo current api version is: $API_VERSION
# #   echo $new_API_VERSION | tee  API_VERSION
# #
# #
# #
# #   git add .
# #   git commit -m "auto commit"
# #   git push
# #
# # }
# # grep "production = true" GIT_CONTROLS/auto_push && auto_push || echo "Auto push deactivated."
#
# function sync_to_stage {
#   rsync -a -P ./* devops@tower2:~/dev/service-factory/
# }
#
# grep "staging = true" GIT_CONTROLS/auto_push && sync_to_stage || echo "Staging deactivated."
