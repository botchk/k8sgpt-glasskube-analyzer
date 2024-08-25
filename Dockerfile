FROM golang:1.23-alpine3.19 AS builder

ENV CGO_ENABLED=0
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# TODO what are "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" 
# TODO what are -s -w -X flags
RUN go build -o /src/analyzer .

FROM gcr.io/distroless/static

WORKDIR /
COPY --from=builder /src/analyzer .
USER 65532:65532

ENTRYPOINT ["/analyzer"]