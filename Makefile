APP = midiserver
VERSION = $(shell cat VERSION)

BIN_APP = bin/$(APP)
CMD_APP = cmd/$(APP)

DVCS_HOST = github.com
ORG = geomyidia
PROJ = erl-midi-server
FQ_PROJ = $(DVCS_HOST)/$(ORG)/$(PROJ)

GO_VERSION_STRING = $(shell go version)
GO_VERSION = $(strip $(subst go, , $(word 3, $(GO_VERSION_STRING))))
GO_ARCH = $(strip $(word 4, $(GO_VERSION_STRING)))
LD_VERSION = -X $(FQ_PROJ)/internal/app.version=$(VERSION)
LD_BUILDDATE = -X $(FQ_PROJ)/internal/app.buildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LD_GITCOMMIT = -X $(FQ_PROJ)/internal/app.gitCommit=$(shell git rev-parse --short HEAD)
LD_GITBRANCH = -X $(FQ_PROJ)/internal/app.gitBranch=$(shell git rev-parse --abbrev-ref HEAD)
LD_GITSUMMARY = -X $(FQ_PROJ)/internal/app.gitSummary=$(shell git describe --tags --dirty --always)
LD_GO_VERSION = -X $(FQ_PROJ)/internal/app.goVersion=$(GO_VERSION)
LD_GO_ARCH = -X $(FQ_PROJ)/internal/app.goArch=$(GO_ARCH)

LDFLAGS = -w -s $(LD_VERSION) $(LD_BUILDDATE) $(LD_GITBRANCH) $(LD_GITSUMMARY) $(LD_GITCOMMIT) $(LD_GO_VERSION) $(LD_GO_ARCH)

default: all

all: build

build: $(BIN_APP)

bin:
	@mkdir ./bin

$(BIN_APP): bin
	@GO111MODULE=on go build -ldflags "$(LDFLAGS)" -o $(BIN_APP) ./$(CMD_APP)

start: build
	@$(BIN_APP)

run:
	@GO111MODULE=on go run ./$(CMD_APP)

clean:
	@rm -f $(BIN_APP)
