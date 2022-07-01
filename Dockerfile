# Step 1: Modules caching
FROM golang:1.18.3-alpine3.15 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.18.3-alpine3.15 as builder

COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/api ./cmd/api

# Step 3: Final
FROM scratch
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/api /api
CMD ["/api", "-port=4000", "-env=development"]