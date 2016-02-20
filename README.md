# simpleauth
Simple HTTP / LDAP bind authentication bridge

Usage:

### From one shell
```
$ go get github.com/peak6/simpleauth
$ LDAP_SERVER=somehost:port PORT=8888 $GOPATH/bin/simpleauth
```

### From another shell
```
$ curl -X POST -d '{ "username": "someone@somewhere.com", "password": "somepassword" }' http://localhost:8888/authenticate
``` 

A 200 OK response means you are authenticated, anything else is an error
