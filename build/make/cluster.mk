cluster_dir=$(if $(project_root),"$(project_root)/build/make",".")

.PHONY: download-signoz-collector
download-signoz-collector:
	@echo "Authenticating to ECR and loading SigNoz collector image..."
	@aws ecr get-login-password --region eu-west-1 | docker login --username AWS --password-stdin 263262308774.dkr.ecr.eu-west-1.amazonaws.com
	@docker pull 263262308774.dkr.ecr.eu-west-1.amazonaws.com/ns/mirror/prd/signoz/signoz-otel-collector:0.111.27
	@kind load docker-image 263262308774.dkr.ecr.eu-west-1.amazonaws.com/ns/mirror/prd/signoz/signoz-otel-collector:0.111.27 -n ai

.PHONY: clean-setup-local-kind-cluster
clean-setup-local-kind-cluster: delete-ai-kind-cluster
clean-setup-local-kind-cluster: setup-local-kind-cluster

.PHONY: setup-local-kind-cluster
setup-local-kind-cluster: create-ai-kind-cluster
setup-local-kind-cluster: setup-local-cluster

.PHONY: delete-ai-kind-cluster
delete-ai-kind-cluster:
	@kind delete clusters ai

create-ai-kind-cluster:
	@kind create cluster --config $(cluster_dir)/../../deployments/k8s/overlays/local/infra/kind-config.yaml
	@kubectx kind-ai

run-setup-scripts:
	@kubectx $(KUBECONFIG_LOCAL)
#	@$(MAKE) download-signoz-collector
	@cd $(cluster_dir)/../../deployments/k8s/overlays/local/infra && \
    		sh -c './0010_install-ingress-controller.sh --version 1.24' && \
    		sh -c './0020_install-cert-manager.sh --version 1.16.2' && \
    		sh -c './0030_install-sealed-secrets.sh --version 0.27.2' && \
    		sh -c './0035_install-self-signed-issuer.sh' && \
    		sh -c './0040_install-signoz.sh' && \
			sh -c './0105_create_namespace.sh'

.PHONY: setup-local-cluster
setup-local-cluster: run-setup-scripts

run-local-kind-load-balancer: install-local-kind-load-balancer
	if [ "$(shell kubectl config current-context)" = "kind-ai" ]; then sudo cloud-provider-kind ; fi

install-local-kind-load-balancer:
ifeq (, $(shell which cloud-provider-kind))
	@echo "Installing cloud-provider-kind"
	@go install sigs.k8s.io/cloud-provider-kind@latest
	@echo "Installing cloud-provider-kind from ~/go/bin/cloud-provider-kind to /usr/local/bin"
	sudo install ~/go/bin/cloud-provider-kind /usr/local/bin
endif

.PHONY: get-local-kind-status
get-local-kind-status:
	docker inspect ai-control-plane | egrep -iC2 'Status|ExitCode|RestartPolicy'
