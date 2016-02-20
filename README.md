# simpleauth
Simple HTTP / LDAP bind authentication bridge

Usage:

### From one shell
```
$ go get github.com/peak6/simpleauth
$ LDAP_SERVER=somehost:port PORT=8888 $GOPATH/bin/simpleauth
```
Default value for PORT is 8080.

There is no default for LDAP_SERVER, you must set it.

### From another shell
```
$ curl -X POST -d '{ "username": "someone@somewhere.com", "password": "somepassword" }' http://localhost:8888/authenticate
``` 

A 200 OK response means you are authenticated, anything else is an error

## Docker instructions

###build with: 
```
$ docker build -t simpleauth $GOPATH/src/github.com/peak6/simpleauth/.  # if you are already in this dir, use .
```

### run with:
```
$ docker run -it --rm --name sa -p 8080:8888 -e LDAP_SERVER=yourserver:port simpleauth
```
-p 8080:8888 tells docker to map simpleauth's default port 8080 to the docker host's port 8888.

Change 8888 to whatever you want to avoid port conflicts on the host.

