FROM golang:1.15 AS builder

RUN mkdir $GOPATH/src/golang
WORKDIR $GOPATH/src/golang

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /bin/srv ./server/main.go

#FROM golang:1.15
#
#COPY --from=builder /srv /bin/srv

ENTRYPOINT ["/bin/srv"]