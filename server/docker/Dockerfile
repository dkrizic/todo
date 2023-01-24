FROM golang:1.19 as builder
ARG ARCH
COPY /api /build/api
COPY /server /build/server
WORKDIR /build/server
ENV GOOS=linux
ENV GOARCH=${ARCH}
ENV CGO_ENABLED=0
RUN go build -ldflags="-extldflags=-static" -o server main.go

FROM scratch
COPY --from=builder /build/server/server /server
ENTRYPOINT ["/server"]
CMD ["-v=2","serve","memory"]
