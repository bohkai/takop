.PHONY:build
build:
	docker build -t takop .

.PHONY:run
run:
	go run ./