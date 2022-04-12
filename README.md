# Node Label Controller

## Overview

Node label Controller is a controller which watches the Kubernetes nodes and attaches a label `k8c.io/uses-container-linux: 'true'` to the Kubernetes Node object when the node uses `Flatcar Container Linux` as an operating system.

## Running controller locally

Prerequisite: Your terminal must in context of a kubernetes cluster.

```terminal
go run main.go
```

## Running controller on a Kubernetes Cluster

### Build docker image

```terminal
docker build ./ -t node-label-controller:v1.0.0
```

After building the image please upload the image to your image repository

### Deploy to cluster

```terminal
```
