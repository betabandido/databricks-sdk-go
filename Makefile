all: test

test: get-deps generate
	go test ./...

get-deps:
	go get -u github.com/golang/glog
	go get -u github.com/stretchr/testify
	go get -u gopkg.in/h2non/gock.v1

generate:
	go generate
