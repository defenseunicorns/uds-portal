variable "region" {
  description = "AWS Region"
  type        = string
  default     = "us-west-2"
}

variable "bucket_prefix" {
  description = "S3 Bucket Prefix"
  type        = string
  default     = "uds-aws-runtime-nightly"
}

variable "permissions_boundary_name" {
  description = "The name of the policy that is used to set the permissions boundary for the role."
  type        = string
  default     = "uds_permissions_boundary"
}
