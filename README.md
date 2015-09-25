# Swarm manager
Manage kernels on a Rackspace Swarm


## Start interlock on the cluster

Use the following command to start Interlock on the swarm:

```
DOCKER run -d \
   --name interlock \
   -p 80:80 \
   --volumes-from swarm-data \
   ehazlett/interlock \
   --swarm-url $DOCKER_HOST \
   --swarm-tls-ca-cert=/etc/docker/ca.pem  \
   --swarm-tls-cert=/etc/docker/server-cert.pem \
   --swarm-tls-key=/etc/docker/server-key.pem \
   --plugin haproxy start
```

You will also need to set up your DNS records to the IP address of the container; be sure to also enable DNS wildcards.  So, if you set up "interlock.example.com," you also need to add "*.interlock.example.com".
