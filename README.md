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
TODO BUG the release should be managed by the github actions bot
TODO change deployment name to not include chart
TODO readiness and liveness probe
TODO describe which version of k8sgpt operator and k8sgpt are needed
  include k8sgpt config
TODO helm chart app version
TODO golang app version
TODO renovate
TODO build gets triggered two times on PRs
TODO when does commit lint get triggered?

## k8sgpt-operator
TODO BUG report k8sgpt-operator should publish latest image or delete existing latest image