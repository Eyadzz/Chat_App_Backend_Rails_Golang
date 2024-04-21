FROM golang:1.22-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY . ./
RUN go mod tidy
RUN go build -o /main
CMD ["/main"]