run-api:
	CGO_ENABLED=0 go run cmd/api/main.go

test:
	CGO_ENABLED=0 GOTOOLCHAIN=go1.25.0+auto go test -race -short -count=1 ./... -gcflags=all=-l
	
test-cover:
	CGO_ENABLED=0 GOTOOLCHAIN=go1.25.0+auto go test -cover -race -short -count=1 -coverprofile=coverage.out ./... -gcflags=all=-l
	go tool cover -html=coverage.out

mock-gen:
	@./generate-mock.sh internal/services
	@./generate-mock.sh internal/repositories
	@./generate-mock.sh internal/pkg/file
	@./generate-mock.sh internal/pkg/event
	@./generate-mock.sh internal/pkg/queue
