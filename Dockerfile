FROM golang:1.21.6 AS builder

WORKDIR /go/bin

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .


RUN go build -o /go/bin/user  *.go


FROM ubuntu

COPY --from=builder /go/bin/user /go/bin/user
COPY --from=builder /go/bin/config/ /go/bin/config/


WORKDIR /go/bin/

EXPOSE 5300

ENTRYPOINT ["/go/bin/user"]
