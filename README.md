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

## Demo

```bash
kubectl create namespace open-match-demo
```

```bash
kubectl apply --namespace open-match-demo \
    -f https://open-match.dev/install/v1.8.0/yaml/02-open-match-demo.yaml
```

```bash
kubectl port-forward --namespace open-match-demo service/om-demo 51507:51507
```

### Uninstall Demo

```bash
kubectl delete namespace open-match-demo
```
