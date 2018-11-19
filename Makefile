DIRS  = $(shell cd services && ls -d */ | grep -v "_output")
ODIR 	:= bin/

export VAR_SERVICES ?= $(DIRS:/=)

compile: 
	@$(foreach svc, $(VAR_SERVICES), \
		mkdir -p services/$(svc)/bin && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o services/$(svc)/bin/$(svc) services/$(svc)/main.go;)