minikube image build . -t grpc-load-balancing-demo/client:v1.0.0-client-load-balancing

kubectl apply -f k8s/