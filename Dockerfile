FROM golang:1.21.3-alpine AS build
WORKDIR $GOPATH/src/github.com/audstanley/david
COPY . .
RUN go build -o /go/bin/david cmd/david/main.go
RUN go build -o /go/bin/bcpt cmd/bcpt/main.go

FROM alpine:latest  
RUN addgroup -g 1000 david
RUN adduser -S -G david -u 1000 david
COPY --from=build /go/bin/bcpt /usr/local/bin
COPY --from=build /go/bin/david /usr/local/bin
USER david
ENTRYPOINT ["/usr/local/bin/david"]
