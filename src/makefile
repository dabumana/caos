APP=caos
VERSION=v.0.0.0

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

run:
	./bin/${APP}-${VERSION}

test:
	go test -cover ./...

update:
	go mod tidy

env:
	openssl genrsa -out localCA.key 4096
	openssl req -x509 -new -nodes -key localCA.key -sha256 -days 1024 -out localCA.crt
	openssl x509 -in localCA.crt -out localCA.pem -outform PEM
