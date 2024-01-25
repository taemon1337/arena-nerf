FROM golang:1.21.5 as builder

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o /serf-cluster

FROM gcr.io/distroless/static-debian12

COPY --from=builder /serf-cluster /

ENTRYPOINT ["/serf-cluster"]
