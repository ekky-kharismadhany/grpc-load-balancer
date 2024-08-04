# Explore: GRPC Load Balancing

This repo holds my exploration on load balancing a grpc system. Currently, my team has one backend application that have high traffic (100k> rps from 9 pods). I'm not satisfied with current solution with the use of kubernetes' headless service. My team still see that there is an imbalance between the pods.

## Setup

I create two simple application to help me demonstrate that the current solution is not working properly. These two app is a client and server model. Each application has an UUID that I used as indentifier. The client has RESTful API endpoint, which upon called, will also make a gRPC call to the server application. This is the gRPC and RESTful contract.

```json
// RESTful contract
  "detail": {
    "from: <server's_app_id>": <number_of_call>
  }
```

```proto
// protocol buffer contract
message Echo {
    string message = 1;
}

service EchoServer {
    rpc CallEcho(Echo) returns(Echo); 
}
```

This is the deployment setup inside kubernetes.
![alt text](image/initial-k8s-deployment-setup.png)

Basically, the client app will count how many variance of the server app that serve its request. With this, we can conclude whether the request is well balanced or not.

## Current Result

When I test the response of the client app. I found this response.

```json
{
  "detail": {
    "from: 2a63ca4d-193f-4ec8-b9d3-49a39548bb95 replied by: 4bcbc60a-5700-4313-8772-c2a3ceab0513": 1000
  }
}
```

This shows that out of 10 server app's pod, only one serving the traffic.

## Client Side Load Balancing

Client side load balancing is performed with client library. In this docs, it will performed with `round_robin` algorithm. Here is how it works:

1. On start up, the client application will issue a name resolution to the DNS. Then, the DNS will return a list of host and a service config.
2. After the client receives the server's host list. It will override the service config as the default DNS system in k8s does not return service config.

The 2 step is known because of this documentation in code WithDefaultServiceConfig returns a DialOption that configures the default service config, which will be used in cases where `WithDisableServiceConfig` is also used, or the name resolver does not provide a service config or provides an invalid service config. The parameter s is the JSON representation of the default service config`

### Is there any code changes?

Yes, there is small code change to the application code. We need to add this line:

```go
...
 opts := []grpc.DialOption{
  insecureCred,
  grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`),
 }
 conn, err := grpc.NewClient(config.serverHost, opts...)
...
```

### Result

Here is the client result with this curl request:

```bash
curl localhost:9090/echo?iter=1000 | jq
```

```json
{
  "detail": {
    "from: 5aa23e94-eafe-4928-b696-5f2efb22f290 replied by: 2077dd21-6440-4787-a8b0-8b9a2c87fade": 100,
    "from: 5aa23e94-eafe-4928-b696-5f2efb22f290 replied by: 24b7647d-7bf7-48a8-965f-cfd3575444b6": 100,
    "from: 5aa23e94-eafe-4928-b696-5f2efb22f290 replied by: 2c2cecfd-28ef-4490-b159-57624356dbfc": 100,
    "from: 5aa23e94-eafe-4928-b696-5f2efb22f290 replied by: 388afcc2-890b-4261-b8ce-7ccab8668ccd": 100,
    "from: 5aa23e94-eafe-4928-b696-5f2efb22f290 replied by: 4d61ae83-a384-46c8-adc9-8b00803f8134": 100,
    "from: 5aa23e94-eafe-4928-b696-5f2efb22f290 replied by: 64b821be-d29c-4fe4-94b5-7b33a4b5558d": 100,
    "from: 5aa23e94-eafe-4928-b696-5f2efb22f290 replied by: 8953f032-8a23-4561-bb3f-34e378a0a274": 100,
    "from: 5aa23e94-eafe-4928-b696-5f2efb22f290 replied by: 912fb514-1966-44dd-a910-728606e78a23": 100,
    "from: 5aa23e94-eafe-4928-b696-5f2efb22f290 replied by: abd1cb27-4b17-4ca2-8849-bd8ea0bfe0f7": 100,
    "from: 5aa23e94-eafe-4928-b696-5f2efb22f290 replied by: b0ee357b-e051-404f-87ae-8f53b4c62f30": 100
  }
}
```

## Scaling Elasticity Issue

I notice that the traffic load is not balanced directly after scale up event. This is because the resolver needs some time to retrieve newly created host and create a subroutine for the new host. This issue does not appear when we do scaling down the pods, as the subroutine's connectivity status is set to `TRANSIENT_FAILURE`, thus the client knows to not send any traffic to the pod.
