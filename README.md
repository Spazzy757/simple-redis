# Simple Redis Operator

This is just a simplified example of a Kubernetes operator that deploys [redis](https://redis.io/) clusters. This is not to be used in a production environment.

This operator covers the following:

- [x] Deploys a master redis instance with networking setup
- [x] Deploys replicas redis instances that are setup to replicate the master instance
- [x] Allows some basic settings of the redis instances
- [x] Validation of input with sensible defaults

Potential roadmap items that could be added, but will not be for this iteration

- [ ] Setup state using a Storage Class
- [ ] Setup automated master election in case of failure of master redis instance
- [ ] Allow various scheduling options such as taints an tolerations 
- [ ] TLS setup between replicas and master
- [ ] Allow for setting various other configurations on the redis instances
- [ ] Multi Master setup

## Description

This is a simple redis operator, it will deploy a single master instance with
as many replicas as specified.

This is not a production grade application but was rather used to show the
different aspects of a Kubernetes operator including:

- Validating and Mutating Webhooks for sensible defaults as well as resource
  validation
- Reconcile logic
- Resource Ownership
- Idempotent Reonciliation

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/simple-redis:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/simple-redis:tag
```

### Webhooks

**Note:** The webhooks setup will require you to setup TLS. Yopu will need to
update the MutatingWebhookConfiguration and the ValidatingWebhookConfiguration
and supply the Certs to the operator. You can read up more on how this works in
the [Deploying Admission Webhooks][1] in the Kubebuilder Book

[1]: https://book.kubebuilder.io/cronjob-tutorial/running-webhook.html

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

