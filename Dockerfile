FROM golang:latest

# Maintainer Info
ARG NAME=service-factory
ARG PRODUCT_VERSION

LABEL maintainer="Cloudputation"
LABEL version=$PRODUCT_VERSION


# ENV values
ENV NAME=$NAME
ENV VERSION=$PRODUCT_VERSION
ENV ROOTDIR="/service-factory"
ENV SF_CONFIG_FILE_PATH=${ROOTDIR}/config/config.hcl
ENV SF_LOG_DIRECTORY="/var/log"
ENV SF_DATA_DIRECTORY=${ROOTDIR}/sf-data
ENV SF_TERRAFORM_DIRECTORY=${ROOTDIR}/terraform

# BUILD values
ENV	TERRAGRUNT_PATH="/usr/bin/terragrunt"
ENV TERRAGRUNT_VERSION="0.53.8"
ENV TERRAGRUNT_URL=https://github.com/gruntwork-io/terragrunt/releases/download/v${TERRAGRUNT_VERSION}/terragrunt_linux_amd64

# Create service directories
RUN mkdir -p /service-factory/config
RUN mkdir -p ${SF_LOG_DIRECTORY}
RUN mkdir -p ${SF_DATA_DIRECTORY}/services
RUN mkdir -p ${SF_DATA_DIRECTORY}/repositories
RUN mkdir -p ${SF_TERRAFORM_DIRECTORY}



# Set the Current Working Directory inside the container
WORKDIR ${ROOTDIR}



RUN apt update
# RUN apt upgrade -y
RUN apt install -y\
	dumb-init\
	gnupg\
	software-properties-common

RUN wget -O- https://apt.releases.hashicorp.com/gpg |\
      gpg --dearmor |\
      tee /usr/share/keyrings/hashicorp-archive-keyring.gpg

RUN gpg --no-default-keyring\
      --keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg\
      --fingerprint

RUN echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg]\
      https://apt.releases.hashicorp.com $(lsb_release -cs) main" |\
      tee /etc/apt/sources.list.d/hashicorp.list

RUN apt update
RUN apt -y install terraform
RUN echo "TERRAFORM HAS BEEN INSTALLED! - \$(terraform version)"



RUN wget ${TERRAGRUNT_URL} -O ${TERRAGRUNT_PATH}
RUN chmod +x ${TERRAGRUNT_PATH}
RUN echo "TERRAGRUNT HAS BEEN INSTALLED! - $(terragrunt --version)"


COPY ./terraform/ ./terraform/
COPY ./API_VERSION ./API_VERSION


COPY ./build/service-factory /bin/service-factory
COPY ./.release/defaults/config.hcl /service-factory/config/config.hcl
COPY .release/docker/docker-entrypoint.sh /bin/docker-entrypoint.sh


# Set service user
RUN addgroup --system ${NAME} && adduser --system --ingroup ${NAME} ${NAME}

# Set service binary ownership
RUN chown -R ${NAME}:${NAME} /service-factory

# Set entry point permissions
RUN chmod +x /bin/docker-entrypoint.sh

# Expose port 48840
EXPOSE 48840

# Command to run the executable
ENTRYPOINT ["/bin/docker-entrypoint.sh"]
###

USER ${NAME}
CMD /bin/${NAME}
