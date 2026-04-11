PLUGIN_ID = com.thecodingflow.hamlibplugin
BINARY_NAME = hamlibplugin
VERSION_NUMBER ?= $(shell git describe --tags 2>/dev/null || echo "v0.0.0-dev" | sed -E 's#v##')

BUILD_DIR = build
PLUGIN_DIR = ${BUILD_DIR}/${PLUGIN_ID}.sdPlugin
DEST_DIR ?= ~/.config/opendeck/plugins
ARCHIVE_DIR:=$(shell pwd)

version_number:
	@echo ${VERSION_NUMBER}

.PHONY: clean
clean:
	rm -rf ${BUILD_DIR}
	mkdir -p ${PLUGIN_DIR}

.PHONY: prepare
prepare: clean
	sed "s/--VERSION--/${VERSION_NUMBER}/g" manifest.json > ${PLUGIN_DIR}/manifest.json
	cp -r icons ${PLUGIN_DIR}/
	cp -r pi ${PLUGIN_DIR}/

.PHONY: build
build:
	echo "$(VERSION_NUMBER)"
	GOARCH=amd64 GOOS=linux go build -v -ldflags "-X main.version=${VERSION_NUMBER}" -o "${PLUGIN_DIR}/${BINARY_NAME}-x86_64-unknown-linux-gnu" cmd/main.go

.PHONY: test
test:
	go test -v -timeout=30s ./...

.PHONY: install
install: uninstall
	cp -r ${PLUGIN_DIR} ${DEST_DIR}/

.PHONY: uninstall
uninstall:
	rm -rf ${DEST_DIR}/${PLUGIN_ID}.sdPlugin

.PHONY: archive
archive: prepare build
	(cd ${BUILD_DIR} && zip -r ${ARCHIVE_DIR}/${PLUGIN_ID}.zip ${PLUGIN_ID}.sdPlugin/)
