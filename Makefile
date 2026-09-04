# Steampipe Plugin Makefile

PLUGIN_NAME  := tencentcloud
BUILD_DIR    := build
BINARY       := $(BUILD_DIR)/$(PLUGIN_NAME).plugin
LOCAL_DIR    := $(HOME)/.steampipe/plugins/local/$(PLUGIN_NAME)
LDFLAGS      := -s -w

.PHONY: build install clean test test-query test-sts fmt vet

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags '$(LDFLAGS)' -o $(BINARY) .

install: build
	@mkdir -p $(LOCAL_DIR)
	cp $(BINARY) $(LOCAL_DIR)/
	@if [ "$$(uname)" = "Darwin" ]; then \
		xattr -cr $(LOCAL_DIR)/$(PLUGIN_NAME).plugin 2>/dev/null; \
		codesign --force --sign - $(LOCAL_DIR)/$(PLUGIN_NAME).plugin 2>/dev/null \
			&& echo "Signed plugin (adhoc) for macOS"; \
	fi
	@echo "Plugin installed to $(LOCAL_DIR)/"
	@echo "Restart Steampipe service: steampipe service restart"

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./... -v

test-query:
	go test -tags=querytest ./tencentcloud-test -v

test-sts:
	go test -tags=integration -run TestSTSEndToEnd -v -timeout 120s ./tencentcloud/utils/

fmt:
	go fmt ./...

vet:
	go vet ./...
