# Create base builder image
FROM golang:1.19.6-alpine3.17 AS builder
WORKDIR /go/src/github.com/chain4travel/camino-signavault
RUN apk add --no-cache alpine-sdk bash git make gcc musl-dev linux-headers git ca-certificates g++ libstdc++

# Build app
COPY . .
RUN if [ -d "./vendor" ];then export MOD=vendor; else export MOD=mod; fi && \
    GOOS=linux GOARCH=amd64 go build -mod=$MOD -o /opt/camino-signavault ./cmd/camino-signavault/*.go

# Create final image
FROM alpine:3.17 as execution
RUN apk add --no-cache libstdc++
VOLUME /var/log/camino-signavault
WORKDIR /opt/camino-signavault

# Copy in and wire up build artifacts
COPY --from=builder /opt/camino-signavault /opt/camino-signavault/camino-signavault
COPY --from=builder /go/src/github.com/chain4travel/camino-signavault/config.yml /opt/camino-signavault/config.yml
COPY --from=builder /go/src/github.com/chain4travel/camino-signavault/db/migrations /opt/camino-signavault/migrations
ENTRYPOINT ["/opt/camino-signavault/camino-signavault"]
