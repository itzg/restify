
GLIDE := ${GOPATH}/bin/glide
GORELEASER := ${GOPATH}/bin/goreleaser

default: build

${GLIDE}:
	curl https://glide.sh/get | sh

${GORELEASER}:
	go get github.com/goreleaser/goreleaser

vendor: ${GLIDE}
	${GLIDE} install

build: vendor
	go build

snapshot: vendor ${GORELEASER}
	${GORELEASER} --snapshot --skip-validate

release: vendor ${GORELEASER}
	${GORELEASER}

.PHONY: default build snapshot release