ARG GO_VERSION=1.11
FROM golang:${GO_VERSION}-alpine AS builder
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
RUN apk add --no-cache ca-certificates git sqlite tzdata gcc musl-dev
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN GOOS=linux go build -ldflags="-w -s" -a -installsuffix 'static' -o /app

FROM alpine AS final
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /
COPY --from=builder /src/migrations /migrations
EXPOSE 8080
VOLUME /storage
USER nobody:nobody
ENTRYPOINT ["/app"]
