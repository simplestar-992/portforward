# PortForward - TCP Port Forwarder

Simple TCP port forwarder/proxy for redirecting network traffic.

## Usage

```bash
# Forward local port 8080 to remote port 80
./portforward -from-port 8080 -to-host 127.0.0.1 -to-port 80

# With verbose output
./portforward -from-port 8080 -to-host 10.0.0.1 -to-port 80 -v
```

## Options

- `-from-host` - Listen address (default: 127.0.0.1)
- `-from-port` - Local port to listen on
- `-to-host` - Remote host to forward to
- `-to-port` - Remote port
- `-v` - Verbose output