NAME=jxcore
PACKAGE_NAME=jx-jxcore
BINDIR=build/bin
PACKAGE_DIR=build/package-$(arch)
EXTRACT_DIR=edge/jxcore
CHANGELOG=build/CHANGELOG.md
RC_TAG=$(shell git tag -l "rc-*" --points-at HEAD)

version=$(shell git tag -l "v*" --points-at HEAD | tail -n 1 | tail -c +2 )
commit=$(shell git rev-parse --short HEAD)
builddate=$(shell date "+%m/%d/%Y %R %Z")
arch?=arm64

REPO=jxcore
GO=CGO_ENABLED=0 GO111MODULE=on go
GOFLAGS=-v -ldflags '-X "$(REPO)/version.Version=$(version)" -X "$(REPO)/version.GitCommit=$(commit)" -X "$(REPO)/version.BuildDate=$(builddate)"'

.PHONY: build debian test clean check_version push_to_source check_rc

check_version:
ifeq ($(version),)
	$(error No version tag found)
endif
ifeq ($(shell git cat-file -t v$(version)),commit)
	$(error Changelog should be in tag message)
endif

check_rc:
ifeq ($(RC_TAG),)
	$(error No rc-* tag found)
endif

build:
	GOARCH=$(arch) $(GO) build $(GOFLAGS) -o $(BINDIR)/$(NAME)-$(arch)

changelog:
	mkdir -p build
	echo "# Changelog\n" > $(CHANGELOG)
	TZ=UTC-8 git tag -l 'v*' --sort=-creatordate --format='## %(refname:short) - %(creatordate:iso-local)%0a### Author: %(taggername) %(taggeremail)%0a%(contents)%0a' >> $(CHANGELOG)

upload_changelog: changelog
	curl -s --fail -F "changelog=@$(CHANGELOG)" "http://packages.debian.jiangxingai.com:8000/api/v1/packages/$(PACKAGE_NAME)/changelog"

debian_base:
	mkdir -p $(PACKAGE_DIR)/$(EXTRACT_DIR)/bin
	cp -r $(BINDIR)/$(NAME)-$(arch) $(PACKAGE_DIR)/$(EXTRACT_DIR)/bin/$(NAME)
	chmod +x $(PACKAGE_DIR)/$(EXTRACT_DIR)/bin/$(NAME)
	cp -r settings.yaml $(PACKAGE_DIR)/$(EXTRACT_DIR)/bin
	cp -r gateway.cfg $(PACKAGE_DIR)/$(EXTRACT_DIR)/bin
	mkdir -p $(PACKAGE_DIR)/etc/systemd/system/
	cp -r scripts/jxcore.service $(PACKAGE_DIR)/etc/systemd/system/
	cp -r template $(PACKAGE_DIR)/$(EXTRACT_DIR)/template
	cp -r DEBIAN $(PACKAGE_DIR)/DEBIAN/

debian: check_version build debian_base changelog
	cp  $(CHANGELOG) $(PACKAGE_DIR)/$(EXTRACT_DIR)/CHANGELOG.md
	echo $(version) > $(PACKAGE_DIR)/$(EXTRACT_DIR)/VERSION
	sed -e "s/REPLACE_VERSION/$(version)/g" \
		-e "s/REPLACE_ARCH/$(arch)/" \
		-e "s/REPLACE_PACKAGE_NAME/$(PACKAGE_NAME)/g" \
		DEBIAN/control > $(PACKAGE_DIR)/DEBIAN/control
	dpkg -b $(PACKAGE_DIR) build/$(NAME)_$(version)_$(arch).deb
	rm -rf $(PACKAGE_DIR)

debian_rc: check_rc build debian_base
	sed -e "s/REPLACE_VERSION/0.1.0-$(RC_TAG)/g" \
		-e "s/REPLACE_ARCH/$(arch)/g" \
		-e "s/REPLACE_PACKAGE_NAME/$(PACKAGE_NAME)-dev/g" \
		DEBIAN/control > $(PACKAGE_DIR)/DEBIAN/control
	dpkg -b $(PACKAGE_DIR) build/$(NAME)-dev_0.1.0-$(RC_TAG)_$(arch).deb
	rm -rf $(PACKAGE_DIR)

push_to_source:
	bash scripts/upload.sh

deploy_prerelease:
	bash scripts/deploy.sh $(PACKAGE_NAME) $(version) prerelease

deploy_rc:
	bash scripts/deploy.sh $(PACKAGE_NAME)-dev 0.1.0-$(RC_TAG) review
	printf "\033[0;31mTo download package $(PACKAGE_NAME), run in development environment: apt install $(PACKAGE_NAME)-dev=0.1.0-$(RC_TAG)\n"

test:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

clean:
	rm -rf build/
