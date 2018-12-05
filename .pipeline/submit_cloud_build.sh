#!/bin/bash

GEN_COMMIT_SHA=$(openssl rand -base64 32 | tr -dc A-Za-z0-9 | head -c 12)

source SUBSTITUTIONS

gcloud builds submit . --config=.pipeline/cloudbuild-ci.yaml \
    --substitutions=COMMIT_SHA="${GEN_COMMIT_SHA}",_GCS_BUCKET="${GCS_BUCKET}",_ZONE="${ZONE}",_CLUSTER_NAME="${CLUSTER_NAME}",_ARTIFACTORY="${ARTIFACTORY}"
