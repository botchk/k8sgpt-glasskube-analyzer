# k8sgpt-glasskube-analyzer
Custom k8sgpt analyzer for Glasskube

# How to run locally

```
go run main.go
```

# Test with grpcurl

```
grpcurl --plaintext localhost:8085 schema.v1.CustomAnalyzerService/Run
```

# TODO

## k8sgpt-glasskube-analyzer
TODO PR checks should be fulfilled before merging is allowed
TODO create github test user to test setup
TODO change deployment name to not include chart
TODO readiness and liveness probe
TODO describe which version of k8sgpt operator and k8sgpt are needed
  include k8sgpt config
TODO build gets triggered two times on PRs
TODO when does commit lint get triggered?
TODO run pod with increased security settings