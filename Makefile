
define gen_commit_sha
$(eval GEN_COMMIT_SHA := $(shell openssl rand -base64 32 | tr -dc A-Za-z0-9 | head -c 12))
endef

define xray_config
$(eval XRAY_CONFIG := $(shell cat charts/files/xray_config.yaml | base64))
endef

.PHONY: build
build:
	@echo "++ Building kube-xray"
	CGO_ENABLED=0 GOOS=linux cd cmd/kube-xray && go build -a -tags netgo -ldflags "$(LDFLAGS) -w -s" -o kube-xray .

.PHONY: cloud
cloud:
	@$(call gen_commit_sha)
	gcloud builds submit . --config=.pipeline/cloudbuild.yaml --substitutions=COMMIT_SHA="${GEN_COMMIT_SHA}"

.PHONY: docker
docker:
	docker build -t kube-xray .

.PHONY: go
go:
	mkdir -p bin
	cd cmd/kube-xray && go build -a --installsuffix cgo --ldflags="-s" -o ../../bin/kube-xray

.PHONY: run
run:
	cd bin && ./kube-xray

.PHONY: dry-run
dry-run:
	@$(call xray_config)
	helm tiller run -- helm upgrade --install kube-xray --namespace kube-xray charts/kube-xray/ --set xrayConfig="${XRAY_CONFIG}" --dry-run --debug
