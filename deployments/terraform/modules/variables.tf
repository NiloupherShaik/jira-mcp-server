variable "fusion-app-name" {
  description = "Name of the app being deployed"
  default = "MCP Server template"
}

variable "aws-region" {
  description = "AWS region in which to apply"
}

variable "aws-account-id" {
  description = "AWS account in which to apply"
}

variable "irsa-oidc-provider" {
  description = ""
}
