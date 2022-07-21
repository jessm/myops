##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY myops/*.go ./
RUN go build -o /myops

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

COPY --from=build /myops /myops

ENTRYPOINT ["/myops"]
