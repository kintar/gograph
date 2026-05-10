BINARY    = gograph
BUILD_DIR = bin
CMD       = ./cmd/gograph
INSTALL   = /usr/local/bin

.PHONY: build test run-build clean bump-patch bump-minor bump-major install release

# Read the current version from .bumpversion.cfg after bumping.
# grep picks the 'current_version = x.y.z' line; awk extracts the last field.
build: bump-patch
	$(eval VERSION := $(shell grep '^current_version' .bumpversion.cfg | awk '{print $$3}'))
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY) $(CMD)
	@echo "Built $(BUILD_DIR)/$(BINARY) v$(VERSION)"

release:
	$(eval VERSION := $(shell grep '^current_version' .bumpversion.cfg | awk '{print $$3}'))
	git tag -a v$(VERSION) -m "Release v$(VERSION)" || true
	git push origin master --tags

# install only copies whatever is already in bin/ — no implicit build.
install:
	@test -f $(BUILD_DIR)/$(BINARY) || (echo "Run 'make build' first — $(BUILD_DIR)/$(BINARY) not found." && exit 1)
	sudo rm -f $(INSTALL)/$(BINARY)
	sudo cp $(BUILD_DIR)/$(BINARY) $(INSTALL)/
	@echo "Installed $(BINARY) to $(INSTALL)/"

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

test-fuzz:
	@echo "Running FuzzConstructors for 5s..."
	go test -fuzz=FuzzConstructors -fuzztime=5s ./internal/search
	@echo "Running FuzzSchema for 5s..."
	go test -fuzz=FuzzSchema -fuzztime=5s ./internal/search

run-build:
	go run $(CMD) build .

clean:
	rm -rf $(BUILD_DIR)

bump-patch:
	bump2version --no-commit patch --allow-dirty

bump-minor:
	bump2version --no-commit minor --allow-dirty

bump-major:
	bump2version --no-commit major --allow-dirty
