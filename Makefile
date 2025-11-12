SDK_ROOT ?= ~/git/aws-sdk-go-v2

.PHONY: set-local

set-local:
	cd ${SDK_ROOT} && ./local-mod-replace.sh -d $(shell pwd)

reset-local:
	git checkout -- go.mod
