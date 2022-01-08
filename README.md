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
  },
  "routes": [
    {
      "dst": "default",
      "gateway": "172.28.0.1",
      "dev": "eth0"
    },
    {
      "dst": "172.28.0.0/16",
      "dev": "eth0",
      "protocol": "kernel"
    }
  ],
  "ifaces": [
    ...skipped...
    {
      "ifindex": 86,
      "ifname": "eth0",
      "flags": [
        "BROADCAST",
        "MULTICAST",
        "UP",
        "LOWER_UP"
      ],
      "mtu": 1500,
      "operstate": "UP",
      "group": "default",
      "link_type": "ether",
      "address": "02:42:ac:1c:00:02",
      "broadcast": "ff:ff:ff:ff:ff:ff",
      "addr_info": [
        {
          "family": "inet",
          "local": "172.28.0.2",
          "prefixlen": 16,
          "broadcast": "172.28.255.255",
          "scope": "global",
          "label": "eth0",
          "valid_life_time": 4294967295,
          "preferred_life_time": 4294967295
        }
      ]
    }
  ],
  "mounts": [
    ...skipped...
    "proc on /proc/sysrq-trigger type proc (ro,relatime)",
    "tmpfs on /proc/asound type tmpfs (ro,relatime,inode64)",
    "tmpfs on /proc/acpi type tmpfs (ro,relatime,inode64)",
    "tmpfs on /proc/kcore type tmpfs (rw,nosuid,size=65536k,mode=755,inode64)",
    "tmpfs on /proc/keys type tmpfs (rw,nosuid,size=65536k,mode=755,inode64)",
    "tmpfs on /proc/latency_stats type tmpfs (rw,nosuid,size=65536k,mode=755,inode64)",
    "tmpfs on /proc/timer_list type tmpfs (rw,nosuid,size=65536k,mode=755,inode64)",
    "tmpfs on /proc/scsi type tmpfs (ro,relatime,inode64)",
    "tmpfs on /sys/firmware type tmpfs (ro,relatime,inode64)"
  ],
  "resolv.conf": [
    "nameserver 127.0.0.11",
    "options edns0 trust-ad ndots:0"
  ]
}

```

**Docker image**
```bash
docker pull fl64/echo-http:latest
docker run --rm -p 8000:8000 fl64/echo-http:latest
```
