
define xray_config
$(eval XRAY_CONFIG := $(shell cat charts/files/xray_config.yaml | base64))
endef

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
	.pipeline/submit_cloud_build.sh

.PHONY: encrypt
encrypt:
	gcloud kms encrypt --key=kubexray-ci --keyring=kubexray-ci --location=global --ciphertext-file=charts/files/artifactory.creds.enc --plaintext-file=charts/files/artifactory.creds
	gcloud kms encrypt --key=kubexray-ci --keyring=kubexray-ci --location=global --ciphertext-file=charts/files/xray_config.yaml.enc --plaintext-file=charts/files/xray_config.yaml
	gcloud kms encrypt --key=kubexray-ci --keyring=kubexray-ci --location=global --ciphertext-file=charts/files/bintray.creds.enc --plaintext-file=charts/files/bintray.creds
