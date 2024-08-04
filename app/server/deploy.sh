minikube image build . -t grpc-load-balancing-demo/server:v1.0.0-client-load-balancing

kubectl apply -f k8s/