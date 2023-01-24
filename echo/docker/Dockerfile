FROM golang:1.19 as builder
ARG ARCH
COPY /api /build/api
COPY /echo /build/echo
WORKDIR /build/echo
ENV GOOS=linux
ENV GOARCH=${ARCH}
ENV CGO_ENABLED=0
RUN go build -ldflags="-extldflags=-static" -o echo main.go

FROM scratch
COPY --from=builder /build/echo/echo /echo
ENTRYPOINT ["/echo"]
