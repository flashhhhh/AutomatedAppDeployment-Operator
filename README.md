# AutomatedAppDeployment-Operator
This is a Kubernetes Operator to automatically deploy an app with a specific CRD specification for this app.

## Problems
Suppose that you have an application that you want to deploy in a Kubernetes cluster. You want to deploy these application pods in a deployment in a specific namespace, and you want to expose these pods with a service.

Instead of manually set up the deployment and service, this project provides an operator that will automatically deploy the application pods in a deployment and expose them with a service.

## How to use
1. Create your own cluster.

### Minikube

### Kind
#### Install Kind
```bash
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.27.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/
```

#### Create a cluster
Suppose that you want to have a cluster with 1 control plane node and 2 worker nodes. You can create a cluster with the following command:

```yaml
# kind-config.yaml
apiVersion: kind.x-k8s.io/v1alpha4
kind: Cluster
name: kind-cluster
nodes:
  - role: control-plane
  - role: worker
  - role: worker
```

Then, you can create the cluster with the following command:
```bash
kind create cluster --config kind-config.yaml
```

2. Create the general CRD (**AutomatedAppDeployment**).
```bash
kubectl apply -f https://github.com/flashhhhh/AutomatedAppDeployment-Operator/blob/main/k8s-specifications/v1.0.0/automatedappdeployment-crd.yaml
```

3. Create the operator controller.
```bash
kubectl apply -f https://github.com/flashhhhh/AutomatedAppDeployment-Operator/blob/main/k8s-specifications/v1.0.0/automatedappdeployment-controller.yaml"
```

4. Create the operator's permission.
```bash
kubectl apply -f https://github.com/flashhhhh/AutomatedAppDeployment-Operator/blob/main/k8s-specifications/v1.0.0/automatedappdeployment-operator-permission.yaml"
```

5. Create your own CRD specification for the app you want to deploy, based on the general CRD **AutomatedAppDeployment**. This is an example of a CRD specification:
```yaml
# automated-postgres.yaml
apiVersion: automation.local.io/v1
kind: AutomatedAppDeployment
metadata:
  name: postgres
  namespace: default
  labels:
    app: postgres
spec:
  replicas: 2
  deployments:
    - image: postgres
      ports:
        - 5432
      envVars:
        POSTGRES_PASSWORD: "12345678"
```

Then, you can create the CRD with the following command:
```bash
kubectl apply -f automated-postgres.yaml
```

6. Check the status of the deployment and service.
```bash
kubectl get deployments -n default
kubectl get services -n default
kubectl get pods -n default
```

7. Check the logs of the operator controller.
```bash
kubectl logs -l app=automatedappdeployment-operator -n automatedappdeployment-operator-system
```