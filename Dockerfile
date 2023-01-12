FROM golang:1.15-alpine as builder

RUN apk update && apk add git make


WORKDIR /app
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download



COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goapp .



# for reduce image size

FROM scratch

COPY --from=builder /app/goapp /app

#COPY --from=builder /app/.env /goapp/app/.env

####### Requires this ##############

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/



ENTRYPOINT ["/app"]