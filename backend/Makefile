.PHONY : build prod run fresh test clean dist deps

BIN := newsletter.bin

HASH := $(shell git rev-parse --short HEAD)
COMMIT_DATE := $(shell git show -s --format=%ci ${HASH})
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
VERSION := ${HASH} (${COMMIT_DATE})

STATIC := ../frontend/static:/static

deps:
	go get -u github.com/knadh/stuffbin/...
	cd ../frontend && yarn install

build:
	go build -o ${BIN} -ldflags="-X 'main.version=${VERSION}' -X 'main.date=${BUILD_DATE}'"

build-frontend:
	cd ../frontend && yarn build

run:
	./${BIN}

fresh: clean build run

test:
	go test

# dist builds the backend, frontend, and uses stuffbin to
# embed all frontend assets into the binary.
dist: build build-frontend
	stuffbin -a stuff -in ${BIN} -out ${BIN} ${STATIC}


clean:
	go clean
	rm -f ${BIN}