VERSION := 0.1.0
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date +"%s")

default: build

test:
	go test -v ./...

build:
	go build -ldflags "\
		 -X main.version=$(VERSION) \
		 -X main.commit=$(COMMIT) \
		 -X main.buildDate=$(BUILD_DATE)"\
		 -o bundles/gohit

container-build: bundles
	docker build -t gohit-build .

binary: container-build
	docker run --rm -v "$(CURDIR)/bundles/container:/go/src/github.com/fabiofalci/gohit/bundles" gohit-build make build

clean:
	rm -rf bundles/

bundles:
	mkdir -p bundles/container

#
# Only works on osx as we can generate both osx and linux binaries in one go.
#
release: clean binary build
	mkdir -p bundles/release/{linux,osx}
	cp bundles/gohit bundles/release/osx/
	cp bundles/container/gohit bundles/release/linux/
	zip -j bundles/release/osx-x86_64-gohit-$(VERSION).zip bundles/release/osx/gohit
	zip -j bundles/release/linux-x86_64-gohit-$(VERSION).zip bundles/release/linux/gohit

