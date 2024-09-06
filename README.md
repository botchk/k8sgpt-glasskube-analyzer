# k8sgpt-glasskube-analyzer
Custom [k8sgpt analyzer](https://github.com/k8sgpt-ai/k8sgpt) for [Glasskube](https://github.com/glasskube/glasskube)

# How to run locally

```
go run main.go
```

# Test with grpcurl

```
grpcurl --plaintext localhost:8085 schema.v1.CustomAnalyzerService/Run
```

# Changelog

[Changelog](CHANGELOG.md)

# TODOs

* PR checks should be fulfilled before merging is allowed
* create github test user to test setup
* change deployment name to not include chart
* readiness and liveness probe
* describe which version of k8sgpt operator and k8sgpt are needed
* include k8sgpt config
* run pod with increased security settings
* publish helm chart