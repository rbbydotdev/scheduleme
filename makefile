.PHONY: build run

build:
		docker build -t go-docker-scheduleme .

run:
	 	docker run --rm go-docker-scheduleme

