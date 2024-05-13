build:
	go build -o gitlab-mr-notifier ./cmd/main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gitlab-mr-notifier ./cmd/main.go

aws-zip: build-linux
	zip -jrm gitlab-mr-notifier.zip gitlab-mr-notifier

test: 
	go test -v -race -count=1 -cover -coverprofile="./coverage.out" ./...

check-escape:
	go build -gcflags '-m' ./...
