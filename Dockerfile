#FROM scratch
#
##COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
##COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo/
#RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
#ADD ca-certificates.crt /etc/ssl/certs/
#ADD main /
#ADD .env /
#CMD ["/main"]

#FROM golang:alpine as build
## Redundant, current golang images already include ca-certificates
#RUN apk --no-cache add ca-certificates
#WORKDIR /go/src/app
#COPY . .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
#
#FROM scratch
## copy the ca-certificate.crt from the build stage
#COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=build /go/bin/app /app
#ENTRYPOINT ["/app"]
FROM golang:1.15-alpine as builder

RUN apk update && apk add git make

ARG db_name
ENV db_name=$db_name
ARG db_pass
ENV db_pass=$db_pass
ARG db_user
ENV db_user=$db_user
ARG db_type
ENV db_type=$db_type
ARG db_host
ENV db_host=$db_host
ARG db_port
ENV db_port=$db_port
ARG bot_token
ENV bot_token=$bot_token


WORKDIR /app

ADD .env /


ENV GO111MODULE=on



COPY go.mod .

COPY go.sum .

RUN go mod download



COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goapp .



# for reduce image size

FROM scratch

COPY --from=builder /app/goapp /app

COPY --from=builder /app/.env /goapp/app/.env

####### Requires this ##############

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/



ENTRYPOINT ["/app"]