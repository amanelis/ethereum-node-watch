docker_build:
	docker build -t 55nodes/reporter .

docker_test:
	go test ./...
