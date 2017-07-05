Cowsay ~~ Meow ~~
=================

This is a demo with golang and RabbitMQ.

Prerequisities
--------------

- [golang](https://golang.org)
- [RabbitMQ](https://www.rabbitmq.com)
- [Docker Compose](https://docs.docker.com/compose/)

Start up
--------

On the root folder, simply type,

```bash
docker-compose up --build
```

As it has [autoreload](https://github.com/handwritingio/autoreload) package installed, you might want to start the watcher services also.

```bash
# watch the cowsay service
cd cowsay
./watcher.sh

#watch the main service
cd server
./watcher.sh
```

Start up production
-------------------

You can start up production by,

```bash
docker-compose -f docker-compose.prod.yml up
```

It will bind the http service on 80 port instead of 8080 on the host machine.
