FROM golang:1.19 as builder
ARG ARCH
COPY /api /build/api
COPY /server /build/server
WORKDIR /build/server
ENV GOOS=linux
ENV GOARCH=${ARCH}
RUN go build -ldflags="-extldflags=-static" -o server server.go

FROM scratch
COPY --from=builder /build/server/server /server
ENTRYPOINT ["/server"]
