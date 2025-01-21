locals {
  bucket_name = "${var.bucket_prefix}-${var.region}-tfstate"
}

resource "aws_kms_key" "objects" {
  # checkov:skip=CKV2_AWS_64: "Ensure KMS key Policy is defined"
  enable_key_rotation     = true
  description             = "KMS key is used to encrypt bucket objects (this one is for the tf state bucket)"
  deletion_window_in_days = 7
}

module "state_bucket" {
  # checkov:skip=CKV2_AWS_64: "Ensure KMS key Policy is defined"
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "4.2.2"

  bucket = local.bucket_name

  server_side_encryption_configuration = {
    rule = {
      apply_server_side_encryption_by_default = {
        kms_master_key_id = aws_kms_key.objects.arn
        sse_algorithm     = "aws:kms"
      }
    }
  }

  versioning = {
    status     = true
    mfa_delete = false
  }
}

module "artifact_bucket" {
  # checkov:skip=CKV2_AWS_64: "Ensure KMS key Policy is defined"
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "4.2.2"

  bucket = "runtime-nightly-artifacts"

  server_side_encryption_configuration = {
    rule = {
      apply_server_side_encryption_by_default = {
        kms_master_key_id = aws_kms_key.objects.arn
        sse_algorithm     = "aws:kms"
      }
    }
  }
}

resource "aws_kms_key" "dynamodb" {
  # checkov:skip=CKV2_AWS_64: "Ensure KMS key Policy is defined"
  enable_key_rotation     = true
  description             = "KMS key used to encrypt DynamoDB Table"
  deletion_window_in_days = 7
}

module "lock_table" {
  source  = "terraform-aws-modules/dynamodb-table/aws"
  version = "4.2.0"

  name     = "uds-state-lock"
  hash_key = "LockID"

  attributes = [
    {
      name = "LockID"
      type = "S"
    }
  ]

  point_in_time_recovery_enabled = true

  server_side_encryption_enabled     = true
  server_side_encryption_kms_key_arn = aws_kms_key.dynamodb.arn
}
