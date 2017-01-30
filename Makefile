.PHONY: all clean get-deps build lint-go lint-jss lint fmt docker install uninstall

UNAME_S := $(shell uname -s)
PKG_MGR := apt-get

all: build

clean:
	rm -rf ./public/assets/node_modules/ ./public/assets/build/ ./dist/

get-deps:
	if [ -z "$(which curl)" ]; then \
		${PKG_MGR} install -y curl; \
	fi
	if [ -z "$(which npm)" ]; then \
		curl -sL https://deb.nodesource.com/setup_6.x | bash - && \
		${PKG_MGR} install -y nodejs; \
	fi
	if [ -z "$(which gulp)" ]; then \
		npm install -g gulp; \
	fi
	if [ -z "$(which go-bindata)" ]; then \
		go get -u github.com/jteeuwen/go-bindata/...; \
	fi
	if [ -z "$(which glide)" ]; then \
		curl https://glide.sh/get | sh; \
	fi

build: get-deps
	glide install
	cd public/assets/ && \
	npm install && \
	gulp
	go-bindata -pkg public -o public/public.go public/assets/build/...
	mkdir -p ./dist/
	go build -o ./dist/maestro-www main.go

lint-go:
	go lint ./...

lint-js:
	cd public/assets/ && \
	gulp lint

lint: lint-go lint-js

fmt:
	go fmt ./...

docker:
	docker build -t maestro-www-build .
	mkdir -p ./dist/
	docker run --rm -it -v `pwd`/dist/:/opt/maestro-www/dist/ maestro-www-build
	docker build -t maestro-www -f Dockerfile_scratch .

install:
	mkdir -p /opt/bin/
	cp ./dist/maestro-www /opt/bin/

uninstall:
	rm /opt/bin/maestro-www

