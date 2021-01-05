# `kubectl grid`

A kubectl plugin to create test clusters and manage test runs on the clusters.

## Usage

### Create a grid

```shell
$ kubectl grid create --from-yaml ./examples/basic/grid.yaml
```

### Deploy an app to all clusters in the grid

```shell
$ kubectl grid deploy --grid eks-existing --application ./examples/basic/kots-app.yaml
```

### List namespaces on one of the clusters in the grid

```shell
$ kubectl grid get grids

$ kubectl grid describe grid eks-existing

$ kuebctl grid get ns --grid eks-existing --cluster 

```

### Execute an experiment on all applications in the grid

```shell

```

### Delete and clean up all resources created

```shell

```

## Questions

**Why not Terraform/Pulumi?**  

Great question. We love all of the mature automation tools that can spin up production-grade, best-practices, rock-solid infrastructure. But a test grid isn't production grade. We wanted to focused on test grid execution speed instea of building production-grade clusters. So, keep using Terraform and other tools to manage your prod cluster, this project is NOT what you want there.



