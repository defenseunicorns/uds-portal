# Release Playboook

This document describes the release process for UDS Runtime.

## Release CI

UDS Runtime uses release-please to manage releases. Configurations are located in [release-please-config.json](../release-please-config.json). Upon each merge into the `main` branch, release-please creates a PR that, when merged, will kick off the release CI (note that release-please also generates release notes based on [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) messages).

### Starting the Release CI

1. Find the release-please PR in Github
2. Add the PR to a Github milestone. Milestones are typically named after the release version (e.g. `v1.0.0`) and must be created before the PR can be added to it.
3. After the release-please PR's CI has passed, approve and merge the PR to kick off the release CI pipeline.

### Release Artifacts

The release CI pipeline will create the following artifacts:

#### Github Release

A Github Release will be created with the following assets:

- `uds-runtime-linux-amd64`
- `uds-runtime-linux-arm64`
- `uds-runtime-darwin-amd64`
- `uds-runtime-darwin-arm64`
- `checksums.txt`

#### Image

An internal-only UDS Runtime multi-arch Linux image located at: `oci://ghcr.io/defenseunicorns/uds-runtime:<version>`

#### Zarf Packages

Internal-only Zarf packages with the `unicorn` flavor will be created for both `amd64` and `arm64` architectures, located at: `oci://ghcr.io/defenseunicorns/packages/private/uds/uds-runtime:<version>-unicorn`

### S3 Artifacts

Because we need to support a public IronBank image and a canary environment, there are 2 S3 buckets that the release CI pushes artifacts to. These buckets are located in the `product_uds_ci_aws` AWS account (384851857288).

1. `runtime-canary`: This holds UDS Runtime Zarf packages that are used in the canary environments (currently `uds.run` only). We put the Zarf packages here because it's difficult to pull from internal GHCR repos in the Terraform userdata script that is used to create the canary EC2 instances.
2. `uds-runtime`: This holds the UDS Runtime Linux amd64 binaries used by the IronBank Dockerfile. This bucket is public and the files are publicly accessible.

## IronBank Release

The IronBank release process is a manual process that involves creating an MR to the [IronBank repo](https://repo1.dso.mil/dsop/opensource/defenseunicorns/uds/runtime), waiting for the IronBank team to merge the PR, and then creating the `registry1` flavor Zarf package. The steps are outlined below:

1. Create an issue in the IronBank repo indicating that a new version of UDS Runtime is ready to be released (extensive docs on how to create MRs are [located here](https://docs-ironbank.dso.mil/tutorials/submit-merge-request/)).
1. Create an MR in the IronBank repo (linked to the above issue) with the following changes:
   - Update the `hardening_manifest.yaml` file to the new version (should be in 3 places: `tags`, `args`, and `resources`)
   - Update the sha256 value (can be found in the `checksums.txt` file in the Github release or you can download the binary from S3 and run `sha256sum`)
1. Assign the `Status::Review` label to the MR

Once the IronBank team merges the MR into the `development` branch, the IronBank team will eventually merge it into `master`. Once this is done, the IronBank image will be available at: `registry1.dso.mil/ironbank/opensource/defenseunicorns/uds/runtime:<version>`

In the case that there are any CVEs caught by the IronBank pipeline, justifications can be added to the findings in [Vat](vat.dso.mil) which the IronBank team will review and once approved, the image will be released.

### Releasing the `registry1` Flavor Zarf Package

After the IronBank team merges the MR into `master`, the `registry1` flavor Zarf package can be created and released with the following steps (_todo_: automate step 1):

1. Create PR to update the `registry1` image version in the following files:
   - `zarf.yaml`
   - `hack/flavors/values/registry1-values.yaml`
1. Once the above PR has been merged, trigger the "Publish Registry1 Flavor of UDS Runtime" workflow in the Actions tab on Github.

## Release Tests

Once the Zarf packages have been released, you can perform post-release smoke tests by navigating to the Actions tab in Github and running the "Post Release Tests" workflow. This workflow will deploy the UDS Runtime Zarf packages to a test environment and run a series of tests to ensure the release is functioning as expected.

## Upgrading Demo Env

We maintain a demo environment at https://runtime-canary.uds.run that is used to showcase the latest features and changes. To upgrade the demo environment to the latest release, follow these steps:

1. Make changes to `.github/test-infra/terraform/setup.sh` to bump the `RUNTIME_VERSION` variable to the latest release version and the version specified in the UDS bundle YAML in that script
1. During a predetermined and agreed upon upgrade window, either:
   - Run the `Deploy Runtime + Core on AWS` workflow from the Actions tab in Github (this will completely destroy and rebuild the infra and will result in downtime)
   - Log into the EC2 instance containing the demo env (product CI AWS account), switch to the `ubuntu` user with `sudo su ubuntu`, download the UDS Runtime Zarf package from S3, and perform the deployment with UDS CLI (this will perform a rolling upgrade with no downtime)
