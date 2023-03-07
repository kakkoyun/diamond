FROM golang:1.20-alpine3.17 as builder
RUN mkdir /.cache && chown nobody:nogroup /.cache && touch -t 202101010000.00 /.cache

ARG VERSION
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download -modcacherw

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY --chown=nobody:nogroup ./main.go ./main.go

RUN mkdir bin
RUN go build -trimpath -a -o ./bin/diamond ./main.go

FROM alpine:3.17

USER nobody

COPY --chown=0:0 --from=builder /app/bin/diamond /bin/diamond

CMD ["/bin/diamond"]
