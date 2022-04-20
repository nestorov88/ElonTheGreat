FROM golang:1.18-alpine3.15

EXPOSE 9000

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .


RUN go build cmd/tweet/main.go
RUN mv main /usr/local/bin/

#RUN go build cmd/importer/main.go
#RUN mv main /usr/local/bin/

CMD ["main"]
