data "aws_partition" "current" {}
data "aws_caller_identity" "current" {}

locals {
  permissions_boundary_policy_arn = "arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:policy/${var.permissions_boundary_name}"
}

# Creates an IAM role that allows the uds-runtime GitHub repo to assume it
module "aws_oidc_github" {
  source  = "unfunco/oidc-github/aws"
  version = "v1.8.1"

  attach_admin_policy           = true
  create_oidc_provider          = true
  iam_role_permissions_boundary = local.permissions_boundary_policy_arn
  github_repositories           = ["defenseunicorns/uds-runtime"]
  iam_role_name                 = "GitHub-OIDC-Role"
}
