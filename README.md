# http-to-https-proxy
A proxy that upgrades HTTP connections to HTTPS for systems which cannot make HTTPS requests.

## Running the proxy

Download the latest binary corresponding to your platfrom from the [releases section](https://github.com/yeokm1/http-to-https-proxy/releases/).

### Default Configuration

```bash
http-to-https-proxy.exe
2023/06/25 20:38:37 HTTP to HTTPS proxy v0.3 listening to 80, forward to 443 with listening buffer 4096
2023/06/25 20:38:41 Request from 192.168.1.112:13519 to host api.openai.com and url /v1/chat/completions
2023/06/25 20:38:43 EOF reached
2023/06/25 20:38:43 End of handler
...
```

By default, proxy will listen to HTTP requests on port `80` and retransmit HTTPS via port `443`. Buffer size of `4096` is the buffer to receive destination server's response chunks before forwarding back original client.

### Arguments

```
  -h  --help      Print help information
  -l  --listen    HTTP port to listen on. Default: 80
  -c  --connect   HTTPS port to connect to. Default: 443
  -b  --buffer    Buffer size. Default: 4096
  -i  --insecure  Allow insecure TLS certificates. Default: false
  -q  --query     Handle requests from query string "q". Default: false
  -d  --debug     Enable debug console logging. Default: false
```

If the server you are connecting to is using expired/insecure TLS certificates. You can add `-i` argument to allow those connections.

Add `-q` argument to make it act like CGI proxy which accepts target URL in `?q=` query parameter.

# Compiling

Just install the latest Go compiler for your platform. The latest at the time of writing is `1.20.5`. THe following was compiled on windows/amd64 platform using Powershell script `build.ps1`.
