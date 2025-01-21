# Runtime Ephemeral Infrastructure

The IaC found in [.github/test-infra/](../.github/test-infra/) is used by the [Deploy Runtime + Core on AWS](../.github//workflows/deploy-demo-env.yaml) and [Nightly Deploy](../.github//workflows/deploy-nightly-env.yaml) workflows, via [uds tasks](../.github/test-infra/tasks/infra.yaml), to create the UDS Runtime [demo env](https://demo.uds.run) and [nightly testing env](https://runtime.exploding.boats).

## How it Works

When the workflow kicks off, it will `tofu init` using the backend variables defined in the workflow, then destroy the currently running EC2 instance and related infra. After removing the old infra, it will create a new EC2 instance, that on startup will run the [setup.tpl](../.github/test-infra/terraform/setup.tpl)

## Custom AMI

The EC2 instance is created with a custom AMI. We use `packer` to define the AMI in [runtime.pkr.hcl](../.github/test-infra/packer/runtime.pkr.hcl) and build / push it to our AWS accounts.

To update the AMI, use the [Build and Push AMI](../.github/workflows/build-and-push-ami.yaml) workflow. The workflow is set up to push the ami to both the Product CI and Nightly Runtime Staging accounts.

> **NOTE**  
> Please delete old instances of the AMI

## Debug with SSM

The ec2 instance has been configured with SSM for debugging running clusters without needing SSH. To start an SSM session from the aws console:

`Systems Manager` > click `Session Manager` under `Node Management` > click `start session` > select `runtime-*` > click `start session`

You will be logged in as the ssm-user. To interact with the cluster you'll need to `sudo su ubuntu` and `cd ~`.

> **NOTE**  
> You can also follow the alternative steps in [OPERATIONS.md](./OPERATIONS.md#accessing-the-demo-environment)

## Demo Deployment

The demo instance, accessible at `https://demo.uds.run`, is the main Runtime deployment used for external demos. While it is not a PROD environment, we do our best to keep it up at all times and do things, like updates, as if it were a PROD instance (e.g. rolling updates instead of tear down and re-deploy).

> **NOTE**  
> The demo instance is currently deployed into the Product CI AWS account. That will most likely change to a separate account managed by Spacelift some time in the near future.

### Updating Demo Deployment

Because there is no SSH to the demo ec2 instance, for executing updates say via Github actions, the current approach is to login to the instance via SSM (see [Debug with SSM](#debug-with-ssm)) and then manually kick off a rolling update. This involves:

- creating a "mini" bundle, containing only Runtime and the gateway / host overrides
- updating the Runtime version in "mini" bundle
- copying the new Runtime release artifact from the `runtime-canary` S3 bucket to the /tmp directory
- creating and deploying the "mini" bundle

```yaml
kind: UDSBundle
metadata:
  name: update-runtime-bundle
  description: for deploying a new version of Runtime
  version: $RUNTIME_VERSION

packages:
  - name: uds-runtime
    path: /tmp
    ref: $RUNTIME_VERSION
    overrides:
      uds-runtime:
        uds-runtime:
          variables:
            - name: GATEWAY
              description: The gateway to use for the runtime
              path: package.gateway
              default: tenant
            - name: HOST
              description: The host to use for the runtime
              path: package.host
              default: demo
```

```console
sudo aws s3 cp s3://runtime-canary/zarf-package-uds-runtime-<arch>-<version>.tar.gz /tmp/
uds create uds-bundle.yaml --confirm
uds deploy uds-bundle-<name-of-bundle>-*.tar.zst --set DOMAIN=uds.run --confirm
```

> **NOTE**  
> We use a mini bundle because if you update the original bundle found in /tmp and re-deploy it, you will destroy the cluster entirely (re-deploy of uds-k3d) and will also need to re-download the TLS cert and key and export those variables.

### Troubleshooting Updates

There are some known issues with doing rolling updates on the demo ec2 instance. For guidance in troubleshooting, please refer to the [Operations Runbook](./OPERATIONS.md)

## Nightly Deployment

The nightly instance, accessible at `https://runtime.exploding.boats`, is our ephemeral testing deployment for the `nightly-unstable` release. This instance allows for faster feedback on new code pushed to `main` daily. This instance is not meant to be long-lived, meaning that every night the instance is torn down and re-created with the latest nightly release. The AWS account this instance is deployed to is `product/staging/uds-runtime-nightly` managed by Spacelift. The IaC for configuring this account is found first in the `it-ops-infra` repo and then in the [`infra`](../infra/) directory in this repo.

### Why Are There Two Infra Directories?

The reason why the AWS account infra is separated from the ec2 instance deployment infra, is because Spacelift automatically runs `plan` and `apply` on commits to `main` in a single directory. That means that if we placed all of the infra in the same directory, Spacelift would only deploy our ec2 instance when we first commit the code to that directory. However, we want to deploy that instance on release cadences (official and nightly). So by keeping them separate we keep Spacelift from deploying our ec2 instance and we continue to control that aspect by using Github actions workflows.

> **WARNING**  
> If you are doing any kind of work in the `infra` directory, please do not remove the EIP resource. The IP address this resource has created has been attached to DNS records for publicly accessing `runtime.exploding.boats`

### Debugging

The nightly ec2 instance is the exact same as the demo ec2 instance, just deploying separate releases. This means that if you need to troubleshoot the deployment you can use SSM to get a terminal session on the instance.
