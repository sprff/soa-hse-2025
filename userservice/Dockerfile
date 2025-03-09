FROM golang:1.24-alpine3.20 AS builder

ENV SERVICE_NAME=userservice
ENV SHARED_NAME=shared

# Prepare depencies
COPY ${SHARED_NAME}/source/go.mod /app/${SHARED_NAME}/source/go.mod
# COPY ${SHARED_NAME}/source/go.sum /app/${SHARED_NAME}/source/go.sum
# WORKDIR /app/${SHARED_NAME}/source
# RUN go mod download

COPY ${SERVICE_NAME}/source/go.mod /app/${SERVICE_NAME}/source/go.mod
COPY ${SERVICE_NAME}/source/go.sum /app/${SERVICE_NAME}/source/go.sum
WORKDIR /app/${SERVICE_NAME}/source
RUN go mod download

# Copy sources
COPY ${SERVICE_NAME}/source /app/${SERVICE_NAME}/source
COPY ${SHARED_NAME}/source /app/${SHARED_NAME}/source

#Enable caching and build
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /app/bin/app /app/${SERVICE_NAME}/source/cmd/main


FROM alpine:3.14 AS runner
WORKDIR /app
COPY --from=builder /app/bin/app /app/bin/app
# COPY config ./config
EXPOSE 8080
WORKDIR /app


ENTRYPOINT ["/app/bin/app"]