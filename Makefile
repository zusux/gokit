MODULE := $(shell go list -m)

ZRPC_SRC := third_party/zusux/zrpc/zrpc.proto
ZRPC_DST_DIR := gserver/zrpc
ZRPC_DST := $(ZRPC_DST_DIR)/zrpc.proto
OPENAPI_OUT := gserver/openapi

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/zusux/gokit/gen/protoc-gen-tag@latest


.PHONY: patch-zrpc
patch-zrpc:
	@echo "==> Copying zrpc.proto..."
	@mkdir -p $(ZRPC_DST_DIR)
	@cp $(ZRPC_SRC) $(ZRPC_DST)
	@echo "==> Rewriting go_package in zrpc.proto..."
	@rm -f $(ZRPC_DST).bak

.PHONY: gen-zrpc
gen-zrpc: patch-zrpc
	@echo "==> Generating zrpc.pb.go..."
	@protoc \
		--proto_path=. \
        --proto_path=./third_party \
        --openapi_out=fq_schema_naming=true,default_response=false:. \
		--go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		--include_imports \
		$(ZRPC_DST)
	@rm -f $(ZRPC_DST)

.PHONY: gen-desc
gen-desc:
	@echo "==> Generating desc/all.desc..."
	@protoc \
		--proto_path=. \
		--proto_path=./third_party \
		--descriptor_set_out=./desc/all.desc \
		--include_imports \
		$(shell find api -name "*.proto")

.PHONY: gen-api
gen-api:
	@echo "==> Generating other api proto files..."
	@find gserver/invoke/api -name '*.proto' ! -name 'zrpc.proto' -print0 | \
	xargs -0 -I{} protoc \
		--proto_path=. \
		--proto_path=./third_party \
		--openapi_out=fq_schema_naming=true,default_response=false:./gserver/invoke/api \
		--go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		--go-http_out=paths=source_relative:. \
		--tag_out=. \
        --include_imports \
		"{}"


.PHONY: gen-router
gen-router:
	@echo '==> Copying router.proto from third_party...'
	mkdir -p api/invoke
	cp third_party/zusux/invoke/router.proto api/invoke/router.proto
	@echo '==> Patching go_package in router.proto...'
	#sed -i.bak -E 's|option go_package = ".*";|option go_package = "$(MODULE)/api/invoke;invoke";|' api/invoke/router.proto
	rm -f api/invoke/router.proto.bak
	@echo '==> Generating invoke pb files...'
	protoc \
		--proto_path=. \
		--proto_path=./third_party \
		--go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		--tag_out=. \
		api/invoke/router.proto

.PHONY: api
api: gen-zrpc gen-api

