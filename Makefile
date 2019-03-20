
.PHONY: image
image:
	@echo "++ Building kubexray docker image..."
	docker build -t kubexray .

.PHONY: build
build: export GOARCH=amd64
build: export CGO_ENABLED=0
build: export GO111MODULE=on
build: export GOPROXY=https://gocenter.io
build:
	@echo "++ Building kubexray go binary..."
	mkdir -p bin
	cd cmd/kubexray && go build -a --installsuffix cgo --ldflags="-s" -o ../../bin/kubexray

.PHONY: cloud
cloud:
	@echo "++ Submiting CI cloud build..."
	.scripts/submit_cloud_build.sh

.PHONY: creds
creds:
	@echo "++ Uploading creds files to GCS bucket..."
	.scripts/upload_creds.sh
