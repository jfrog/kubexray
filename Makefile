
.PHONY: image
image:
	@echo "++ Building kubexray docker image..."
	rm -f cmd/kubexray/go.sum && docker build -t kubexray .

.PHONY: build
build:
	@echo "++ Building kubexray go binary..."
	mkdir -p bin
	cd cmd/kubexray && rm -f go.sum && go build -a --installsuffix cgo --ldflags="-s" -o ../../bin/kubexray

.PHONY: cloud
cloud:
	@echo "++ Submiting CI cloud build..."
	.scripts/submit_cloud_build.sh

.PHONY: creds
creds:
	@echo "++ Uploading creds files to GCS bucket..."
	.scripts/upload_creds.sh
