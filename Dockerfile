FROM docker.io/golang:1.20-alpine as builder

COPY . /build
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -a -o pulla .

# generate clean, final image for end users
FROM alpine:3.17.3

RUN apk add --no-cache ca-certificates git

COPY --from=builder /build/pulla .

ENTRYPOINT [ "./pulla" ]
CMD []
