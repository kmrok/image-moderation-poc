FROM golang:latest AS builder
WORKDIR /go/src/image-moderation-poc
ENV GO111MODULE=on
COPY . .
RUN make build

FROM alpine:latest AS api
COPY --from=builder /go/src/image-moderation-poc/cmd/bin/app /usr/local/bin/app
CMD ["app"]
