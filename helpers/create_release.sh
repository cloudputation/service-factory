#!/bin/bash
RELEASE_TAG=$(cat VERSION)
GITHUB_SF_TOKEN="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"
GITHUB_REMOTE_URL="https://${GITHUB_SF_TOKEN}@github.com/cloudputation/service-factory.git"
GITLAB_REMOTE_URL="git@gitlab.com:cloudputation/service-factory.git"
git remote remove origin
git remote add origin $GITHUB_REMOTE_URL
git tag -a v${RELEASE_TAG} -m "Creating release for version: ${RELEASE_TAG}"
git push origin v${RELEASE_TAG}
git remote remove origin
git remote add origin $GITLAB_REMOTE_URL
