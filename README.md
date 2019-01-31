# SServe
_Static content on https connection trusted by localhost_

### The idea
sserve is a simple tool for serving static content on https thanks to locally-trusted development certificates.
It requires no configuration.

### About
Serving static content on localhost in a trusted SSL connection is not so simple.  
It requires to manually generate and trust certificates, with complicate commands and many manual steps.

sserve, serves static content using a locally-trusted certificate, generated with the well-knowed [mkcert](https://github.com/FiloSottile/mkcert) tool.

When you install sserve it automatically creates and installs a local CA in the system (and browsers) root store, and generates the certificate for you.  
No configuration is required, just lunch the tool and we take care of everything you need.

## Installation
**Warning:** the `rootCA-key.pem` file that mkcert automatically generates when installing sserve gives complete power to intercept secure requests from your machine. Do not share it.

**Requires Go:** follow [this instructions](https://golang.org/doc/install) to install it.

### MacOS
```
go get github.com/daquinoaldo/sserve
```

### Linux
First install `certutil`
```
sudo apt install libnss3-tools
    -or-
sudo yum install nss-tools
    -or-
sudo pacman -S nss
```

Then install sserve
```
go get github.com/daquinoaldo/sserve
```

### Windows
```
go get github.com/daquinoaldo/sserve
```

## Troubleshooting
If you're running into permission problems try running `sserve` as an Administrator.

For example, on some operating system is not possible to use the 443 port without sudo permissions.

## Supported root stores
_The supported root stores are the one supported by mkcert.  
Checkout the updated list [here](https://github.com/FiloSottile/mkcert/blob/master/README.md#supported-root-stores)._

**Here there is a handy copy:**
- macOS system store
- Windows system store
- Linux variants that provide either
    - `update-ca-trust` (Fedora, RHEL, CentOS) or
    - `update-ca-certificates` (Ubuntu, Debian) or
    - `trust` (Arch)
- Firefox (macOS and Linux only)
- Chrome and Chromium
- Java (when `JAVA_HOME` is set)
