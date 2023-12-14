# david - The simple WebDAV server... extended.

## Introduction

_david_ is a simple WebDAV server that provides the following features:

- Single binary that runs under Windows, Linux and OSX.
- Authentication via HTTP-Basic.
- CRUD operation permissions
- TLS support - if needed.
- A simple user management which allows user-directory-jails as well as full admin access to
  all subdirectories.
- Live config reload to allow editing of users without downtime.
- A cli tool to generate BCrypt password hashes.

It perfectly fits if you would like to give some people the possibility to upload, download or
share files with common tools like the OSX Finder, Windows Explorer or Nautilus under Linux
([or many other tools](https://en.wikipedia.org/wiki/Comparison_of_WebDAV_software#WebDAV_clients)).

The project david is an extension from the project ![david](https://github.com/audstanley/david)

## Table of Contents

- [Configuration](#configuration)
  * [First steps](#first-steps)
  * [TLS](#tls)
  * [Behind a proxy](#behind-a-proxy)
  * [User management](#user-management)
  * [Logging](#logging)
  * [Live reload](#live-reload)
- [Installation](#installation)
  * [Binary-Installation](#binary-installation)
  * [Build from sources](#build-from-sources)
  * [Build and run with Docker](#build-and-run-with-docker)
- [Connecting](#connecting)
- [Contributing](#contributing)
- [License](#license)

## Configuration

The configuration is done in form of a yaml file. _david_ will scan the
following locations for the presence of a `config.yaml` in the following
order:

- The directory `./config`
- The directory `$HOME/.swd` (swd was the initial project name of david)
- The directory `$HOME/.david`
- The current working directory `.`

Alternatively, the path to a configuration file can be specified on the
command-line:

```sh
david --config /path/to/config.yaml
```

### First steps

Here an example of a very simple but functional configuration:

```yaml
address: "127.0.0.1"    # the bind address
port: "8000"            # the listening port
dir: "/home/webdav"     # the provided base dir
prefix: "/webdav"       # the url-prefix of the original url
users:
  user:                 # with password 'foo' and jailed access to '/home/webdav/user'
    password: "$2a$10$yITzSSNJZAdDZs8iVBQzkuZCzZ49PyjTiPIrmBUKUpB0pwX7eySvW"
    subdir: "/user"
    permissions: "cru" # This user won't be able to delete from the server.
  admin:                # with password 'foo' and access to '/home/webdav'
    password: "$2a$10$DaWhagZaxWnWAOXY0a55.eaYccgtMOL3lGlqI3spqIBGyM0MD.EN6"
    permissions: "crud"
```

With this configuration you'll grant access for two users and the WebDAV
server is available under `http://127.0.0.1:8000/webdav`.

### TLS

At first, use your favorite toolchain to obtain a SSL certificate and
keyfile (if you don't  already have some).

Here an example with `openssl`:

```sh
# Generate a keypair
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365
# Remove the passphrase from the key file
openssl rsa -in key.pem -out clean_key.pem
```

Now you can reference your keypair in the configuration via:

```yaml
address: "127.0.0.1"    # the bind address
port: "8000"            # the listening port
dir: "/home/webdav"     # the provided base directory
tls:
  keyFile: clean_key.pem
  certFile: cert.pem
users:
  ...
...
```

The presence of the `tls` section is completely enough to let the server
start with a TLS secured HTTPS connection.

In the current release version you must take care, that the private key
doesn't need a passphrase. Otherwise starting the server will fail.

### Cross Origin Resource Sharing (CORS)

In case you intend to operate this server from a web browser based application,
you might need to allow [CORS](https://en.wikipedia.org/wiki/Cross-origin_resource_sharing)
access. To achieve that, you can configure the host you want to grant access to:

```yaml
cors:
  origin: "*"        # the origin to allow, or '*' for all
  credentials: true  # whether to allow credentials via CORS
```

Note however that this has security implications, so be careful in production
environments.

### Behind a proxy

_david_ will also work behind a reverse proxy. Here is an example
configuration with `apache2 httpd`'s `mod_proxy`:

```xml
<Location /webdav>
  ProxyPass           https://webdav-host:8000/
  ProxyPassReverse    https://webdav-host:8000/
</Location>
```

Here is an example of david using a [json caddyfile](https://caddyserver.com/docs/json/) for a reverse proxy:
```json
{
    "admin": {
      "disabled": false,
      "listen": "0.0.0.0:2019",
      "enforce_origin": false,
      "origins": [
        "127.0.0.1"
      ],
      "config": {
        "persist": false
      }
    },
    "apps": {
      "http": {
        "servers": {
          "MyServers": {
            "listen": [
              ":443"
            ],
            "routes": [
              {
                "match": [
                  {
                    "host": [
                      "files.example.com"
                    ]
                  }
                ],
                "handle": [
                  {
                    "handler": "reverse_proxy",
                    "upstreams": [
                      {
                        "dial": ":8000"
                      }
                    ]
                  }
                ]
              }
            ]
          }
        }
      }
    }
  }
```

### User management

User management in _david_ is very simple, but optional. You don't have to add users if it's not
necessary for your use case. But if you do, each user in the `config.yaml` **must** have a
password and **can** have a subdirectory.

The password must be in form of a BCrypt hash. You can generate one calling the shipped cli
tool `davidcli passwd`.

If a subdirectory is configured for a user, the user is jailed within it and can't see anything
that exists outside of this directory. If no subdirectory is configured for an user, the user
can see and modify all files within the base directory.

### Logging

You can enable / disable logging for the following operations:

- **C**reation of files or directories
- **R**eading of files or directories
- **U**pdating of files or directories
- **D**eletion of files or directories

You can also enable or disable the error log.

All file-operation logs are disabled per default until you will turn it on via the following
config entries:

```yaml
address: "127.0.0.1"    # the bind address
port: "8000"            # the listening port
dir: "/home/webdav"     # the provided base directory
log:
  error: true
  create: true
  read: true
  update: true
  delete: true
...
```

Be aware, that the log pattern of an attached tty differs from the log pattern of a detached tty.

Example of an attached tty:

	INFO[0000] Server is starting and listening              address=0.0.0.0 port=8000 security=none

Example of a detached tty:

	time="2018-04-14T20:46:00+02:00" level=info msg="Server is starting and listening" address=0.0.0.0 port=8000 security=none

### Live reload

There is no need to restart the server itself, if you're editing the user or log section of
the configuration. The config file will be re-read and the application will update it's own
configuration silently in background.


## Installation

### Build from sources

#### Setup

3. Clone the repository (or your fork)

```sh
git clone https://github.com/micromata/david.git
```

```sh
cd cmd/david && go build . && mv ./david ~/go/bin/david
```

#### magefile

You can also use mage to build the project.

Please ensure you've got [mage](https://magefile.org) installed. This can be done with the
following steps:

```sh
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
```

The [golint](https://github.com/golang/lint) utility is also required for `mage check`:

```sh
go install golang.org/x/lint/golint
```

Now you can call `mage install` to build and install the binaries. If you just call `mage`,
you'll get a list of possible targets:

	Targets:
	  build            Builds david and davidcli and moves it to the dist directory
	  buildReleases    Builds david and davidcli for different OS and package them to a zip file for each os
	  check            Runs golint and go tool vet on each .go file.
	  clean            Removes the dist directory
	  fmt              Formats the code via gofmt
	  install          Installs david and davidcli to your $GOPATH/bin folder
	  installDeps      Runs dep ensure and installs additional dependencies.

### Build and run with Docker

The image of david is available on Docker Hub as [`micromata/david`](https://hub.docker.com/r/micromata/david).

If you like to build it for your own, just execute the following lines:

```sh
git clone https://github.com/micromata/david.git
cd david
docker build -t micromata/david:latest .
```

You can run david with the following lines:

```sh
# create webdav home
mkdir webdav-home

docker run -d \
	-p 8000:8000 \
	-v $(pwd)/examples/config-sample.yaml:/config.yaml:ro \
	-v $(pwd)/webdav-home:/tmp:rw \
	micromata/david:latest
```

## Connecting

You could simply connect to the WebDAV server with an HTTP(S) connection and a tool that
allows the WebDAV protocol.

For example: Under OSX you can use the default file management tool *Finder*. Press _CMD+K_,
enter the server address (e.g. `http://localhost:8000`) and choose connect.

## Contributing

Everyone is welcome to create pull requests for this project. If you're new to github, take
a look [here](https://help.github.com/categories/collaborating-with-issues-and-pull-requests/)
to get an idea of it.

If you'd like to contribute, please make sure to use the [magefile](#magefile) and execute and
check the following commands before starting a PR:

```sh
mage fmt
mage check
```

If you've got an idea of a function that should find it's way into this project, but you
won't implement it by yourself, please create a new issue.

## License

Please be aware of the licenses of the components we use in this project. Everything else that has
been developed by the contributions to this project is under the [Apache 2 License](LICENSE.txt).
