echo "Starting Deployment"

echo "Building Docker Image"
docker build -t snake-and-ladders:1.0 . 

echo "Pushing Docker Image to Docker Hub"
docker tag snake-and-ladders:1.0 dhruvgupta742/snake-and-ladders:latest
docker push dhruvgupta742/snake-and-ladders:latest

echo "Applying Kubernetes Manifests"
kubectl apply -f k8s/namespace.yml
kubectl apply -f k8s/secret.yml

kubectl apply -f k8s/deployment.yml
kubectl apply -f k8s/hpa.yml

kubectl apply -f k8s/service.yml
kubectl apply -f k8s/ingress.yml

echo "Deployed Successfully"