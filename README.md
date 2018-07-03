# klstr - A μ-services platform

`klstr` is a platform to build, deploy, manage and monitor your services. klstr is built on top of
Kubernetes and abstracts over the primitives such as pods, deployments,
services, statefulsets etc., to provide a meaningful interface to your
services. 
To install `klstr` you can use homebrew if you are using a Mac.

    $ brew tap klstr/klstr
    $ brew install klstr

You also need to install [minikube](https://kubernetes.io/docs/setup/minikube/). Once you have your kubernetes cluster running you can adopt the cluster. Adopting a cluster enables `klstr` to install the necessary controllers on the cluster to create kubernetes resources such as secrets, deployments, services, ingress resources and so on.

    $ kubectl get nodes # ensure that the cluster is ready before running this.
    $ klstr adopt --klstr-name=dev --default

You can inspect whats running and the version of klstr by using the following command.

    $ klstr status --klstr-name=dev
    klstr client version v0.1.0
    klstr server version v0.1.0

    components
    - ✅ logging (oklog)
    - ✅ monitoring (prometheus-controller v0.8.9)
    - ✅ grafana v0.3.2
    - ✅ jaeger v0.3.4

Go ahead and clone the following project. http://github.com/klstr/muservice

## Muservice readme.

MuService is a sample golang app that runs on port 4323 and expects a PostgreSQL
database and a redis instance. We can describe muservice using the following yaml.

    version: io.klstr/v1
    # muservice is a CRD that defines a micro service
    kind: Muservice
    metadata:
      labels:
        app: muservice
      name: muservice
    spec:
      image: quay.io/klstr/muservice:v0.1.0
      replicas: 3
      ports:
      - port: 4323
      # the expose directive creates an ingress resource matching
      # to expose the service to the outside world using the domain name
      expose: muservice.minikube.local
      environment:
        - name: RABBITMQ_URL
          valueFrom:
            configMap:
              name: muservice-confmap
              key: RABBITMQ_URL
      # These services are managed by klstr and will be deleted
      # if the parent service is deleted.
      services:
        - name: samplecache
          type: redis
          # this creates a redis stateful set deployment
          # and exposes environment variables for SAMPLECACHE_REDIS_HOST
          # SAMPLECACHE_REDIS_PORT for the app to connect to
      databases:
        - name: mysampledb
          type: postgres
          instance: dev
          # this creates a user and a database and sets
          # a secret environment variable for LOCATION_DATABASE_URI,
          # LOCATION_PG_HOST, LOCATION_PG_PORT, LOCATION_PG_USER and
          # LOCATION_PG_PASSWORD


To deploy the service, run the following command.

    $ klstr deploy -f muservice-with-ingress.yaml

To access the service ensure that you add the relevant `/etc/hosts` entry to point `muservice.minikube.local` to minikube's IP address which is usually on `192.168.99.100`. 

To deploy a new version of the service.

    klstr deploy muservice --image quay.io/repo/mysample:0.1.2


To scale out the deployment, you can do the following.

    klstr deploy muservice --scale=3
