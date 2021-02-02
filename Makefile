build:
	go build
run:
	go run .
test:
	go test -v -cover ./...
clean:
	rm dispatch-simulation
docker:
	docker build -t hub.docker.com/fguy/dispatch-simulation .
mock:
	mockgen -destination=mocks/go.uber.org/fx/lifecycle.go go.uber.org/fx Lifecycle
	go generate ./...
