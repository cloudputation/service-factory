# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
ARG NAME=service-factory
ARG PRODUCT_VERSION

LABEL maintainer="Cloudputation"
LABEL version=$PRODUCT_VERSION

# Set ARGs as ENV so that they can be used in ENTRYPOINT/CMD
ENV NAME=$NAME
ENV VERSION=$PRODUCT_VERSION

# Set the Current Working Directory inside the container
WORKDIR /service-factory

RUN apt update
RUN apt update && apt install -y gnupg software-properties-common
RUN wget -O- https://apt.releases.hashicorp.com/gpg | \
      gpg --dearmor | \
      tee /usr/share/keyrings/hashicorp-archive-keyring.gpg

RUN gpg --no-default-keyring \
      --keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg \
      --fingerprint

RUN echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] \
      https://apt.releases.hashicorp.com $(lsb_release -cs) main" | \
      tee /etc/apt/sources.list.d/hashicorp.list

RUN apt update
RUN apt -y install terraform
RUN echo "TERRAFORM HAS BEEN INSTALLED! - $(terraform version)"

RUN addgroup --system ${NAME} && adduser --system --ingroup ${NAME} ${NAME}

COPY ./terraform/ ./terraform/
COPY ./API_VERSION ./API_VERSION

# Copy the source from the current directory to the Working Directory inside the container
COPY ./build/service-factory /bin/service-factory
COPY ./.release/defaults/config.hcl /service-factory/config/config.hcl
COPY .release/docker/docker-entrypoint.sh /bin/docker-entrypoint.sh

# Add permissions to entry script
RUN chmod +x /bin/docker-entrypoint.sh

RUN mkdir -p /service-factory/config \
	&& chown -R ${NAME}:${NAME} /service-factory

# Expose port 48840 to the outside
EXPOSE 48840

# Command to run the executable
ENTRYPOINT ["/bin/docker-entrypoint.sh"]
###

USER ${NAME}
CMD /bin/${NAME}
