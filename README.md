# UDS Portal

<img align="right"  alt="Unicorn Delivery Service logo" src="ui/static/uds.svg"  height="256" />

[![UDS Documentation](https://img.shields.io/badge/docs-uds.defenseunicorns.com-775ba1)](https://uds.defenseunicorns.com/docs/)

UDS Portal is the landing page for all UDS users, it serves as a single point of discovery for all UDS-deployed applications.

<br><br>

## Quickstart Deploy

> [!NOTE]  
> UDS Portal is in early alpha, follow the instructions below to deploy a prototype of the app.

To deploy the UDS Portal in your UDS cluster, run:

```
uds run test:e2e-setup
```

After a successful deployment, you can access the UDS Portal at https://apps.uds.dev

### Pre-requisites

Recommended:

- [UDS-CLI](https://github.com/defenseunicorns/UDS-CLI#install)

If building locally:

- `Go >= 1.22.0`
- `Node >= v21.1.0`

## Authentication Model

When deployed in-cluster, UDS Portal runs behind AuthService (from UDS Core).

- AuthService is the trust boundary for identity: it mints and validates JWTs.
- UDS Portal does not perform independent JWT signature verification.
- UDS Portal consumes identity/group claims that AuthService has already validated.
