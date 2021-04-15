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
    "HOSTNAME": "709ae3979eba",
    "PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
  },
  "hostdata": {
    "args": "/app/server",
    "hostname": "709ae3979eba"
  },
  "ipaddr": [
    "127.0.0.1",
    "192.169.25.2"
  ]
}
```
