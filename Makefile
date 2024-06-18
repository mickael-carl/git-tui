.PHONY: install lint

install:
	go install

lint:
	golangci-lint run
