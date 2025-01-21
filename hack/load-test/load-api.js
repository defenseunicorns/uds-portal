// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { K8s, kind } from "kubernetes-fluent-client";

const namespace = "load-api" + Math.floor(Math.random() * 1000);

async function createNamespace() {
  try {
    await K8s(kind.Namespace).Apply({
      metadata: {
        name: namespace,
      },
    });
    console.log(`Namespace ${namespace} created`);
  } catch (error) {
    console.error(`Error creating namespace ${namespace}:`, error);
  }
}

function createDaemonSet() {
  return K8s(kind.DaemonSet).Apply({
    metadata: {
      name: "load-test-daemonset",
      namespace: namespace,
    },
    spec: {
      selector: {
        matchLabels: {
          app: "load-test",
        },
      },
      template: {
        metadata: {
          labels: {
            app: "load-test",
          },
        },
        spec: {
          containers: [
            {
              name: "load-test",
              image: "nginx",
            },
          ],
        },
      },
    },
  });
}

function createStatefulSet() {
  return K8s(kind.StatefulSet).Apply({
    metadata: {
      name: "load-test-statefulset",
      namespace: namespace,
    },
    spec: {
      serviceName: "load-test",
      replicas: 3,
      selector: {
        matchLabels: {
          app: "load-test",
        },
      },
      template: {
        metadata: {
          labels: {
            app: "load-test",
          },
        },
        spec: {
          containers: [
            {
              name: "load-test",
              image: "nginx",
            },
          ],
        },
      },
    },
  });
}

function createCronJob() {
  return K8s(kind.CronJob).Apply({
    metadata: {
      name: "load-test-cronjob",
      namespace: namespace,
    },
    spec: {
      schedule: "*/1 * * * *",
      jobTemplate: {
        spec: {
          template: {
            spec: {
              containers: [
                {
                  name: "load-test",
                  image: "nginx",
                },
              ],
              restartPolicy: "Never",
            },
          },
        },
      },
    },
  });
}

function createJob() {
  return K8s(kind.Job).Apply({
    metadata: {
      name: "load-test-job",
      namespace: namespace,
    },
    spec: {
      template: {
        spec: {
          containers: [
            {
              name: "load-test",
              image: "nginx",
            },
          ],
          restartPolicy: "Never",
        },
      },
    },
  });
}

function createHPA() {
  return K8s(kind.HorizontalPodAutoscaler).Apply({
    metadata: {
      name: "load-test-hpa",
      namespace: namespace,
    },
    spec: {
      scaleTargetRef: {
        apiVersion: "apps/v1",
        kind: "Deployment",
        name: "load-test",
      },
      minReplicas: 1,
      maxReplicas: 10,
      metrics: [
        {
          type: "Resource",
          resource: {
            name: "cpu",
            target: {
              type: "Utilization",
              averageUtilization: 50, // Replaces targetCPUUtilizationPercentage
            },
          },
        },
      ],
    },
  });
}

function createLimitRange() {
  return K8s(kind.LimitRange).Apply({
    metadata: {
      name: "default-limits",
      namespace: namespace,
    },
    spec: {
      limits: [
        {
          type: "Container",
          max: {
            cpu: "2",
            memory: "1Gi",
          },
          min: {
            cpu: "200m",
            memory: "128Mi",
          },
          default: {
            cpu: "500m",
            memory: "256Mi",
          },
          defaultRequest: {
            cpu: "300m",
            memory: "128Mi",
          },
        },
      ],
    },
  });
}

function createResourceQuota() {
  return K8s(kind.ResourceQuota).Apply({
    metadata: {
      name: "resource-quota",
      namespace: namespace,
    },
    spec: {
      hard: {
        pods: "10",
        services: "5",
        replicationcontrollers: "5",
        secrets: "20",
        persistentvolumeclaims: "4",
        "requests.cpu": "4",
        "requests.memory": "16Gi",
        "limits.cpu": "8",
        "limits.memory": "32Gi",
      },
    },
  });
}

// can't make this resource because kubernetes-client-node says resource is not found in cluster, but can make manually with a yaml ....
function createRuntimeClass() {
  return K8s(kind.RuntimeClass).Apply({
    metadata: {
      name: "load-test-runtimeclass",
    },
    handler: "load-test-containers",
  });
}

const createResources = async () => {
  await createNamespace();
  await createDaemonSet();
  await createStatefulSet();
  await createCronJob();
  await createJob();
  await createHPA();
  await createLimitRange();
  await createResourceQuota();
  // await createRuntimeClass();
};

createResources()
  .then(() => {
    console.log("All resources created");
  })
  .catch((error) => {
    console.error("Error creating resources:", error);
  });
