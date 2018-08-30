NOVENDOR = $(shell go list ./... | grep -v vendor | grep -v node_modules)
NOVENDOR_LINTER = $(shell go list ./... | grep -v vendor | grep -v ptypes | grep -v node_modules)

all: build

fix:
	go fix $(NOVENDOR)
.PHONY: fix

vet:
	go vet $(NOVENDOR)
.PHONY: vet

lint:
	printf "%s\n" "$(NOVENDOR)" | xargs -I {} sh -c 'golint -set_exit_status {}'
.PHONY: lint

test:
	go test -v -cover $(NOVENDOR)
.PHONY: test

metalinter:
	gometalinter --config .gometalinter.json $(NOVENDOR_LINTER)
.PHONY: metalinter

clean:
	rm -rf ./bin
.PHONY: clean

build: clean fix vet lint test
	mkdir bin
	GOOS=linux GOARCH=386 go build -v -o ./bin/dapi main.go
	GOOS=windows GOARCH=amd64 go build -v -o ./bin/dapi.exe main.go
.PHONY: build

docker: build
	docker build -f Dockerfile -t gcr.io/rf-services/dapi:latest .
	docker push gcr.io/rf-services/dapi:latest
	docker system prune --volumes -f
.PHONY: docker

deploy: docker
	scp update.sh dapi:dapi/update.sh
	scp docker-compose.yml dapi:dapi/docker-compose.yml
.PHONY: deploy
