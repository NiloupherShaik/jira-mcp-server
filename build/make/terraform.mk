# DEPLOYMENT TERRAFORM

.PHONY: deploy-dev-k8s
deploy-dev-k8s: ENVIRONMENT=dev
deploy-dev-k8s: REGION=eu-west-1
deploy-dev-k8s: .set-kustomize-images
deploy-dev-k8s: .run-terraform

.PHONY: deploy-stg-k8s
deploy-stg-k8s: ENVIRONMENT=stg
deploy-stg-k8s: REGION=eu-west-1
deploy-stg-k8s: .set-kustomize-images
deploy-stg-k8s: .run-terraform

.run-terraform:
	cd deployments/terraform/modules && \
	tfswitch && \
	terraform init -backend-config="../environments/${ENVIRONMENT}/${REGION}/backend-config.tf" && \
	terraform apply -auto-approve -var-file="../environments/${ENVIRONMENT}/${REGION}/terraform.tfvars"

# PLAN TERRAFORM

.PHONY: plan-dev-k8s
plan-dev-k8s: ENVIRONMENT=dev
plan-dev-k8s: REGION=eu-west-1
plan-dev-k8s: .set-kustomize-images
plan-dev-k8s: .plan-terraform

.PHONY: plan-stg-k8s
plan-stg-k8s: ENVIRONMENT=stg
plan-stg-k8s: REGION=eu-west-1
plan-stg-k8s: .set-kustomize-images
plan-stg-k8s: .plan-terraform

.plan-terraform:
	cd deployments/terraform/modules && \
	tfswitch && \
	terraform init -backend-config="../environments/${ENVIRONMENT}/${REGION}/backend-config.tf" && \
	terraform plan -var-file="../environments/${ENVIRONMENT}/${REGION}/terraform.tfvars"

.set-kustomize-images:
	cd deployments/k8s/overlays/$(ENVIRONMENT)/base && \
	kustomize edit set image \
		${APP_NAME}=$(API_SERVER_DOCKER_REPO):$(VERSION)
