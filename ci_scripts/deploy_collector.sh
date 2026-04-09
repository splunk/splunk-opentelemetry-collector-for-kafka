#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CHART_DIR="${SCRIPT_DIR}/../helm-chart/splunk-opentelemetry-collector-for-kafka"
VALUES_FILE="${SCRIPT_DIR}/ci_values.yaml"

CI_SPLUNK_HEC_TOKEN="${CI_SPLUNK_HEC_TOKEN:-00000000-0000-0000-0000-0000000000000}"

# Resolve the node IP so pods (without hostNetwork) can reach
# Kafka and Splunk pods (which use hostNetwork).
NODE_IP=$(sudo microk8s kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')
if [ -z "${NODE_IP}" ]; then
  echo "ERROR: Could not resolve node IP"
  exit 1
fi
echo "Resolved node IP: ${NODE_IP}"

# Prepare values with substituted placeholders (pipe delimiter avoids issues with / in values)
TEMP_VALUES=$(mktemp /tmp/ci_values.XXXXXX.yaml)
cp "${VALUES_FILE}" "${TEMP_VALUES}"
sed -i "s|__NODE_IP__|${NODE_IP}|g" "${TEMP_VALUES}"
sed -i "s|__SPLUNK_HEC_TOKEN__|${CI_SPLUNK_HEC_TOKEN}|g" "${TEMP_VALUES}"

echo "--- Rendered CI values (token redacted) ---"
sed 's/token: .*/token: "***REDACTED***"/' "${TEMP_VALUES}"
echo "--- End of CI values ---"

# Clean up any previous deployment
if sudo microk8s helm3 ls --short | grep -q .; then
  echo "Cleaning previous deployments..."
  sudo microk8s helm3 ls --short | xargs -r sudo microk8s helm3 uninstall || true
fi

echo "Deploying SOC4Kafka chart..."
sudo microk8s helm3 install ci-soc4kafka -f "${TEMP_VALUES}" "${CHART_DIR}"

echo "Waiting for SOC4Kafka deployment to be ready..."
DEP_NAME=$(sudo microk8s kubectl get deployment -l app.kubernetes.io/instance=ci-soc4kafka -o jsonpath='{.items[0].metadata.name}')
sudo microk8s kubectl rollout status "deployment/${DEP_NAME}" --timeout=180s

echo "SOC4Kafka deployment is ready."
sudo microk8s kubectl get pods -l app.kubernetes.io/instance=ci-soc4kafka

rm -f "${TEMP_VALUES}"
