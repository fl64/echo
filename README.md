echo-http
=========

Simple service for debugging requests and container env

**Localbuild**
```bash
make up
# or
make build
```

**Usage example**
```bash
curl -H "TestHeader: somevalue" localhost:8000 | jq .
```

**Output:**
```json
{
  "host": "localhost:8000",
  "url": "/",
  "method": "GET",
  "headers": {
    "Accept": [
      "*/*"
    ],
    "Testheader": [
      "somevalue"
    ],
    "User-Agent": [
      "curl/7.68.0"
    ]
  },
  "body": "",
  "env": {
    "HOME": "/root",
    "HOSTNAME": "c96c590b45f7",
    "PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
  },
  "hostdata": {
    "args": "/app/server",
    "hostname": "c96c590b45f7"
  },
  "ipaddr": [
    "127.0.0.1",
    "192.169.6.2"
  ],
  "RemoteAddr": "192.169.6.1:39112"
}
```

**Docker image**
```bash
docker pull fl64/echo-http:latest
docker run --rm -p 8000:8000 fl64/echo-http:latest
```
