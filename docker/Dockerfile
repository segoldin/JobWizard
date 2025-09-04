# Dockerfile to build the job_wizard CLI executable image
# Sally Goldin, created 2025-08-20
# This version clones the entire repo because we were having problems with checkout

# build stage
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git
ENV GOPATH=/go

WORKDIR /app

# easy way to get the source files
RUN git clone https://github.com/segoldin/JobWizard.git

WORKDIR /app/JobWizard/job_wizard

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

# Install necessary packages
RUN apt-get update && apt-get install -y ca-certificates tzdata sqlite3 && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the built binary
COPY --from=builder /app/JobWizard/job_wizard/job_wizard .

# Copy the .env file
COPY --from=builder /app/JobWizard/job_wizard/.env_jobwizard .

# Copy the database
RUN mkdir database
WORKDIR /app/database
COPY --from=builder /app/JobWizard/job_wizard/database/jobwizard_db .

# Change ownership of the app directory to the non-root user
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER 1000:1000

# Set the environment variable for the application
ENV ENV=${ENV}

CMD [ "./job_wizard" ]
