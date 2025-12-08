#!/bin/bash
set -e

# -------------------------
# CONFIG
# -------------------------
CLUSTER_NAME=kind
NAMESPACE=snl

BACKEND_IMAGE=snl-backend:latest
FRONTEND_IMAGE=snl-frontend:latest

# -------------------------
# STEP 0: Check kind cluster exists
# -------------------------
if ! kind get clusters | grep -q "$CLUSTER_NAME"; then
    echo "Kind cluster '$CLUSTER_NAME' not found. Please create it first."
    exit 1
fi

# -------------------------
# STEP 1: Build Docker images
# -------------------------
echo "Building backend image..."
docker build -t $BACKEND_IMAGE ./backend

echo "Building frontend image..."
docker build -t $FRONTEND_IMAGE ./frontend

# -------------------------
# STEP 2: Load images into kind
# -------------------------
echo "Loading images into kind cluster '$CLUSTER_NAME'..."
kind load docker-image $BACKEND_IMAGE --name $CLUSTER_NAME
kind load docker-image $FRONTEND_IMAGE --name $CLUSTER_NAME

# -------------------------
# STEP 3: Apply Kubernetes manifests
# -------------------------
echo "Applying namespace..."
kubectl apply -f k8s/namespace.yml

echo "Applying secrets..."
kubectl apply -n $NAMESPACE -f k8s/secret.yml

echo "Deploying MongoDB..."
kubectl apply -n $NAMESPACE -f k8s/mongo/secret.yml
kubectl apply -n $NAMESPACE -f k8s/mongo/statefulset.yml
kubectl apply -n $NAMESPACE -f k8s/mongo/service.yml

echo "Deploying Redis..."
kubectl apply -n $NAMESPACE -f k8s/redis/deployment.yml
kubectl apply -n $NAMESPACE -f k8s/redis/service.yml

# wait for redis and mongo to be ready
kubectl wait --for=condition=ready pod -l app=redis -n $NAMESPACE
kubectl wait --for=condition=ready pod -l app=mongo -n $NAMESPACE

echo "Deploying backend..."
kubectl apply -n $NAMESPACE -f k8s/backend/deployment.yml
kubectl apply -n $NAMESPACE -f k8s/backend/service.yml
kubectl apply -n $NAMESPACE -f k8s/backend/hpa.yml

echo "Deploying frontend..."
kubectl apply -n $NAMESPACE -f k8s/frontend/deployment.yml
kubectl apply -n $NAMESPACE -f k8s/frontend/service.yml

# echo "Deploying ingress..."
# kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml
# kubectl apply -n $NAMESPACE -f k8s/ingress.yml

# -------------------------
# STEP 4: Done
# -------------------------
echo "✅ Deployment complete!"

# -------------------------
# STEP 5: Setup Istio service mesh
# -------------------------
./setup-istio.sh