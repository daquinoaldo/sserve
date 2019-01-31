# SServe
_Static content on https connection trusted by localhost_

```
$ sserve ~/myproj
Serving ~/myproj on port 443. Checkout at https://localhost.
```
Easy. Handy. Free.

### The idea
sserve is a lightweight tool for serving static content on https thanks to locally-trusted development certificates.  
It requires you no configuration.

#### Why and how it works
Serving static content on localhost in a trusted SSL connection is not so simple.  
It requires to manually generate and trust certificates, with complicate commands and many manual steps.

sserve, serves static content using a locally-trusted certificate, generated with the well-knowed [mkcert](https://github.com/FiloSottile/mkcert) tool.

When you install sserve it automatically creates and installs a local CA in the system (and browsers) root store, and generates the certificate for you.  
No configuration is required, just lunch the tool and we take care of everything you need.


## Installation
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


## Things to know

### Warning
The `rootCA-key.pem` file that mkcert automatically generates when installing sserve gives complete power to intercept secure requests from your machine. Do not share it.

### License
Is released under [AGPL-3.0 - GNU Affero General Public License v3.0](LICENSE).

#### Briefly:
- modification and redistribution allowed for both private and **commercial use**
- you must **grant patent rigth to the owner and to all the contributors**
- you must **keep it open source** and distribute under the **same license**
- changes must be documented
- include a limitation of liability and it **does not provide any warranty**

### Warranty
THIS TOOL IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND.
THE ENTIRE RISK AS TO THE QUALITY AND PERFORMANCE OF THE PROGRAM IS WITH YOU.
