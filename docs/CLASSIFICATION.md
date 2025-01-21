# Classification Configuration

Runtime can be configured to display a classification banner.

## Configure

To display the banner, which can appear as just a header or as header and footer, you need to set two environment variables:

- CLASSIFICATION_BANNER_ENABLED
- CLASSIFICATION_BANNER_TEXT
- CLASSIFICATION_BANNER_FOOTER

If you're running Runtime locally (for development), you can export these variables from your terminal:

```bash
export CLASSIFICATION_BANNER_ENABLED=true
export CLASSIFICATION_BANNER_FOOTER=false
export CLASSIFICATION_BANNER_TEXT="unclassified"
```

When deploying Runtime into a cluster, however, these variables are controlled by [Helm values](../chart/values.yaml)

```yaml
classificationBanner:
  enabled: false
  addFooter: false
  text: ""
```

The easiest way to configure these values is through bundle overrides (unless you're building the package and then can use the helm values file directly):

```yaml
kind: UDSBundle
metadata:
  name: runtime-demo-bundle
  version: 0.0.1

packages:
  - name: uds-runtime
    repository: ghcr.io/defenseunicorns/packages/private/uds/uds-runtime
    ref: <tag>-unicorn
    overrides:
      uds-runtime:
        uds-runtime:
          variables:
            - name: CLASSIFICATION_ENABLED
              description: enable classification banner
              path: classificationBanner.enabled
            - name: CLASSIFICATION_TEXT
              description: text to display in classification banner
              path: classificationBanner.text
            - name: CLASSIFICATION_FOOTER
              description: display footer too
              path: classificationBanner.addFooter
```

You can either add default values in the bundle variable definitions, or you can use any of the established ways to set values for bundle variables via UDS CLI.

## Allowed Values

Only a set of known classifications are accepted in the helm chart. A type definition, checked by Helm at deploy time, exists in the [values.schema.json](../chart/values.schema.json) file.

Current classifications that we're looking for:

```console
  'UNCLASSIFIED',
  'CUI',
  'CONFIDENTIAL',
  'SECRET',
  'TOP SECRET',
  'TOP SECRET//SCI',
  'UNKNOWN'
```

> Note  
> The text does not have to be uppercase.

## Specifications

The banner follows the specs from https://www.astrouxds.com/components/classification-markings/
