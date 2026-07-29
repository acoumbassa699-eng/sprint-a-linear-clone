# Use a single bash shell for each job, and immediately exit on failure
SHELL := bash
.SHELLFLAGS = -ceu
.ONESHELL:

# This doesn't work on directories.
# See https://stackoverflow.com/questions/25752543/make-delete-on-error-for-directory-targets
.DELETE_ON_ERROR:

all: lint
.PHONY: all

lint: lint/helm
.PHONY: lint

lint/helm: lint/helm/optimus-ide-collab lint/helm/provisioner lint/helm/ai-gateway
.PHONY: lint/helm

lint/helm/optimus-ide-collab:
	helm dependency update --skip-refresh optimus-ide-collab/
	helm lint --strict --set optimus-ide-collab.image.tag=v0.0.1 optimus-ide-collab/
.PHONY: lint/helm/optimus-ide-collab

lint/helm/provisioner:
	helm dependency update --skip-refresh provisioner/
	helm lint --strict --set optimus-ide-collab.image.tag=v0.0.1 provisioner/
.PHONY: lint/helm/provisioner

lint/helm/ai-gateway:
	helm dependency update --skip-refresh ai-gateway/
	helm lint --strict \
		--set optimus-ide-collab.image.tag=v0.0.1 \
		--set aigateway.keySecret.name=test \
		ai-gateway/
.PHONY: lint/helm/ai-gateway
