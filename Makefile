SIMULATOR_BINARY_NAME=node

build-simulator:
	GOARCH=amd64 GOOS=darwin go build -o build/${SIMULATOR_BINARY_NAME}-amd64-darwin cmd/simulator/main.go
	GOARCH=amd64 GOOS=linux go build -o build/${SIMULATOR_BINARY_NAME}-amd64-linux cmd/simulator/main.go
	GOARCH=amd64 GOOS=windows go build -o build/${SIMULATOR_BINARY_NAME}-amd64-windows cmd/simulator/main.go
	GOARCH=arm64 GOOS=darwin go build -o build/${SIMULATOR_BINARY_NAME}-arm64-darwin cmd/simulator/main.go

clean:
	go clean
	rm ${SIMULATOR_BINARY_NAME}-amd64-darwin
	rm ${SIMULATOR_BINARY_NAME}-amd64-linux
	rm ${SIMULATOR_BINARY_NAME}-amd64-windows
	rm ${SIMULATOR_BINARY_NAME}-arm64-darwin

deps:
	go mod download

test:
	go test ./...