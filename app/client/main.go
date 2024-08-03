package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	echo_rpc "github.com/ekky-kharismadhany/grpc-load-balancer/client-demo/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port = "port"
const serverHost = "serverHost"

type config struct {
	port       string
	serverHost string
}

func loadConfig() config {
	config := config{}
	config.port = os.Getenv(port)
	config.serverHost = os.Getenv(serverHost)
	return config
}

var app_id = uuid.NewString()

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

type resultCall struct {
	Detail map[string]int `json:"detail"`
}

func countEachCall(messages []string) resultCall {
	resultCall := resultCall{
		Detail: map[string]int{},
	}
	for _, message := range messages {
		_, ok := resultCall.Detail[message]
		if !ok {
			resultCall.Detail[message] = 1
			continue
		}
		resultCall.Detail[message]++
	}
	return resultCall
}

func handleEchoCall(client echo_rpc.EchoServerClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		iter := r.URL.Query().Get("iter")
		iterInt, err := strconv.Atoi(iter)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Info("receive echo call", "iter", iterInt)
		var messages []string
		for range iterInt {
			response, err := client.CallEcho(context.Background(), &echo_rpc.Echo{
				Message: app_id,
			})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			messages = append(messages, response.Message)
		}
		jsonByte, err := json.Marshal(countEachCall(messages))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(jsonByte)
	}
}

func main() {
	config := loadConfig()
	host := fmt.Sprintf(":%s", config.port)
	insecureCred := grpc.WithTransportCredentials(insecure.NewCredentials())
	logger.Info("dialling server", "server host", config.serverHost)
	conn, err := grpc.NewClient(config.serverHost, insecureCred)
	if err != nil {
		logger.Error(err.Error())
		panic(err.Error())
	}
	defer conn.Close()
	client := echo_rpc.NewEchoServerClient(conn)
	http.HandleFunc("/echo", handleEchoCall(client))
	logger.Info("start client app", "host", host)
	if err := http.ListenAndServe(host, nil); err != nil {
		logger.Error(err.Error())
		panic(err.Error())
	}
}
