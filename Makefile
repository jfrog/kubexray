.PHONY: build-local
build-local:
	mkdir -p bin
	cd cmd/kube-xray && rm -f go.sum && go build -a --installsuffix cgo --ldflags="-s" -o ../../bin/kube-xray

.PHONY: build
build:
	@echo "++ Building kube-xray linux amd64"
	mkdir -p bin/
	cd cmd/kube-xray && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o ../../bin/kube-xray

.PHONY: image
image:
	@echo "++ Building kube-xray docker image"
	rm -f cmd/kube-xray/go.sum && docker build -t kube-xray .
