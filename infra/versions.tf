terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.83.0"
    }
  }
  required_version = "~> 1.8.0"
}

provider "aws" {
  region = "us-west-2"
  default_tags {
    tags = {
      terraform           = true
      github_repository   = "defenseunicorns/uds-runtime"
      owner               = "uds-runtime"
      PermissionsBoundary = "uds_permissions_boundary"
    }
  }
}
