docker_build:
	docker build -t 55nodes/controller .

docker_run:
	docker run -p 80:3000 -it 55nodes/controller:latest

docker_test:
	go test ./...
