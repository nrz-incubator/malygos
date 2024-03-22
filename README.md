# Malygos - Kubernetes in Kubernetes provisioning on multiple proviers

Malygos is a tool to orchestrate Kubernetes provisionners installed on multiple Kubernetes management
clusters.

It permits to have decentralized Kubernetes cluster spawning, per region for example.

Note: This is currently in heavy development and experimental, use it at your own risk.

## Components

* `malygos-controller`: The main controller that will orchestrate regional controllers
* `kamaji`: The regional controller that will spawn Kubernetes clusters on a provider
* `cert-manager`: The certificate manager that will provide TLS certificates for the clusters