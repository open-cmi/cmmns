ROOT := $(shell pwd)

TARGET=$(ROOT)/main/cmmns

.PHONY:build
build:
	cd main && go build -ldflags "-s -w" -o $(TARGET) main.go

PREFIX?=/usr/local
.PHONY:install
install:
	mkdir -p ${PREFIX}/bin
	cp -rfp ${TARGET} ${PREFIX}/bin/

.PHONY:clean
clean:
	rm -r build/*
