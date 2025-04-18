.PHONY: all clean

all: deps certs build

deps:
	go mod tidy


clean:
	rm -rf assets/certs
	rm -f bin/server

certs:
	mkdir -p assets/certs tmp
	openssl genpkey -algorithm RSA -out assets/certs/ca.key
	openssl req -x509 -new -key assets/certs/ca.key -out assets/certs/ca.crt -subj "/CN=MyCA"
	openssl genpkey -algorithm RSA -out assets/certs/server.key
	openssl req -new -key assets/certs/server.key -out assets/certs/server.csr -subj "/CN=localhost"
	echo "subjectAltName = DNS:localhost,IP:127.0.0.1" > tmp/extfile.cnf
	openssl x509 -req -in assets/certs/server.csr -CA assets/certs/ca.crt -CAkey assets/certs/ca.key -out assets/certs/server.crt -extfile tmp/extfile.cnf
	rm -Rf tmp/
	openssl genpkey -algorithm RSA -out assets/certs/client.key
	openssl req -new -key assets/certs/client.key -out assets/certs/client.csr -subj "/CN=client"
	openssl x509 -req -in assets/certs/client.csr -CA assets/certs/ca.crt -CAkey assets/certs/ca.key -out assets/certs/client.crt

build:
	go build -o bin/server cmd/server/main.go
