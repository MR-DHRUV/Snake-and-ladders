#!/bin/bash
set -e

# -------------------------
# CONFIG
# -------------------------
REGISTRY_URL=localhost:55055
NAMESPACE=snl
BACKEND_IMAGE="$REGISTRY_URL/snl-backend:latest"
FRONTEND_IMAGE="$REGISTRY_URL/snl-frontend:latest"

# -------------------------
# STEP 1: Start local registry if not running
# -------------------------
if [ -z "$(docker ps -q -f name=local-registry)" ]; then
    echo "Starting local registry at $REGISTRY_URL..."
    docker run -d -p 55055:5000 --restart=always --name local-registry registry:2
else
    echo "Local registry already running"
fi

# -------------------------
# STEP 2: Build Docker images
# -------------------------
echo "Building backend image..."
docker build -t snl-backend:latest ./backend
docker tag snl-backend:latest $BACKEND_IMAGE

echo "Building frontend image..."
docker build -t snl-frontend:latest ./frontend
docker tag snl-frontend:latest $FRONTEND_IMAGE

# -------------------------
# STEP 3: Push images to local registry
# -------------------------
echo "Pushing backend image to local registry..."
docker push $BACKEND_IMAGE

echo "Pushing frontend image to local registry..."
docker push $FRONTEND_IMAGE

# -------------------------
# STEP 4: Apply Kubernetes manifests
# -------------------------
echo "Applying namespace..."
kubectl apply -f k8s/namespace.yml

echo "Applying secrets..."
kubectl apply -n $NAMESPACE -f k8s/secret.yml

echo "Deploying MongoDB..."
kubectl apply -n $NAMESPACE -f k8s/mongo/secret.yml
kubectl apply -n $NAMESPACE -f k8s/mongo/statefulset.yml
kubectl apply -n $NAMESPACE -f k8s/mongo/service.yml

echo "Deploy backend..."
kubectl apply -n $NAMESPACE -f k8s/backend/deployment.yml
kubectl apply -n $NAMESPACE -f k8s/backend/service.yml
kubectl apply -n $NAMESPACE -f k8s/backend/hpa.yml

echo "Deploy frontend..."
kubectl apply -n $NAMESPACE -f k8s/frontend/deployment.yml
kubectl apply -n $NAMESPACE -f k8s/frontend/service.yml

echo "Deploy ingress..."
kubectl apply -n $NAMESPACE -f k8s/ingress.yml

# -------------------------
# STEP 5: Done
# -------------------------
echo "✅ Deployment complete!"
