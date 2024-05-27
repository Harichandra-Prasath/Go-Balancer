run: build
	@./bin/balancer

build:
	@ go build -o "bin/balancer"