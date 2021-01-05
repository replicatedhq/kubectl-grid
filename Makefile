
export GO111MODULE=on

.PHONY: test
test:
	go test ./pkg/... ./cmd/... 

.PHONY: bin
bin: fmt vet
	go build -o bin/kubectl-grid github.com/replicatedhq/kubectl-grid/cmd/kubectl-grid

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet ./pkg/... ./cmd/...

