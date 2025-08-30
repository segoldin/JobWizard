# Dockerfile to build the job_wizard CLI executable
# Sally Goldin, 2025-08-20
# build stage
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git
ENV GOPATH=/go

WORKDIR /go/src

# Download Go modules
RUN mkdir JobWizard
COPY go.mod go.sum JobWizard
WORKDIR /go/src/JobWizard
RUN mkdir job_wizard
COPY JobWizard/job_wizard/* ./job_wizard
RUN go mod download

WORKDIR /go/src/JobWizard/job_wizard

LABEL author="SallyGoldin"

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o job_wizard .

# Final stage: Create a minimal image with the built binary
FROM ubuntu:22.04

# Set environment variables to prevent interactive prompts
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=UTC

# Create a non-root user with specific UID and GID
RUN groupadd -g 1000 appgroup && \
    useradd -u 1000 -g appgroup -s /bin/bash -m appuser

# Set the working directory
WORKDIR /app

# Valid options for ENV:
# - prod (default)
# - dev
# - staging
ARG ENV=prod

COPY --from=builder /go/src/job_wizard .
COPY --from=builder /go/src/.env .

# Change ownership of the app directory to the non-root user
RUN chown -R appuser:appgroup /app

# Create directories that the application might need to write to
RUN mkdir -p /app/logs /tmp && \
    chown -R appuser:appgroup /app/logs

# Create a directory to store the templates
RUN mkdir -p /app/templates && \
    chown -R appuser:appgroup /app/templates

# Switch to the non-root user
USER 1000:1000

# Set the environment variable for the application
ENV ENV=${ENV}

CMD [ "./job_wizard" ]
