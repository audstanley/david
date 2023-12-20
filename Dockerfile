# Not tested yet
FROM golang:1.21.3-alpine AS build
WORKDIR $GOPATH/src/github.com/audstanley/david
COPY . .
RUN cd cmd/david && go build . && mv david ~/go/bin
RUN cd cmd/bcpt && go build . && mv bcpt ~/go/bin

FROM alpine:latest  
RUN addgroup -g 1000 david
RUN adduser -S -G david -u 1000 david
COPY --from=build /go/bin/bcpt /usr/local/bin
COPY --from=build /go/bin/david /usr/local/bin
USER david
ENTRYPOINT ["/home/david/go/bin/david"]
