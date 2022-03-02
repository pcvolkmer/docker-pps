# docker-pps

The purpose of this tool is to show a list of processes running within docker containers.
It is some kind of `docker ps` for processes.

## Usage

```shell
$ docker-pps
```
shows
```
CONTAINER ID   IMAGE    PID        UID        COMMAND
0123456789ab   alpine   12345      1000       sh
123456789abc   nginx    13570      root       nginx: master process nginx -g daemon off;
123456789abc   nginx    13571      101        nginx: worker process
123456789abc   nginx    13572      101        nginx: worker process
```

To filter processes by user, provide `--uid` argument with comma separated usernames or IDs.

```shell
$ docker-pps --uid=root,101
```
shows
```
CONTAINER ID   IMAGE    PID        UID        COMMAND
123456789abc   nginx    13570      root       nginx: master process nginx -g daemon off;
123456789abc   nginx    13571      101        nginx: worker process
123456789abc   nginx    13572      101        nginx: worker process
```

Use `-q`/`--quiet` argument to show a list of PIDs

```shell
$ docker-pps -q
$ docker-pps --quiet
```
shows
```
12345
13570
13571
13572
```

To select another docker host, use `-H`/`--host` argument.

```shell
$ docker-pps -H http://otherhost:1234
$ docker-pps --host=ssh://otherhost:22
```