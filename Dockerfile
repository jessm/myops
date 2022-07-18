##
## Build
##
FROM golang:1.16-buster AS build

WORKDIR /app

COPY cmd/ping/main.go ./
COPY cmd/ping/go.mod ./

RUN go build -o /ping

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /ping /ping

EXPOSE 8123

ENTRYPOINT ["/ping"]
