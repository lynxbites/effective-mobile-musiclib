FROM golang:alpine

RUN mkdir /app
WORKDIR /app

COPY . /app
COPY .env /app

RUN go get -d ./...

RUN go build -o app ./cmd/server/main.go 

CMD [ "./app" ]