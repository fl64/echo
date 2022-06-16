echo-http
=========

Simple service for debugging requests, container env and networking

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

**Example output:**
```json
{
  "request": {
    "host": "localhost:8000",
    "url": "/",
    "method": "GET",
    "headers": {
      "Accept": [
        "*/*"
      ],
      "User-Agent": [
        "curl/7.76.1"
      ]
    },
    "body": "",
    "remoteaddr": "172.28.0.1:37520"
  },
  "envs": {
    "env": {
      "HOME": "/root",
      "HOSTNAME": "944803f3b04a",
      "PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    }
  },
  "hostdata": {
    "args": "/app/server",
    "hostname": "944803f3b04a"
  }
}

```

**Docker image**
```bash
docker pull fl64/echo-http:latest
docker run --rm -p 8000:8000 fl64/echo-http:latest
```
