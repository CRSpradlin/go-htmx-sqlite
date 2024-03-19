FROM golang:1.22

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]
# Build Docker image by running command: docker build -t go-htmx-sqlite:0.0.1 .
# Run Docker container by running command: docker run -it -p 80:8080 --name running-go-htmx-sqlite  go-htmx-sqlite:0.0.1