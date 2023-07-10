APP=caos
VERSION=v.0.0.0
PKG_LIST=$(go list ./... | grep -v /vendor/ | tr '\n' ' ')

build:
	go build -o ./bin/${APP}-${VERSION}

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

run:
	./bin/${APP}-${VERSION}

test:
	go test ./...

coverage:
	go test -covermode=count -coverprofile coverage.out ${PKG_LIST}

env:
	openssl genrsa -out localCA.key 4096
	openssl req -x509 -new -nodes -key localCA.key -sha256 -days 1024 -out localCA.crt
	openssl x509 -in localCA.crt -out localCA.pem -outform PEM
