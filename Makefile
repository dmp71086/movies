# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN := $(CURDIR)/bin

# Добавляем bin в текущей директории в PATH при запуске protoc
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

# Путь до protobuf файлов
PROTO_PATH := $(CURDIR)/pkg/proto

# Путь до завендореных protobuf файлов
VENDOR_PROTO_PATH := $(CURDIR)/vendor.protobuf

# устанавливаем необходимые плагины
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	$(info Installing binary dependencies...)

	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/bufbuild/buf/cmd/buf@v1.32.2

# генерация .go файлов с помощью protoc
protoc-generate:
	$(PROTOC) -I $(CURDIR) --go_out=$(CURDIR) --go-grpc_out=$(CURDIR)  \
	$(PROTO_PATH)/movies.proto
	

# go mod edit -replace=google.golang.org/grpc=github.com/grpc/grpc-go@latest
# go mod tidy

# Генерация кода из protobuf
generate: .bin-deps .protoc-generate .tidy
	
# Объявляем, что текущие команды не являются файлами и
# интсрументируем Makefile не искать изменения в файловой системе
.PHONY: \
	.bin-deps \
	.protoc-generate \
	.tidy \
	generate \
	build

vendor:	.vendor-reset .vendor-google-protobuf .vendor-protovalidate 

.vendor-reset:
	rm -rf $(VENDOR_PROTO_PATH)
	mkdir -p $(VENDOR_PROTO_PATH)

# Устанавливаем proto описания google/protobuf
.vendor-google-protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf $(VENDOR_PROTO_PATH)/protobuf
	mv $(VENDOR_PROTO_PATH)/protobuf/src/google /usr/bin/include
	rm -rf $(VENDOR_PROTO_PATH)/protobuf

.vendor-protovalidate:
	git clone -b main --single-branch --depth=1 --filter=tree:0 \
		https://github.com/bufbuild/protovalidate $(VENDOR_PROTO_PATH)/protovalidate && \
	mv $(VENDOR_PROTO_PATH)/protovalidate/proto/protovalidate/buf /usr/bin/include
	rm -rf $(VENDOR_PROTO_PATH)/protovalidate