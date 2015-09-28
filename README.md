# Jupyter Kernel Manager for Rackspace Swarm

This is a simple API for starting and stopping containers running Jupyter on the Rackspace swarm.

## Start interlock on the cluster

Use the following command to start Interlock on the swarm:

```
docker run -d \
   --name interlock \
   -p 80:80 \
   -P \
   --hostname i2.odewahn.com \
   --volumes-from swarm-data \
   ehazlett/interlock \
   --swarm-url $DOCKER_HOST \
   --swarm-tls-ca-cert=/etc/docker/ca.pem  \
   --swarm-tls-cert=/etc/docker/server-cert.pem \
   --swarm-tls-key=/etc/docker/server-key.pem \
   --plugin haproxy start
```

You will also need to set up your DNS records to the IP address of the container; be sure to also enable DNS wildcards.  So, if you set up "interlock.example.com," you also need to add "*.interlock.example.com".

## Building swarm-manager

```
GOOS=linux go build -a -installsuffix cgo -o swarm-manager .
docker build --no-cache -t swarm-manager .
```

## Running the swarm-manager

```
docker run -d  \
   -p 3000 \
   -P \
   --volumes-from swarm-data \
   --hostname swarm-manager.i2.odewahn.com \
   --name swarm_manager \
   swarm-manager
```

## Managing Notebooks on the Swarm

### Spawn container

Submit a `POST` request to `/SPAWN` with the following values:
* `image`.  The name of the Docker image to launch
* `user`.  The user to assign it to.

The `spawn` call will immediately return a hostname, although it may take a while for the container to actually start.  Poll the `status` firled of the `/container/{hostname}` route to check the availability of the container.  The field will change to `ACTIVE` then the container has started.

```
$ curl --data "image=zischwartz/notebook&user=odewahn" http://swarm-manager.i2.odewahn.com/spawn
{
  "hostname": "ki60dje40khk",
  "domainname": "i2.odewahn.com",
  "image": "zischwartz/notebook",
  "url": "ki60dje40khk.i2.odewahn.com",
  "ContainerId": "",
  "Status": "",
  "StartTime": "2015-09-28T16:55:04.715248432Z",
  "User": "odewahn"
}
```

```
$ curl --data "image=ipython/scipystack&user=rmadsen" http://swarm-manager.i2.odewahn.com/spawn
{
  "hostname": "rstnjd3zyu25",
  "domainname": "i2.odewahn.com",
  "image": "ipython/scipystack",
  "url": "rstnjd3zyu25.i2.odewahn.com",
  "container_id": "",
  "status": "",
  "start_time": "2015-09-28T17:00:10.557520925Z",
  "user": "rmadsen"
}
```

### Check container status

Once a container spawns, use the `/container/{hostname}` route to monitor it's progress.  The `Status`
```
$ curl http://swarm-manager.i2.odewahn.com/container/ki60dje40khk
{
  "hostname": "ki60dje40khk",
  "domainname": "i2.odewahn.com",
  "image": "zischwartz/notebook",
  "url": "ki60dje40khk.i2.odewahn.com",
  "container_id": "ddeeb16fd22abec6047f57b317f2aa605e2aec7b50ab60854c10097d6101cb7a",
  "status": "ACTIVE",
  "start_time": "2015-09-28T16:55:04.715248432Z",
  "user": "odewahn"
}
```

### Kill a Container

```
$ curl http://swarm-manager.i2.odewahn.com/container/ki60dje40khk/kill
{
  "hostname": "ki60dje40khk",
  "domainname": "i2.odewahn.com",
  "image": "zischwartz/notebook",
  "url": "ki60dje40khk.i2.odewahn.com",
  "container_id": "ddeeb16fd22abec6047f57b317f2aa605e2aec7b50ab60854c10097d6101cb7a",
  "status": "DELETING",
  "start_time": "2015-09-28T16:55:04.715248432Z",
  "user": "odewahn"
}
```

### See all containers

```
$ curl http://swarm-manager.i2.odewahn.com/containers
[
  {
    "hostname": "80vgjpwxeg3u",
    "domainname": "i2.odewahn.com",
    "image": "zischwartz/notebook",
    "url": "80vgjpwxeg3u.i2.odewahn.com",
    "container_id": "",
    "status": "REMOVED",
    "start_time": "0001-01-01T00:00:00Z",
    "user": "odewahn"
  },
  {
    "hostname": "u1w8pgp2p9gz",
    "domainname": "i2.odewahn.com",
    "image": "zischwartz/notebook",
    "url": "u1w8pgp2p9gz.i2.odewahn.com",
    "container_id": "",
    "status": "REMOVED",
    "start_time": "0001-01-01T00:00:00Z",
    "user": "odewahn"
  },
  {
    "hostname": "vfgkh0j3r2my",
    "domainname": "i2.odewahn.com",
    "image": "ipython/scipystack",
    "url": "vfgkh0j3r2my.i2.odewahn.com",
    "container_id": "b054930b8aed4941ee87b7fdb83f87267c3888f51e93b743a58565cc658cab4f",
    "status": "ACTIVE",
    "start_time": "2015-09-28T13:04:51.502399801-04:00",
    "user": "rmadsen"
  },

  {
    "hostname": "vha5h0s1agd4",
    "domainname": "i2.odewahn.com",
    "image": "zischwartz/notebook",
    "url": "vha5h0s1agd4.i2.odewahn.com",
    "container_id": "",
    "status": "REMOVED",
    "start_time": "0001-01-01T00:00:00Z",
    "user": "odewahn"
  }
]
```

### Web interface to manage containers

There is a super-duper, ugly web interface to allow you to view and kill containers:

```
http://swarm-manager.i2.odewahn.com/manage
```
