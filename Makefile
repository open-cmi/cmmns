ROOT := $(shell pwd)

TARGET=$(ROOT)/main/cmmns

.PHONY:build
build:
	cd main && go build -ldflags "-s -w" -o $(TARGET) main.go

BUILDDIR?=/usr/local
.PHONY:install
install:
	mkdir -p ${BUILDDIR}/bin
	cp -rfp ${TARGET} ${BUILDDIR}/bin/

.PHONY:clean
clean:
	rm -r build/*
