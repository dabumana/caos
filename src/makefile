APP=caos
VERSION=v0.0.0

build:
	mkdir bin
	mkdir bin/caos
	mkdir report
	mkdir report/${VERSION}
	cp .env ./bin/caos
	go build -o ./bin/caos/${APP}

clean:
	go clean
 	# Remove cache directories
	rm -rf 'bin'
	rm -rf 'log'
	rm -rf 'export'
	rm -rf 'report'
	rm -rf 'training'
	rm -rf '.cookies'
	rm -rf 'localCA.key'
	rm -rf 'localCA.crt'
	rm -rf 'localCA.pem'
	rm -rf 'rootCA.pem'
	rm -rf 'coverage.out'

coverage:
	go test ./test/... -coverprofile=report/${VERSION}/coverage.out -coverpkg ./...

env:
	openssl genrsa -out localCA.key 4096
	openssl req -x509 -new -nodes -key localCA.key -sha256 -days 1024 -out localCA.crt
	openssl x509 -in localCA.crt -out localCA.pem -outform PEM

run:
	./bin/caos/${APP}

test: coverage
	go test ./test/...

tidy:
	go mod tidy

vendor:
	go mod vendor
