package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"

	rpc "buf.build/gen/go/k8sgpt-ai/k8sgpt/grpc/go/schema/v1/schemav1grpc"
	"github.com/botchk/k8sgpt-glasskube-analyzer/pkg/analyzer"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	Local    bool   `mapstructure:"local"`
	Port     string `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
}

func main() {
	conf, err := initConfig()
	if err != nil {
		panic(err)
	}

	logLevel, err := parseLogLevel(conf.LogLevel)
	if err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(handler))
	slog.Info("init", "config", conf)

	address := fmt.Sprintf(":%s", conf.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	var config *rest.Config
	if conf.Local {
		kubeconfigFilePath, err := getKubeconfigFilePath()
		if err != nil {
			panic(err)
		}
		slog.Info("init", "kubeconfig", kubeconfigFilePath)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigFilePath)
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

func initConfig() (*Config, error) {
	viper.SetEnvPrefix("K8SGPT")
	viper.MustBindEnv("local")
	viper.SetDefault("local", false)

	viper.MustBindEnv("port")
	viper.SetDefault("port", "8085")

	viper.MustBindEnv("log_level")
	viper.SetDefault("log_level", "INFO")

	conf := Config{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func getKubeconfigFilePath() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".kube", "config"), nil
}

func parseLogLevel(s string) (slog.Level, error) {
	bytes := []byte(s)

	var level slog.Level
	err := level.UnmarshalText(bytes)

	return level, err
}
