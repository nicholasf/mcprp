.PHONY: build deploy

build:
	go build -o service

deploy: build
	scp service 134.199.175.220
