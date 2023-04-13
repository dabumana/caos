APP=''
ARCH=''
ENTRYPOINT=main.go
VERSION=''

build: clean
	GOARCH=${ARCH} GOOS=darwin go build -o ./bin/${APP}-${VERSION} ${ENDPOINT}
	GOARCH=${ARCH} GOOS=linux go build -o ./bin/${APP}-${VERSION} ${ENDPOINT}
	GOARCH=${ARCH} GOOS=windows go build -o ./bin/${APP}-${VERSION} ${ENDPOINT}

clean:
	go clean
 	# Remove cache directories
	rm -rf 'bin'
	rm -rf 'log'
	rm -rf 'export'
	rm -rf 'training'

coverage:
	go test --cover ./...

run: build
	./bin/${APP}-${VERSION}

test:
	go test ./...

update:
	go mod tidy