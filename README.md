# Malygos - Kubernetes in Kubernetes provisioning on multiple proviers

Malygos is a tool to orchestrate Kubernetes provisionners installed on multiple Kubernetes management
clusters.

It permits to have decentralized Kubernetes cluster spawning, per region for example.

Note: This is currently in heavy development and experimental, use it at your own risk.

## Components

* `malygos-controller`: The main controller that will orchestrate regional controllers
* `kamaji`: The regional controller that will spawn Kubernetes clusters on a provider
* `cert-manager`: The certificate manager that will provide TLS certificates for the clusters

## How to develop

### Prerequisites

* Go 1.21+

### Build (with direnv/nix)

If you have nix and direnv (really, you should !), all is ready for you to start hacking.

```bash
make malygos
```

### Build (without direnv/nix)

* Ensure GOPATH variable is set properly with a writable folder

```bash
make malygos
```
