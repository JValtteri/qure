# syntax=docker/dockerfile:1.4

# ------------------------------------------------------------------
# Build the front-end
# ------------------------------------------------------------------

FROM node:22-alpine AS frontend

WORKDIR /app
# copy package*.json separately to improve cache reusability
COPY ./client/package*.json ./
# Faster than 'npm install' when package-lock.json exists
RUN npm ci
COPY ./client ./
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
    linux/amd64)  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ;; \
    linux/arm64)  CGO_ENABLED=0 GOOS=linux GOARCH=arm64 ;; \
    linux/arm/v7) CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 ;; \
    *) echo "Unsupported TARGETPLATFORM: ${TARGETPLATFORM}" ;; \
  esac
RUN go build .

# ------------------------------------------------------------------
# Final scratch image
# ------------------------------------------------------------------

FROM scratch AS Final

WORKDIR /app/server

COPY --from=backend /app/server /app/server/server
COPY ./server/ATTRIBUTION /app/server/ATTRIBUTION

COPY --from=frontend /app/dist /app/client/dist

EXPOSE 8080
ENTRYPOINT ["/app/server/server"]
