# Explore: GRPC Load Balancing

This repo holds my exploration on load balancing a grpc system. Currently, my team has one backend application that have high traffic (100k> rps from 9 pods). I'm not satisfied with current solution with the use of kubernetes' headless service. My team still see that there is an imbalance between the pods.

## Setup

I create two simple application to help me demonstrate that the current solution is not working properly. These two app is a client and server model. Each application has an UUID that I used as indentifier. The client has RESTful API endpoint, which upon called, will also make a gRPC call to the server application. This is the gRPC and RESTful contract.

```json
// RESTful contract
  "detail": {
    "from: 2a63ca4d-193f-4ec8-b9d3-49a39548bb95 replied by: 4bcbc60a-5700-4313-8772-c2a3ceab0513": 1000
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

Basically, the client app will count how many variance of the server app that serve its request. With this, we can conclude whether the request is well balanced or not.
