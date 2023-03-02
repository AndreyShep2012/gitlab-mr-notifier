FROM golang:1.19-alpine

WORKDIR /app

COPY . ./

RUN go build -o gitlab-mr-notifier ./cmd/main.go

CMD [ "./gitlab-mr-notifier" ]

