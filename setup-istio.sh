set -e

# Download and install Istio
echo "Setting up Istio service mesh..."
curl -L https://istio.io/downloadIstio | sh -
cd istio-*
export PATH=$PWD/bin:$PATH

# Install Istio
cd ..
istioctl install -f k8s/istio.yml -y

# Inject envoy sidecar into all pods in the namespace
echo "Injecting Istio sidecars into containers..."
kubectl label namespace snl istio-injection=enabled

# Install Istio base components
echo "Installing Istio base components..."
kubectl get crd gateways.gateway.networking.k8s.io &> /dev/null || \
{ kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v1.3.0" | kubectl apply -f -; }

# Set up Istio ingress gateway
echo "Setting up Istio ingress gateway..."
kubectl apply -f k8s/gateway.yml

# expose the application
kubectl port-forward svc/snl-gateway-istio -n snl 5000:80 &

echo "Istio setup complete."

