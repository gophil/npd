.PHONY: all clean

BUILDFILE = ./example/sample_wait_args.go
GOOS = darwin
GOARCH= amd64

all: format test build

test:
	go test -v . 

format:
	gofmt -w .

build:
	mkdir -p builds
	# 设置交叉编译参数:
	# GOOS为目标编译系统, mac os则为 "darwin", window系列则为 "windows"
	# 生成二进制执行文件 akbs , 如在windows下则为 akbs.exe
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -o builds/npd $(BUILDFILE)

clean:
	go clean -i