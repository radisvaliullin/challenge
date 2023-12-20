BINARY_NAME=store

build:
	GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin ./cmd/store/main.go
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux ./cmd/store/main.go

test:
	go test ./...

test-cover:
	go test ./... -cover

docker-build:
	docker build -f ./deploy/Dockerfile -t store .

docker-run:
	docker run --name=store -p 8080:8080 -d store

docker-down:
	docker stop store
	docker rm store

