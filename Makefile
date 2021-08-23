ROOT := $(shell pwd)

ifdef BUILD_DIR
TARGET=$(BUILD_DIR)/bin/
else
TARGET=$(ROOT)/
endif

.PHONY:build
build:
	cd cmd && go build -ldflags "-s -w" -o $(TARGET)/cmmns main.go

clean:
	rm -r build/*
