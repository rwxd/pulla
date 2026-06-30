FROM docker.io/golang:1.23-alpine as builder

COPY . /build
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -a -o pulla .

# generate clean, final image for end users
FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

RUN apk add --no-cache ca-certificates git

COPY --from=builder /build/pulla .

ENTRYPOINT [ "./pulla" ]
CMD []
