HOSTNAME = registry.terraform.io
NAMESPACE = PaloAltoNetworks
NAME = prismacloudcompute
BINARY = terraform-provider-${NAME}

VERSION ?= 1.0.0-alpha
OS_ARCH ?= darwin_amd64

# For Apple silicon:
#OS_ARCH ?= darwin_arm64 

default: install

format:
	gofmt -l -w .

build:
	go build -gcflags="all=-N -l" -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

clean:
	rm -rf ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}
