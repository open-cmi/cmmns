ROOT := $(shell pwd)

TARGET=$(ROOT)/main/cmmns

.PHONY:build
build:
	cd internal/translation && gotext -srclang=en-US update -out=catalog-gen.go -lang=en-US,zh-CN github.com/open-cmi/cmmns/main
	cd main && go build -ldflags "-s -w" -o $(TARGET) main.go

BUILDDIR?=/usr/local
.PHONY:install
install:
	mkdir -p ${BUILDDIR}/bin
	cp -rfp ${TARGET} ${BUILDDIR}/bin/

.PHONY:clean
clean:
	rm -r build/*
