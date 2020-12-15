export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

LDFLAGS := "-s -w"

.PHONY: all
all: clean fmt build


.PHONY: build
build: 
	@echo "build binary file"; \
 	go generate; \
	GOOS=linux GOARCH=arm64 go build -trimpath -ldflags $(LDFLAGS) -o build/OpenWrtDingTalkSender_linux_arm64 .; \
	upx -9 build/OpenWrtDingTalkSender_linux_arm64; \
	packr2 clean

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: clean
clean: 
	@echo "clean"; \
	rm -rf ./build; \
	packr2 clean

