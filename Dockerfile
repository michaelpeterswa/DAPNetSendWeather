# -=-=-=-=-=-=- Compile Image -=-=-=-=-=-=-

FROM golang:1.17 AS stage-compile

WORKDIR /go/src/app
COPY . .

RUN go get -d -v . && CGO_ENABLED=0 GOOS=linux go build .

# -=-=-=-=-=-=- Final Image -=-=-=-=-=-=-

FROM alpine:3.17.2

WORKDIR /root/
COPY --from=stage-compile /go/src/app/dapnetsendweather ./

RUN apk --no-cache add ca-certificates=20161130-r0

ENTRYPOINT [ "./dapnetsendweather" ]  