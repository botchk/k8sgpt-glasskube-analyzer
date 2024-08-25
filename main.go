package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	rpc "buf.build/gen/go/k8sgpt-ai/k8sgpt/grpc/go/schema/v1/schemav1grpc"
	"github.com/botchk/k8sgpt-glasskube-analyzer/pkg/analyzer"
	"github.com/glasskube/glasskube/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// TODO allow setting of log level from env
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))

	// TODO port should be configureable
	var port string = "8085"
	slog.Info("startup", "port", port)
	address := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	// TODO allow setting localMode from env
	// TODO allow setting kubeconfig path from env
	var localMode bool = true
	var config *rest.Config
	if localMode {
		config, err = clientcmd.BuildConfigFromFlags("", "/home/daugustin/.kube/config")
		if err != nil {
			panic(err)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}

	var pkgClient client.PackageV1Alpha1Client
	if pkgClient, err = client.New(config); err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	aa := analyzer.Analyzer{Handler: &analyzer.Handler{Client: pkgClient}}
	rpc.RegisterCustomAnalyzerServiceServer(grpcServer, aa.Handler)

	if err := grpcServer.Serve(
		lis,
	); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return
	}
}
