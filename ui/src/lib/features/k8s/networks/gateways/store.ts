// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { ResourceStore, transformResource } from '$features/k8s/store'
import { type ColumnWrapper, type CommonRow, type ResourceStoreInterface } from '$features/k8s/types'
import type { VirtualService } from 'uds-core-types/src/pepr/operator/crd/generated/istio/virtualservice-v1beta1'

export interface Row extends CommonRow {
  hosts: {
    list: string[]
  }
}

type ServerType = {
  hosts: string[]
}

interface Resource extends VirtualService {
  spec: VirtualService['spec'] & {
    servers: ServerType[]
  }
}

export type Columns = ColumnWrapper<Row>

export function createStore(): ResourceStoreInterface<Resource, Row> {
  const url = `/api/v1/resources/networks/gateways?fields=.metadata,.spec.servers`

  const transform = transformResource<Resource, Row>((r) => {
    const hosts = r.spec?.servers?.map((gateway) => gateway.hosts[0])

    return {
      hosts: {
        list: Array.from(new Set(hosts)).sort(),
      },
    }
  })

  const store = new ResourceStore<Resource, Row>(url, transform, 'namespace')

  return {
    ...store,
    start: store.start.bind(store),
    sortByKey: store.sortByKey.bind(store),
  }
}
