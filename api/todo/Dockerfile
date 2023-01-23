FROM golang:1.10 AS build
WORKDIR /go/src
COPY . ./.
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o todo .

FROM scratch AS runtime
COPY --from=build /go/src/todo ./
EXPOSE 8080/tcp
ENTRYPOINT ["./todo"]
