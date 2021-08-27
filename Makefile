APP = midiserver
VERSION = $(shell cat VERSION)

BIN_APP = bin/$(APP)
CMD_APP = cmd/$(APP)

DVCS_HOST = github.com
ORG = geomyidia
PROJ = $(APP)
FQ_PROJ = $(DVCS_HOST)/$(ORG)/$(PROJ)

GO_VERSION_STRING = $(shell go version)
GO_VERSION = $(strip $(subst go, , $(word 3, $(GO_VERSION_STRING))))
GO_ARCH = $(strip $(word 4, $(GO_VERSION_STRING)))
LD_VERSION = -X $(FQ_PROJ)/pkg/version.version=$(VERSION)
LD_BUILDDATE = -X $(FQ_PROJ)/pkg/version.buildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LD_GITCOMMIT = -X $(FQ_PROJ)/pkg/version.gitCommit=$(shell git rev-parse --short HEAD)
LD_GITBRANCH = -X $(FQ_PROJ)/pkg/version.gitBranch=$(shell git rev-parse --abbrev-ref HEAD)
LD_GITSUMMARY = -X $(FQ_PROJ)/pkg/version.gitSummary=$(shell git describe --tags --dirty --always)
LD_GO_VERSION = -X $(FQ_PROJ)/pkg/version.goVersion=$(GO_VERSION)
LD_GO_ARCH = -X $(FQ_PROJ)/pkg/version.goArch=$(GO_ARCH)

LDFLAGS = -w -s $(LD_VERSION) $(LD_BUILDDATE) $(LD_GITBRANCH) $(LD_GITSUMMARY) $(LD_GITCOMMIT) $(LD_GO_VERSION) $(LD_GO_ARCH)

MAINS = cmd/%/main.go
CMDS = $(wildcard cmd/*/main.go)
BINS = $(patsubst $(MAINS),bin/%,$(CMDS))

default: all

# For Github Action release builds
linux:
	@sudo apt-get install libasound2-dev libsndfile-dev
	$(MAKE) all

# For Github Action release builds
darwin: all

goversion:
	@echo $(GO_VERSION)

goarch:
	@echo $(GO_ARCH)

all: $(BINS)

bin/%: $(MAINS)
	@echo ">> Building $@ ..."
	@go build -race -ldflags "$(LDFLAGS)" -o ./$@ ./$<

clean:
	@echo ">> Removing $(BINS) ..."
	@rm -f $(BINS)

serve: all
	@echo ">> Serving from compiled binary ..."
	@$(BIN_APP) -loglevel debug -daemon

run:
	@echo ">> Running ..."
	@GO111MODULE=on go run ./$(CMD_APP)

help: all
	@echo ">> Getting binary version info ..."
	@$(BIN_APP) -loglevel error -h

version: all
	@echo ">> Getting binary version info ..."
	@$(BIN_APP) -loglevel error -version

test:
	@echo ">> Running unit tests ..."
	@export PATH=$$PATH:~/go/bin && richgo test -race -v ./... || echo "Uh-oh ... 🔥"

rebuild: clean all

.PHONY: test
