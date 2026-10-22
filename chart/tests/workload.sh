#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/../.."
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT

workload_values=chart/tests/workload-values.yaml

render() {
  uds zarf tools helm template uds-portal chart --namespace uds-portal \
    --show-only templates/deployment.yaml "$@"
}

check() {
  uds zarf tools yq --exit-status "$1" "$2" > /dev/null
}

compare() {
  diff -u \
    <(uds zarf tools yq -o=json "$1 | sort_keys(..)" "$2") \
    <(uds zarf tools yq -o=json "$3 | sort_keys(..)" "$4")
}

uds zarf tools helm lint chart
render > "$test_dir/default.yaml"
for key in nodeSelector affinity tolerations; do
  check ".spec.template.spec | has(\"$key\") == false" "$test_dir/default.yaml"
done
check '.spec.template.spec.containers[0].env | length == 6' "$test_dir/default.yaml"

uds zarf tools helm lint chart --values "$workload_values"
render --values "$workload_values" > "$test_dir/custom.yaml"
for key in nodeSelector affinity tolerations; do
  compare ".spec.template.spec.$key" "$test_dir/custom.yaml" ".$key" "$workload_values"
done
compare '.spec.template.spec.containers[0].env[6:]' "$test_dir/custom.yaml" \
  '.extraEnvs' "$workload_values"
compare '.spec.template.spec.containers[0].env[0:6]' "$test_dir/custom.yaml" \
  '.spec.template.spec.containers[0].env' "$test_dir/default.yaml"

# Reject incorrect value types before producing an invalid workload.
for entry in 'nodeSelector=[]' 'nodeSelector={"nodepool":42}' 'affinity=[]' \
  'tolerations={}' 'tolerations=["invalid"]' 'extraEnvs={}' 'extraEnvs=["invalid"]'; do
  if render --set-json "$entry" > "$test_dir/invalid.log" 2>&1; then
    echo "Expected schema validation to reject $entry" >&2
    exit 1
  fi
  if ! grep -q "${entry%%=*}" "$test_dir/invalid.log"; then
    cat "$test_dir/invalid.log" >&2
    exit 1
  fi
done

echo "Chart workload tests passed."
