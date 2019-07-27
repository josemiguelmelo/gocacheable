go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=on go mod vendor

test:
	@echo "  >  Testing gocacheable..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=on go test -v -race -coverprofile=gocacheable.coverprofile ./...
