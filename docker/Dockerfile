# Stage 1: Build image
FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY ./src/ .
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/rabbit_travels ./main.go
# Stage 2: Deployment image
FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /var/app
RUN chmod -R 775 /var/app
COPY --from=build /go/src/app/bin /var/app
EXPOSE 1323
ENTRYPOINT /var/app/rabbit_travels