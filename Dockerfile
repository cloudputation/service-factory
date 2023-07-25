# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Your Name <your-email@domain.com>"

# Set the Current Working Directory inside the container
WORKDIR /app


RUN apt update
RUN apt update &&  apt install -y gnupg software-properties-common
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


# Copy go mod and sum files
COPY go.mod go.sum ./
COPY ./terraform/ ./terraform/
COPY ./API_VERSION/ ./API_VERSION

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ./main.go .

# Build the Go app
RUN go build -o sf .
RUN cp sf /usr/bin/sf


RUN terraform -chdir="terraform" init
# Expose port 48840 to the outside
EXPOSE 48840

# Command to run the executable
CMD ["/usr/bin/sf"]
