APP=caos
VERSION=v.0.0.0
PKG_LIST=$(go list ./... | grep -v /vendor/ | tr '\n' ' ')

build:
	mkdir bin
	mkdir bin/caos
	cp .env ./bin/caos
	cp -r template ./bin/caos
	go build -o ./bin/caos/${APP}


clean:
	go clean
 	# Remove cache directories
	rm -rf 'bin'
	rm -rf 'log'
	rm -rf 'export'
	rm -rf 'training'
	rm -rf '.cookies'
	rm -rf 'localCA.key'
	rm -rf 'localCA.crt'
	rm -rf 'localCA.pem'
	rm -rf 'rootCA.pem'
	rm -rf 'coverage.out'

coverage:
	go test -covermode=count -coverprofile coverage.out ${PKG_LIST}

env:
	openssl genrsa -out localCA.key 4096
	openssl req -x509 -new -nodes -key localCA.key -sha256 -days 1024 -out localCA.crt
	openssl x509 -in localCA.crt -out localCA.pem -outform PEM

run:
	./bin/caos/${APP}

test:
	go test ./...

tidy:
	go mod tidy

vendor:
	go mod vendor
