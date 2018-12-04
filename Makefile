
define gen_commit_sha
$(eval GEN_COMMIT_SHA := $(shell openssl rand -base64 32 | tr -dc A-Za-z0-9 | head -c 12))
endef

define xray_config
$(eval XRAY_CONFIG := $(shell cat charts/files/xray_config.yaml | base64))
endef

.PHONY: image
image:
	@echo "++ Building kubexray docker image"
	rm -f cmd/kubexray/go.sum && docker build -t kubexray .

.PHONY: build
build:
	mkdir -p bin
	cd cmd/kubexray && rm -f go.sum && go build -a --installsuffix cgo --ldflags="-s" -o ../../bin/kubexray

.PHONY: cloud
cloud:
	@$(call gen_commit_sha)
	gcloud builds submit . --config=.pipeline/cloudbuild.yaml --substitutions=COMMIT_SHA="${GEN_COMMIT_SHA}"

.PHONY: encrypt
encrypt:
	gcloud kms encrypt --key=kubexray-ci --keyring=kubexray-ci --location=global --ciphertext-file=artifactory.creds.enc --plaintext-file=artifactory.creds
	gcloud kms encrypt --key=kubexray-ci --keyring=kubexray-ci --location=global --ciphertext-file=charts/files/xray_config.yaml.enc --plaintext-file=charts/files/xray_config.yaml
