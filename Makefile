build:
	go build -o gitlab-mr-notifier ./cmd/main.go

test: 
	go test -v -race -count=1 -cover -coverprofile="./coverage.out" ./...
	