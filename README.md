# Target environment

- ubuntu 18.04.5 LTS  
- minikube 1.6.2  
- Helm <= 2.16.1  
- go 1.15.6  
- make 4.1  
 
## Minikube

Add the following to enable the [nginx ingress controller](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/):

```shell
minikube addons enable ingress
```

Please host `eetest.testdomain.com` at your minikube address, for example (by amending `/etc/hosts`):

```shell
192.168.39.187 eetest.testdomain.com
```

## Utilities

Running `eval $(minikube docker-env)` before building docker builds the image into the minikube environemnt for immediate use. This must precede the docker build step and lasts for a terminal session.

See `Makefile` for utility commands. To deploy into a configured minikube environemnt:

`make helm-deploy`

To build the docker image:

```shell
make docker
```

## Immediate dependencies

Libraries use (docs): 
- [Gorilla Mux](https://github.com/gorilla/mux)
- [Go Flags](https://github.com/jessevdk/go-flags)
- [Testify](https://github.com/stretchr/testify)
