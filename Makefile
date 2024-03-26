.PHONY: openapi openapi-codegen-install

check-env:
ifndef GOPATH
	$(error GOPATH is undefined)
endif

openapi-codegen-install: check-env
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

openapi: check-env openapi-codegen-install
	${GOPATH}/bin/oapi-codegen  -old-config-style -package="api" pkg/api/openapi.yaml > pkg/api/openapi.go
