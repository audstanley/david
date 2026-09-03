# Not tested yet
FROM golang:1.27.0-alpine3.24 AS build
WORKDIR $GOPATH/src/github.com/audstanley/david
COPY . .
RUN cd cmd/david && go build . && mv david /go/bin
RUN cd cmd/dcrypt && go build . && mv dcrypt /go/bin

FROM alpine:latest  
RUN addgroup -g 1000 david
RUN adduser -S -G david -u 1000 david
COPY --from=build /go/bin/dcrypt /usr/local/bin
COPY --from=build /go/bin/david /usr/local/bin
USER david
ENTRYPOINT ["/usr/local/bin/david"]
