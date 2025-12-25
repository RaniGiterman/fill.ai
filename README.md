# Fill.ai

Fill JSON based on HTML code, using LLM functionality. Currently using Ollama llama2.3 model, with modelfile to answer conscisely


### Sealed Secrets
App using sealed secret to conceal the openai api key.

```sh
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml

echo -n bar | kubectl create secret generic mysecret --dry-run=client --from-file=foo=/dev/stdin -o json >mysecret.json

kubeseal -f mysecret.json -w mysealedsecret.json

kubectl create -f mysealedsecret.json

kubectl get secret mysecret

kubectl port-forward service/fill-ai-service 8080:8080
```

