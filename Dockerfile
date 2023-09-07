FROM golang:latest AS build

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GO111MODULE=on

WORKDIR /mybili

ADD . /mybili

RUN go build -tags netgo -o mybili_server



FROM alpine

ENV GIN_MODE="release"
ENV PORT=3000

WORKDIR /mybili

COPY --from=build /mybili/mybili_server /mybili/mybili_server
COPY --from=build /mybili/.env /mybili/.env

RUN chmod +x /mybili/mybili_server

ENTRYPOINT ["./mybili_server"]