# Node Label Controller

## Overview

Node Label Controller is a controller which watches the Kubernetes nodes and attaches a label `k8c.io/uses-container-linux: 'true'` to the Kubernetes Node object when the node uses `Flatcar Container Linux` as an operating system.

## Running controller locally

Prerequisite: Your terminal must be in context of a kubernetes cluster.

```terminal
go run main.go
```

## Running controller on a Kubernetes Cluster

### Build docker image

```terminal
docker build ./ -t node-label-controller:v1.0.0
```

After building the image please upload the image to your image repository and update the [deployment](https://github.com/abdulhaseeb2/node-label-controller/blob/master/manifests/deployment.yaml#L25) or you can use my image from dockerhub [abdulhaseeb2/node-label-controller:v0.0.1](https://hub.docker.com/repository/docker/abdulhaseeb2/node-label-controller/general)

### Deploy to cluster

First create a namespace with the name `node-label-controller`.

```terminal
$ kubectl create namespace node-label-controller
namespace/node-label-controller created
```

Now, create all resources under manifest folder.

```terminal
$ kubectl apply -f manifests/
deployment.apps/node-label-controller created
configmap/node-label-config created
service/node-label-controller-metrics-service created
serviceaccount/node-label-controller created
role.rbac.authorization.k8s.io/node-label-controller-leader-election-role created
clusterrole.rbac.authorization.k8s.io/node-label-role created
clusterrole.rbac.authorization.k8s.io/node-label-controller-metrics-reader created
clusterrole.rbac.authorization.k8s.io/node-label-controller-proxy-role created
rolebinding.rbac.authorization.k8s.io/node-label-controller-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/node-label-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/node-label-controller-proxy-rolebinding created
```

Once the controller has been deployed check the status of the controller by viewing the pod in `node-label-controller` namespace.

```terminal
$ kubectl get pods -n node-label-controller
NAME                                     READY   STATUS              RESTARTS   AGE
node-label-controller-54bdf54754-xb9hq   2/2     Running             0          22s
```

Finally we can see labels being populated on Nodes containing `FlatCar Container Linux` OS.

```terminal
$ kubectl get nodes -l k8c.io/uses-container-linux=true
NAME                 STATUS   ROLES                  AGE   VERSION
kind-control-plane   Ready    control-plane,master   94m   v1.20.2
```

## Future prospects

- Add support for custom labels
- Add support for other OS types
