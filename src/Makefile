# Variable for the environment mode 'production' or 'development'
ENV ?= production

GOBIN ?= $(shell go env GOPATH)/bin
NODEBIN ?= ./node_modules/.bin

TMP_DIR := ./tmp
BUILD_DIR := ./build
CMD_DIR := ./cmd
VIEW_DIR := ./view
JS_DIR := $(VIEW_DIR)/js
REPOSITORY_DIR := ./repository

GO_FILES := $(shell find . -path ./node_modules -prune -o -path $(VIEW_DIR) -prune -o -name '*.go' -print)
TEMPL_FILES := $(shell find $(VIEW_DIR) -name '*.templ')
TEMPL_GO_FILES = $(TEMPL_FILES:.templ=_templ.go)
JS_FILES := $(shell find $(JS_DIR) -name '*.js')
SQL_FILES := $(wildcard ./db/migrations/*.sql) $(wildcard ./db/query/*.sql)

MAINJS_OUT := ./assets/static/js/main.js
TAILWINDCSS_OUT := ./assets/static/css/tailwind.css

TAILWINDCSS_LOG := tailwind.log
SQLC_SENTINEL := $(REPOSITORY_DIR)/.sqlc_generated

$(TMP_DIR)/%: $(BUILD_DIR)/% $(MAINJS_OUT) $(TAILWINDCSS_OUT)
	touch "$@"

# The final build step
ifeq ($(ENV),development)
$(BUILD_DIR)/%: $(GO_FILES) $(TEMPL_GO_FILES) $(SQLC_SENTINEL)
	go build -tags=dev -o "$@" "$(CMD_DIR)/$*"
else
$(BUILD_DIR)/%: $(GO_FILES) $(TEMPL_GO_FILES) $(SQLC_SENTINEL) $(MAINJS_OUT) $(TAILWINDCSS_OUT)
	go build -o "$@" "$(CMD_DIR)/$*"
endif

# Build step for templ source
%_templ.go: %.templ
	"$(GOBIN)/templ" generate -f "$<"

# Build step for sqlc source
$(SQLC_SENTINEL): $(SQL_FILES) sqlc.yml
	sqlc generate
	touch "$@"

# Build step for main.js
$(MAINJS_OUT): $(JS_FILES)
ifeq ($(ENV),development)
	"$(NODEBIN)/esbuild" "$(JS_DIR)/main.js" --outfile="$@" --bundle
else
	"$(NODEBIN)/esbuild" "$(JS_DIR)/main.js" --outfile="$@" --bundle --minify
endif

# Build step for tailwind.css
$(TAILWINDCSS_OUT): $(TEMPL_FILES) tailwind.config.cjs tailwind.css
ifeq ($(ENV),development)
	./scripts/tailwind.sh "$(TMP_DIR)" "$(TAILWINDCSS_LOG)"
else
	"$(NODEBIN)/tailwindcss" build -i tailwind.css -o "$@" --minify
endif
	touch "$@"

# Live reload
.PHONY: live
live:
ifeq ($(ENV),development)
	mkdir -p "$(TMP_DIR)"
	"$(NODEBIN)/tailwindcss" build -i tailwind.css -o "$(TAILWINDCSS_OUT)" --watch=always &> "$(TMP_DIR)/$(TAILWINDCSS_LOG)" &
	node ./scripts/live-reload.cjs &
	"$(GOBIN)/air"
else
	@echo "Must run with ENV=development"
endif

.PHONY: clean
clean:
	@if [ -z "$(VIEW_DIR)" ] || [ -z "$(BUILD_DIR)" ] || [ -z "$(REPOSITORY_DIR)" ]; then \
		echo "Error: VIEW_DIR, BUILD_DIR or REPOSITORY_DIR is not set."; \
		exit 1; \
	fi
	find "$(VIEW_DIR)" -type f -name '*_templ.go' -exec rm -f {} +
	rm -f "$(REPOSITORY_DIR)"/*.go "$(SQLC_SENTINEL)"
	rm -f "$(TAILWINDCSS_OUT)" "$(MAINJS_OUT)"
	rm -rf "$(BUILD_DIR)"
	rm -rf "$(TMP_DIR)"

.SECONDARY:
