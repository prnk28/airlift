# Airlift

Airlift is a self-hosted file upload and sharing service. The clients upload
files to the server and return a nice link for you to share. Just bring your
own server and domain.

You should use [deuiore/load.link](https://github.com/deuiore/load.link)
instead of this if...

- ...you like PHP;
- ...you don't like me;
- ...you're on a free/cheap shared host that doesn't allow long-running
  processes.

#### Clients

- Web interface (included in server)
- [Cross-plaform (CLI)][cli]
- [OS X][osx]

[cli]: https://github.com/moshee/lift
[osx]: https://github.com/moshee/AirliftOSX

# airlift

`airlift` is the Airlift server. You drop the server on any dedicated, VPS,
shared host, whatever, as long as it supports running a binary and gives you
access to ports or frontend server reverse proxying. A client sends files to it
and recieves nice URLs to share. The server itself also provides a web-based
client to upload files from, as well as manage existing uploads and customize
some behaviors.

The server is packaged as a statically compiled binary with a few text assets
with no system dependencies apart from maybe libc for networking. Just download
(or clone and build), add to your init system of choice, and run.

You can choose to run it behind a frontend server or standalone. 

### Installing

#### If you just want a binary

`airlift` [linux/amd64](http://static.displaynone.us/airlift/airlift-linux_amd64.tar.bz2)

(proper cross-compiling is not really worth it because this imports the
`os/user` package)

#### Or if you want to build it yourself

1. [Install Go](http://golang.org/doc/install) and git
2. `$ mkdir ~/go && export GOPATH=~/go` (you can use any place as your GOPATH)
3. `$ go get -d -u ktkr.us/pkg/airlift`

I haven't tried to build or run it on Windows, YMMV. Works on OS X and
GNU+Linux.

`go get` should clone this repo along with any dependencies and place it
in your `$GOPATH`, which should be set (see [here][GOPATH] for more info).
This can be anywhere that isn't `$GOROOT`; you can set it to any arbitrary
place like `~/go`. Assuming that's what it is, then

```
~/go/src/ktkr.us/pkg/airlift$ go build
~/go/src/ktkr.us/pkg/airlift$ ./airlift
```

to build and run.

[GOPATH]: https://github.com/golang/go/wiki/GOPATH

The binary produced by `go build` will be in your working directory at the
moment you built it. By default, `go get` will install the binary to
`$GOPATH/bin` after building. It isn't very useful there, because...

### Usage

...the server must be run with the `templates` and `static` subdirectories in
its working directory. Whatever else you do is up to you. Just `$ ./airlift` to
run it in your terminal. Use your favorite tools to background it.

When you start the server for the first time, it will generate a dotfolder in
your home directory for local configuration. Visit
`http(s)://<yourhost>/config` to set up a password and change other
configuration parameters. On the first setup, an empty password will not be
accepted.

**Host** []: The base URL that links will be returned on. This includes domain
and path.

If you are proxying the server behind a frontend at a certain subdirectory,
make sure you rewrite the leading path out of the request URL so that the URLs
sent to `airlift` are rooted. Unfortunately, since URLs are rewritten, the
redirecting behavior of /login and /config won't work properly, so you'll have
to do your configuration on the internal port (60606 or whatever). Could use a
meta redirect instead of internal redirect to fix this, but that doesn't play
well with how sessions and stuff are set up in here.

Leaving the host field empty will cause the server to return whatever host the
file was posted to.

**Port** [60606]: This is the port the server executable listens on.

The `-p` flag overrides the configured port.

If you are using e.g. nginx, you can just add a
`proxy_pass http://localhost:60606;` directive inside a server block for the
host you choose.

**Directory** [~/.airlift-server/uploads]: This is where uploaded files will be
stored.

**Max upload age** [0]: If this value is greater than 0, uploads older than
that many days will be automatically deleted.

**Max upload size** [0]: If this value is greater than 0, the oldest uploads
will be pruned on every new upload until the total size is less than that many
megabytes.

If you manually edit the config file while the server is running, you should
send the server process a USR2 signal to force a config reload.

If the server fails to start with a config error, you probably want to delete
`~/.airlift-server/config` and reconfigure from scratch.

#### HTTPS

In order to use SSL/TLS standalone, set the following environment variables:

 Variable       | Value
----------------|---------------------------------------------
 `GAS_TLS_PORT` | The port for the secure server to listen on
 `GAS_TLS_CERT` | The path to your certificate
 `GAS_TLS_KEY`  | The path to your key
 `GAS_PORT`     | *Optional:* set this to -1 if you **only** want HTTPS, not regular HTTP.

If both HTTP and HTTPS are enabled, they will both serve from the same
executable and HTTP requests will redirect to HTTPS.
