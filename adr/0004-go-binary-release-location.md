# 4. UDS Runtime Go Binary Release Location

Date: 8 November 2024

## Status

Accepted

## Context

We are transitioning the UDS Runtime repository from public to private. Previously, we released the UDS Runtime Go binary as an artifact of our releases. Since Iron Bank uses this binary to create the registry1 UDS Runtime image, we need to find an alternative location for the binary to ensure Iron Bank can still access it, now that the repository will be private.

### Considerations

#### Access
The Go binary needs to be publicly accessible (read-only) and versioned

## Decision
We will create a public read-only s3 bucket in our `product_uds_ci_aws` account that will host the UDS Runtime Go binary. Our release pipeline will take the release tag, append the tag to the binary name (example: uds-runtime-linux-amd64-0.9.1), and upload the binary to the s3 bucket. The binary will be accessible at the following URL: https://uds-runtime.s3.amazonaws.com. This will allow us to make the UDS Runtime repo private while still providing Iron Bank with access to the UDS Runtime go binary.

## Consequences
A public version of the UDS Runtime binary will exist, despite it being closed source.
