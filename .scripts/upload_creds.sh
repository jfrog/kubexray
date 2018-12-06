#!/bin/bash

source SUBSTITUTIONS

# files list
FILES_LIST=( artifactory.creds bintray.creds slack.creds xray_config.yaml )

cd credentials
mkdir -p backups
DATE_STAMP=$(date +"%Y-%m-%d_%H-%M-%S")
echo "Archiving creds files to backups/pre-upload-credentials-${DATE_STAMP}.tgz"
tar --exclude='./backups' -czvf backups/pre-upload-credentials-${DATE_STAMP}.tgz .

for i in "${FILES_LIST[@]}"
do
    echo
    echo "Encrypting file ${i}"
    gcloud kms encrypt --key=kubexray-ci --keyring=kubexray-ci --location=global --ciphertext-file="${i}".enc --plaintext-file="${i}"
    echo "Uploading file ${i}.enc to GCS bucket"
    gsutil cp ${i}.enc gs://"${GCS_BUCKET}"
done
echo

echo "Uploading file ingress-values.yaml to GCS bucket"
gsutil cp ingress-values.yaml gs://"${GCS_BUCKET}"
