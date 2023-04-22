APP=caos
ARCH=amd64
VERSION=v.0.0.0

build: clean
	GOARCH=${ARCH} go build -o ./bin/${APP}-${VERSION}

clean:
	go clean
 	# Remove cache directories
	rm -rf 'bin'
	rm -rf 'log'
	rm -rf 'export'
	rm -rf 'training'

run: build
	./bin/${APP}-${VERSION}

test:
	go test -cover ./...

update:
	go mod tidy