FROM alpine:3


ARG NAME=service-factory
ARG SERVICE_USERNAME=sf
ARG PRODUCT_VERSION

ENV NAME=$NAME
ENV VERSION=$PRODUCT_VERSION
ENV ROOTDIR="/sf"
ENV SF_CONFIG_FILE_PATH=${ROOTDIR}/config/config.hcl
ENV SF_LOG_DIRECTORY=${ROOTDIR}/log
ENV SF_DATA_DIRECTORY=${ROOTDIR}/sf-data
ENV TERRAFORM_PATH="/usr/local/bin/terraform"
ENV TERRAGRUNT_PATH="/usr/local/bin/terragrunt"

WORKDIR ${ROOTDIR}

# Install runtime dependencies
RUN apk add --no-cache dumb-init git

# Create service directories
RUN mkdir -p /sf/config \
    && mkdir -p ${SF_LOG_DIRECTORY} \
    && mkdir -p ${SF_DATA_DIRECTORY}/services \
    && mkdir -p ${SF_DATA_DIRECTORY}/repositories

# Set service user
RUN addgroup -g 991 ${SERVICE_USERNAME} \
    && adduser -D -u 991 -G ${SERVICE_USERNAME} ${SERVICE_USERNAME}

# Copy artifacts from builder
COPY ./API_VERSION ./API_VERSION
COPY ./artifacts/terraform ${TERRAFORM_PATH}
COPY ./artifacts/terragrunt ${TERRAGRUNT_PATH}
COPY ./build/service-factory /bin/service-factory
COPY ./.release/defaults/config.hcl /sf/config/config.hcl
COPY .release/docker/docker-entrypoint.sh /bin/docker-entrypoint.sh

# Set permissions
RUN chown -R ${SERVICE_USERNAME}:${SERVICE_USERNAME} ${ROOTDIR} \
    && chmod +x /bin/docker-entrypoint.sh \
    && chmod +x ${TERRAFORM_PATH} \
    && chmod +x ${TERRAGRUNT_PATH}

# Expose port 48840
EXPOSE 48840

# Set user
USER ${SERVICE_USERNAME}

# Entrypoint to run the executable
ENTRYPOINT ["/bin/docker-entrypoint.sh"]
CMD ["/bin/service-factory"]
