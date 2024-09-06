#NOTE; Keep build in sync with .gorelease.yml
FROM golang:1.23-alpine3.19 AS builder

ARG VERSION
ARG COMMIT
ARG DATE

ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /src/analyzer -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" .

FROM gcr.io/distroless/static

LABEL org.opencontainers.image.created="${DATE}" \
    org.opencontainers.image.revision="${COMMIT}" \
    org.opencontainers.image.version="${VERSION}" \
    org.opencontainers.image.title="k8sgpt glasskube analyzer" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.source="https://github.com/botchk/k8sgpt-glasskube-analyzer/" \
    org.opencontainers.image.authors="https://github.com/botchk/"

WORKDIR /
COPY --from=builder /src/analyzer .
USER 65532:65532

EXPOSE 8085 

ENTRYPOINT ["/analyzer"]