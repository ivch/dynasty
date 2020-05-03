################################
# STEP 1 build executable binary
################################

FROM golang:1.13.8-stretch as builder

RUN apt update && apt install -y make gcc musl-dev git ca-certificates && update-ca-certificates && mkdir -p /app

WORKDIR /app

ADD go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go test -mod=vendor -cover -race -v ./...
RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -a -o app

############################
# STEP 2 build a small image
############################

FROM scratch

ADD zoneinfo.tar.gz /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cmd/app /app
COPY /_ui /_ui

ENTRYPOINT ["./app"]
