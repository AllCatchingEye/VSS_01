.PHONY: regenerate-docker
 
CWD=$(shell pwd |  sed 's/ /\\ /g')
 
regenerate-docker:
	docker run --rm -it \
		-v $(CWD):/messages -w /messages \
		gitlab.lrz.de:5005/hm/dockerhub/go-protoc:latest \
		protoc -I=. --go_out=. --go_opt=paths=source_relative --proto_path=. messages.proto
