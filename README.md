# UDS Runtime

<img align="right"  alt="zarf logo" src="ui/static/doug.svg"  height="256" />

[![UDS Documentation](https://img.shields.io/badge/docs-uds.defenseunicorns.com-775ba1)](https://uds.defenseunicorns.com/docs/)

UDS Runtime is the frontend for all things UDS, providing views and insights into your UDS cluster.

<br><br>

## Quickstart Deploy

### Pre-requisites

Recommended:

- [UDS-CLI](https://github.com/defenseunicorns/UDS-CLI#install)

If building locally:

- `Go >= 1.22.0`
- `Node >= v21.1.0`

### In Cluster

Assumes a K8s cluster is running and the appropriate K8s context has been selected

```bash
uds zarf package deploy oci://ghcr.io/defenseunicorns/packages/private/uds/uds-runtime:<tag> --confirm
```

\*_See [all tags](https://github.com/defenseunicorns/uds-runtime/pkgs/container/packages%2Fuds%2Fuds-runtime)_

#### Resource Requirements

When running in cluster, the Runtime pod will need more or less resources based on the number of resources in the cluster it will be watching. The [current defaults](./chart/values.yaml) work for a cluster running mainly UDS Core (about 44 pods). For running in larger clusters, UDS Core + SWF + Leapfrog for example (150+ pods), the `resource.limits.memory` will need to be more like `2Gi`.

### Locally (Out of Cluster)

1. clone this repo
1. compile: `uds run compile`
1. run: `./build/uds-runtime`

## Quickstart Development

For a full guide on developing for UDS Runtime, please read the [CONTRIBUTING.md](./CONTRIBUTING.md)

### To start the backend development server, run the following command:

**With UDS-CLI**

```bash
uds run dev-server
```

**Without UDS-CLI**

```bash
air
```

> **NOTE**: If you do not have air installed, you can find instructions for how to install at [here](https://github.com/air-verse/air)

### To start the frontend server, run the following command:

**With UDS-CLI**

```bash
uds run dev-ui
```

**Without UDS-CLI**

```bash
cd ui
npm ci
npm run dev
```

### UDS Security Hub Integration

To view the "Security/CVE Report" tab, devs will need to have the UDS Security Hub Sqlite database and metadata on their local filesystem. This can be done by running the following command from the root of the repo:

```bash
uds run setup:download-sechub-db
```

This will download the following files

```
./uds-security-hub.db
./artifacts/security-hub-metadata.json
```

Devs will need to set the following environment variable to be the root directory of the UDS Security Hub data. This value will most likely be `.`

```bash
SECURITY_HUB_DATA_PATH=.
```

## Nightly Releases

UDS Runtime publishes a canary release of latest changes every night tagged `nightly-unstable`

```bash
uds zarf package deploy oci://ghcr.io/defenseunicorns/packages/private/uds/uds-runtime:nightly-unstable
```

## Tech Stack

- Backend:

  - [Golang](https://go.dev/)
  - [Chi HTTP Router](https://github.com/go-chi/chi)
  - [K8s client-go](https://github.com/kubernetes/client-go)

- Frontend:

  - [Sveltekit](https://kit.svelte.dev/)
  - [Vite](https://vitejs.dev/)
  - [Typescript](https://typescriptlang.org/)
  - [TailwindCSS](https://tailwindcss.com/) ([Flowbite](https://flowbite.com/))
  - [Carbon Icons](https://www.carbondesignsystem.com/guidelines/icons/library)
  - [svelte-apexcharts](https://github.com/bn3t/svelte-apexcharts)

- Networking:

  - [Server Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
  - [REST API](https://restfulapi.net/)
  - [K8s Shared Informers](https://pkg.go.dev/k8s.io/client-go/informers)
