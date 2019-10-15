NAME=jxcore
BINDIR=build/bin
PACKAGE_DIR=build/package-$(arch)
EXTRACT_DIR=edge/jxcore

version=$(shell git tag -l "v*" --points-at HEAD | tail -n 1 | tail -c +2 )
commit=$(shell git rev-parse --short HEAD)
builddate=$(shell date "+%m/%d/%Y %R %Z")
arch?=arm64

REPO=jxcore
GO=CGO_ENABLED=0 GO111MODULE=off go
GOFLAGS=-v -ldflags '-X "$(REPO)/version.Version=$(version)" -X "$(REPO)/version.GitCommit=$(commit)" -X "$(REPO)/version.BuildDate=$(builddate)"'

.PHONY: build debian test clean check_version push_to_source

check_version:
ifeq ($(version),)
	$(error No version tag found)
endif
	
build: check_version
	$(GO) get
	GOARCH=$(arch) $(GO) build $(GOFLAGS) -o $(BINDIR)/$(NAME)-$(arch)

debian: check_version build
	mkdir -p $(PACKAGE_DIR)/$(EXTRACT_DIR)/bin
	cp -r $(BINDIR)/$(NAME)-$(arch) $(PACKAGE_DIR)/$(EXTRACT_DIR)/bin/$(NAME)
	cp -r settings.yaml $(PACKAGE_DIR)/$(EXTRACT_DIR)/bin
	cp -r scripts/jxcore_service.sh $(PACKAGE_DIR)/$(EXTRACT_DIR)/bin
	cp -r doc $(PACKAGE_DIR)/$(EXTRACT_DIR)/template
	echo $(version) > $(PACKAGE_DIR)/$(EXTRACT_DIR)/VERSION
	
	mkdir -p $(PACKAGE_DIR)/DEBIAN/
	sed -e "s/REPLACE_VERSION/$(version)/g" -e "s/REPLACE_ARCH/$(arch)/" DEBIAN/control > $(PACKAGE_DIR)/DEBIAN/control
	dpkg -b $(PACKAGE_DIR) build/$(NAME)_$(version)_$(arch).deb

	rm -rf $(PACKAGE_DIR)

push_to_source:
	bash scripts/upload.sh

test:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

clean:
	rm -rf build/