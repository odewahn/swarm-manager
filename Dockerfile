#
# This will make a super tiny docker image
# See https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/
#
# GOOS=linux go build -a -installsuffix cgo -o swarm-manager .
# docker build -t swarm-manager .
#
FROM scratch
ADD swarm-manager /swarm-manager
ADD .env.prod /.env
CMD ["/swarm-manager"]
