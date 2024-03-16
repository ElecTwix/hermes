# Container image that runs your code
FROM golang:1.22

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify


COPY . .
RUN ls -la
RUN go build -v -o /usr/local/bin/app /usr/src/app/cmd/hermes/hermes.go 
RUN ls -la /usr/local/bin/app


ENTRYPOINT ["./entrypoint.sh"]
