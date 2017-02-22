# Example for go distribution

## What does this test ?

This is a simple master/node example compiled in only one binary and launched in deveral docker containers.

Each node container will have is own ip address, and will register to "master:10000". Docker allows node to contact "master" with that name that is registered in `docker-compose.yml` file

By scaling up and down "node" containers, you'll see registration and deletion. 

Curl command on "/sum" endpoint will give you the "Arith.Sum" call response using random integers.

## How to

To make that test:

```bash
cd $GOPATH/github.com/metal3d/_example
make build
docker-compose up -d
docker-compose logs
```

Open a second terminal and try:

```bash
for i in $(seq 4); do curl -s localhost:10000/sum  & done; wait
```

You can scale up nodes, in the terminal where you launched docker-compose:

```bash
CTRL+C # stop logs
docker-compose scale node=4
docker-compose logs
```

And redo the command in second terminal:

```bash
for i in $(seq 4); do curl -s localhost:10000/sum  & done; wait
```

You can scale-down, and retry again...

To stop and cleanup docker:

```bash
CTRL+C
docker-compose down -v
```


