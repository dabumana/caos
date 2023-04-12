APP=caos
ENTRYPOINT=main.go

build:
	GOARCH=amd64 GOOS=darwin go build -o ./src/bin/${APP}-darwin ${ENDPOINT}
	GOARCH=amd64 GOOS=linux go build -o ./src/bin/${APP}-linux ${ENDPOINT}
	GOARCH=amd64 GOOS=windows go build -o ./src/bin/${APP}-windows ${ENDPOINT}

clean:
	go clean go.sum
 	# Remove cache directories
	rm -rf 'bin'
	rm -rf 'log'
	rm -rf 'export'
	rm -rf 'training'

coverage:
	go test --cover ./...

run: build
	./${APP}

test:
	go test ./...

update:
	go mod tiny 