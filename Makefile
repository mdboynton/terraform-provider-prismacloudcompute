HOSTNAME = registry.terraform.io
NAMESPACE = PaloAltoNetworks
NAME = prismacloudcompute
BINARY = terraform-provider-${NAME}

VERSION ?= 0.0.0-dev
OS_ARCH ?= darwin_amd64

# For Apple silicon machines:
#OS_ARCH ?= darwin_arm64 

default: install

format:
	gofmt -l -w .

build:
	go build -o ${BINARY}
	#go build -gcflags="all=-N -l" -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

clean:
	rm -rf ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}

acctest: build
	go test -v ./internal/acceptance/ -count=1
