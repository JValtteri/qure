# syntax=docker/dockerfile:1.4

# ------------------------------------------------------------------
# Build the front-end
# ------------------------------------------------------------------

FROM node:22-alpine AS frontend

WORKDIR /app
COPY ./client ./
RUN npm install
# creates /app/dist
RUN npx vite build

# ------------------------------------------------------------------
# Build the Go back-end
# ------------------------------------------------------------------

FROM golang:alpine AS backend

WORKDIR /app
COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/. .

# Target-specific Go environment variables
ARG TARGETPLATFORM
ENV CGO_ENABLED=0

# Check the variable is valid
RUN echo "TARGETPLATFORM=${TARGETPLATFORM}"

# Compile the Go binary for the requested platform
RUN case "${TARGETPLATFORM}" in \
    linux/amd64)  GOOS=linux GOARCH=amd64 ;; \
    linux/arm64)  GOOS=linux GOARCH=arm64 ;; \
    linux/arm/v7) GOOS=linux GOARCH=arm GOARM=7 ;; \
    *) echo "Unsupported TARGETPLATFORM: ${TARGETPLATFORM}" ;; \
  esac
RUN go build .

# ------------------------------------------------------------------
# Final image - Alpine + back-end + static front-end
# ------------------------------------------------------------------

FROM alpine:latest

WORKDIR /app/server

RUN apk add --no-cache ca-certificates

COPY --from=backend /app/server /app/server/server
COPY --from=frontend /app/dist /app/client/dist

EXPOSE 8080
ENTRYPOINT ["/app/server/server"]
