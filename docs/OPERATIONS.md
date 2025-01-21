# UDS Runtime Operations Runbook

This document describes the operations runbook for the UDS Runtime [demo environment](https://demo.uds.run).

** If this is an emergency, follow the steps in [Deploying the Demo Environment](#redeploying-the-demo-environment) to quickly redeploy the demo environment. **

### Table of Contents

- [Monitoring, Observability, and SLA](#monitoring-observability-and-sla)
- [Remediation Steps](#remediation-steps)
- [Accessing the Demo Environment](#accessing-the-demo-environment)
- [Deploys And Upgrades](#deploys-and-upgrades)
  - [Blue / Green Deploy (recommended)](#blue--green-deploy-recommended)
  - [Redeploying the Demo Environment](#redeploying-the-demo-environment)

## Monitoring, Observability, and SLA

The demo env is expected to be up during business hours Monday through Friday. The monitoring and observability of the demo env is performed by [Gatus](https://defenseunicorns-dev.gatus.io/). If Gatus detects that the demo env is down, it will send an alert to the `uds-runtime-dev` Slack channel.

If a team member receives notification that the demo env is down or observes a Gatus alert, the team member will follow the steps in the [Remediation Steps](#remediation-steps) section to investigate and remediate the issue.

## Remediation Steps

If you have received a notification that the demo env is down, either from Gatus or another source, follow the steps below to investigate and remediate the issue:

1. Navigate to https://demo.uds.run to verify that the demo environment is down (use Google SSO to log in).

### If the environment is up:

1. Log into [Gatus](https://defenseunicorns-dev.gatus.io/) and check the status of the demo environment, it could be a false positive or the environment self-healed (ie. the Runtime pod came back up)
1. Investigate the cluster and pod status in the demo environment's K8s cluster. You can do this by following the steps in the [Accessing the Demo Environment](#accessing-the-demo-environment) section or just using the existing UDS Runtime deployment located at https://demo.uds.run.
1. If any issues or pod failures are found, document the behavior and any error messages in a Github Issue.

### If the environment is down:

1. If this is an emergency, follow the steps in the [Redeploying the Demo Environment](#redeploying-the-demo-environment) section.
1. Otherwise, follow the steps in the [Accessing the Demo Environment](#accessing-the-demo-environment) section to access the demo environment. The goal is to investigate the root cause of the issue and gather as much information as possible so it can fixed permanently.
1. Check the status of the demo environment's K8s cluster (note UDS CLI is installed in the environment)
1. Gather as much information as possible about why the environment became inaccessible
1. Create a Github Issue with the details of the issue and the steps taken to investigate the issue.
1. You can either attempt to re-deploy portions of the environment, or (recommended) simply recreate the environment using the [Deploy Runtime + Core on AWS](../workflows/deploy-demo-env.yaml) workflow, instructions for doing this are located in the [Redeploying the Demo Environmen](#redeploying-the-demo-environment) section.

## Accessing the Demo Environment

We recommended using the AWS web console to access the demo environment. The demo environment is hosted in the `product_uds_ci_aws` (384851857288) AWS account, and the demo env is is hosted on an EC2 instance named `runtime-demo-<some hash>` in `us-west-2`. To interact with the demo environment:

1. From the AWS web console, navigate to the EC2 instances page.
1. Find the instance named `runtime-demo-<some hash>` (ensure you are in `us-west-2`).
1. Select the instance and click the `Connect` button.
1. Ensure you are in the `Session Manager` tab.
1. Click the `Connect` button.
1. You are now connected to the demo environment and logged in as the `ssm-user`, you will need to change to the `ubuntu` user to interact with the demo environment, do this with the following command: `sudo su ubuntu`.
1. You can now interact with the demo environment as needed.

### Helpful Tips

- Logs for the cloud-init script that created the environment (including the bundle deployment) are located at `/var/log/cloud-init-output.log`, you can `tail -f` this file to monitor the deployment progress.
- The UDS Core bundle and UDS Runtime Zarf package artifacts are located in `/tmp`.

## Deploys And Upgrades

### Blue / Green Deploy (recommended)

To upgrade UDS Runtime demo environment with minimal to zero downtime, we perform a blue / green deploy of UDS Runtime:

1. Follow the steps outlined in [DEPLOYMENT_INFRASTRUCTURE.md](./DEPLOYMENT_INFRASTRUCTURE.md#updating-demo-deployment)

### Redeploying the Demo Environment

Follow the steps below to redeploy the [demo environment's IaC](../.github/test-infra/README.md), note that this will destroy the current demo environment and create a new one, resulting in downtime:

1. Navigate to the [Actions tab](https://github.com/defenseunicorns/uds-runtime/actions) on the UDS Runtime Github repo.
1. Under `All workflows` click on `Deploy Runtime + Core on AWS` (it should be pinned)
1. Run the workflow by clicking the green `Run workflow` button in the `Run workflow` dropdown.

Note that even if the workflow completes succesfully, the demo environment will not be fully up until the [cloud init script](../.github/test-infra/terraform/setup.sh) finishes executing, which can take up to 15 minutes. You can monitor the rollout of the demo env by following the steps outlined in the [Accessing the Demo Environment](#accessing-the-demo-environment) section.
