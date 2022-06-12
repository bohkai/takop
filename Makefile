.PHONY:build
build:
	docker build -t takop .

.PHONY:run
run:
	docker run --platform=linux/arm64 takop