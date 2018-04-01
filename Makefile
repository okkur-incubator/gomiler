SHELL := /bin/bash

GO=$(shell which go)
BUILD=$(GO) build
CLEAN=$(GO) clean
TEST=$(GO) test
GET=$(GO) get
    
all: test build
build:
	$(BUILD) main.go

test:
	$(TEST) 

clean:
	$(CLEAN)

get:
	$(GET) github.com/peterhellberg/link
	$(GET) gopkg.in/jarcoal/httpmock.v1	

