FROM golang:1.21.4 as build
WORKDIR /var/app
COPY go.mod go.sum ./
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /exception-control cmd/exception-control/main.go

FROM scratch
COPY --from=build exception-control /bin/exception-control
EXPOSE 80
ENTRYPOINT ["/bin/exception-control"]