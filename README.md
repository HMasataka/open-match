# Open Match

## Install

- [Reference](https://open-match.dev/site/docs/installation/yaml/)
- [demo](https://open-match.dev/site/docs/getting-started/)

```bash
kubectl apply --namespace open-match \
    -f https://open-match.dev/install/v1.8.0/yaml/01-open-match-core.yaml
```

```bash
kubectl apply --namespace open-match \
  -f https://open-match.dev/install/v1.8.0/yaml/06-open-match-override-configmap.yaml \
  -f https://open-match.dev/install/v1.8.0/yaml/07-open-match-default-evaluator.yaml
```

## Deploy

```bash
kubectl apply -f manifests/matchmaker.yaml
```

### Uninstall

```bash
kubectl delete namespace open-match
```

## Dev

```bash
skaffold dev --port-forward
```
