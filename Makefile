
# Image URL to use all building/pushing image targets
COMPONENT        ?= armada-operator
VERSION          ?= 0.0.1
DHUBREPO         ?= kubekit99/${COMPONENT}-dev
DOCKER_NAMESPACE ?= kubekit99
IMG              ?= ${DHUBREPO}:v${VERSION}

all: vet

setup:
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	# dep ensure

clean:
	rm -fr vendor
	rm -fr cover.out
	rm -fr build/_output
	rm -fr config/crds

# Run tests
unittest: setup fmt vet
	go test ./armada/... -coverprofile cover.out

# Run go fmt against code
fmt: setup
	go fmt ./armada/... 

# Run go vet against code
vet: setup
	go vet ./armada/...
