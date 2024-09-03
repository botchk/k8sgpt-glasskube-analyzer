package analyzer

import (
	"context"
	"fmt"
	"log/slog"

	rpc "buf.build/gen/go/k8sgpt-ai/k8sgpt/grpc/go/schema/v1/schemav1grpc"
	v1 "buf.build/gen/go/k8sgpt-ai/k8sgpt/protocolbuffers/go/schema/v1"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	pkgv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Handler struct {
	rpc.CustomAnalyzerServiceServer
	Client client.PackageV1Alpha1Client
}
type Analyzer struct {
	Handler *Handler
}

func (a *Handler) Run(context.Context, *v1.RunRequest) (*v1.RunResponse, error) {
	results := []v1.Result{}
	if err := a.analyzePackageRepositories(&results); err != nil {
		return nil, err
	}
	if err := a.analyzeClusterPackages(&results); err != nil {
		return nil, err
	}
	if err := a.analyzePackages(&results); err != nil {
		return nil, err
	}

	// TODO why can there be only one result each analysis? Is it needed to register one analyzer for each resource?
	response := &v1.RunResponse{}
	return response, nil
}

func (a *Handler) analyzePackageRepositories(results *[]v1.Result) error {
	var repoList v1alpha1.PackageRepositoryList
	if err := a.Client.PackageRepositories().GetAll(context.TODO(), &repoList); err != nil {
		return err
	}

	for _, repo := range repoList.Items {
		for _, condition := range repo.Status.Conditions {
			slog.Debug("packagerepository", "name", repo.Name, "condition.type", condition.Type, "condition.reason",
				condition.Reason, "condition.message", condition.Message)
			if condition.Status == pkgv1.ConditionFalse {
				// TODO the result should probably include full crd types to enable better AI analysis
				*results = append(*results, v1.Result{
					Name: "k8sgpt-glasskube-analyzer",
					Error: []*v1.ErrorDetail{
						{
							Text: fmt.Sprintf("%s has condition of type %s, reason %s: %s",
								repo.Name, condition.Type, condition.Reason, condition.Message),
						},
					},
					Kind: repo.Kind,
				})
			}
		}
	}
	return nil
}

func (a *Handler) analyzeClusterPackages(results *[]v1.Result) error {
	var pkgList v1alpha1.ClusterPackageList
	if err := a.Client.ClusterPackages().GetAll(context.TODO(), &pkgList); err != nil {
		return err
	}

	for _, pkg := range pkgList.Items {
		for _, condition := range pkg.Status.Conditions {
			slog.Debug("clusterpackage", "name", pkg.Name, "condition.type", condition.Type, "condition.reason",
				condition.Reason, "condition.message", condition.Message)
			//TODO implement
			*results = append(*results, v1.Result{})
		}
	}

	return nil
}

func (a *Handler) analyzePackages(results *[]v1.Result) error {
	var pkgList v1alpha1.PackageList
	//TODO check if supplying an empty string as namespace returns all packages over all namespaces
	if err := a.Client.Packages("").GetAll(context.TODO(), &pkgList); err != nil {
		return err
	}

	for _, pkg := range pkgList.Items {
		for _, condition := range pkg.Status.Conditions {
			slog.Debug("package", "name", pkg.Name, "condition.type", condition.Type, "condition.reason",
				condition.Reason, "condition.message", condition.Message)
			//TODO implement
			*results = append(*results, v1.Result{})
		}
	}

	return nil
}
