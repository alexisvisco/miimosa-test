generate-proto:
	protoc --go_out=./pkg --go_opt=paths=import --go-grpc_out=./pkg/ --go-grpc_opt=paths=import proto/sessions.proto

tests:
	go test ./...

docker-build:
	docker build -t miimosa-test .

docker-run: docker-build
	docker run -p 3123:3123 miimosa-test

run:
	go run cmd/miimosa-test/main.go
