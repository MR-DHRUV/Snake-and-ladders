#!/bin/bash
set -e

echo "Deploying observability stack..."

# Apply Jaeger (Tracing)
echo "Applying Jaeger..."
kubectl apply -f k8s/observability/jaeger.yml

# Apply Prometheus (Metrics, required for Kiali)
echo "Applying Prometheus..."
kubectl apply -f k8s/observability/prometheus.yml

# Apply Kiali (Dashboard)
echo "Applying Kiali..."
kubectl apply -f k8s/observability/kiali.yml

# Apply Telemetry configuration (Enable Jaeger provider)
echo "Applying Telemetry configuration..."
kubectl apply -f k8s/observability/telemetry.yml

echo "Observability setup complete!"
echo "--------------------------------------------------"
echo "To access Kiali Dashboard:"
echo "1. Run: kubectl port-forward svc/kiali -n istio-system 20001:20001"
echo "2. Open: http://localhost:20001/kiali"
echo "--------------------------------------------------"
