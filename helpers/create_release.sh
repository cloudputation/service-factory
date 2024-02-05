#!/bin/bash
RELEASE_TAG=$(cat VERSION)
git tag -a v${RELEASE_TAG} -m "Creating release for version: ${RELEASE_TAG}"
git push origin v${RELEASE_TAG}
