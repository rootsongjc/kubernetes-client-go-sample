# Kubernetes client-go sample
This is a simple code sample for using kubernetes [client-go](https://github.com/kubernetes/client-go) package.

You can use this code to update a specify deployment's application image (More than one container in a pod).

## Build

```bash
go build main.go
```

## Usage

```bash
Usage of ./main:
  -app string
    	application name (default "app")
  -deployment string
    	deployment name
  -image string
    	new image name
  -kubeconfig string
    	(optional) absolute path to the kubeconfig file (default "/Users/jimmy/.kube/config")
```

- `-image`: new image name
- `-deployment`: deployment name
- `-app`: (optional) application container name (default: app)
- `-kubeconfig`: (optional) absolute path to the kubeconfig file (default "$HOME/.kube/config")

## Example

```bash
./update-deployment-image -image test:Build_8 -deployment filebeat-test
Found deployment
name -> filebeat-test
Old image -> test:Build_7
New image -> test:Build_8
```



