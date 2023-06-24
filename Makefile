.PHONY:build
build:
	docker build -t takop .

.PHONY:docker-run
docker-run:
	docker run --name takop -it takop

.PHONY:run
run:
	go run ./