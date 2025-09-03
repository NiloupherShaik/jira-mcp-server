# Modify this to point to your local context
KUBECONFIG_LOCAL="kind-ai" # "minikube" "docker-desktop"

# Dev
KUBECONFIG_DEV_EUWE1="arn:aws:eks:eu-west-1:616350704275:cluster/slr-cdo-dev-02-euwe1"

# Staging
KUBECONFIG_STG_EUWE1="arn:aws:eks:eu-west-1:086488590575:cluster/slr-cdo-stg-02-euwe1" # TODO: Update to the correct ARN
KUBECONFIG_STG_USWE2="arn:aws:eks:us-west-2:086488590575:cluster/slr-cdo-stg-02-uswe2" # TODO: Update to the correct ARN

# Production
KUBECONFIG_PRD_USWE2="arn:aws:eks:us-west-2:326553260002:cluster/slr-cdo-prd-02-uswe2" # TODO: Update to the correct ARN
KUBECONFIG_PRD_EUWE1="arn:aws:eks:eu-west-1:326553260002:cluster/slr-cdo-prd-02-euwe1" # TODO: Update to the correct ARN
KUBECONFIG_PRD_EUCE1="arn:aws:eks:eu-west-1:326553260002:cluster/slr-cdo-prd-02-euce1" # TODO: Update to the correct ARN
KUBECONFIG_PRD_APSE2="arn:aws:eks:eu-west-1:326553260002:cluster/slr-cdo-prd-02-apse2" # TODO: Update to the correct ARN
