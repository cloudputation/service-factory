#!/usr/bin/dumb-init /bin/sh
# Copyright (c) Cloudputation, Inc.

set -e

# # Initiate terraform
# for d in $(ls ./terraform);
# do
#   terraform -chdir="${d}" init
# done
  terraform -chdir="terraform/" init


# If the user is trying to run service-factory directly with some arguments,
# then pass them to service-factory.
# On alpine /bin/sh is busybox which supports the bashism below.
if [ "${1:0:1}" = '-' ]; then
	set -- /bin/service-factory "$@"
fi

# If user is trying to run service-factory with no arguments (daemon-mode),
# docker will run '/bin/sh -c /bin/${NAME}'. Check for the full command since
# running 'bin/sh' is a common pattern
if [ "$*" = '/bin/sh -c /bin/${NAME}' ]; then
	set -- /bin/service-factory
fi

# Matches VOLUME in the Dockerfile, for importing config files into image
SF_CONFIG_DIR=/service-factory/config

# Set the configuration directory
if [ "$1" = '/bin/service-factory' ]; then
	shift
	set -- /bin/service-factory agent #-config-dir="$SF_CONFIG_DIR" "$@"
fi

exec "$@"
