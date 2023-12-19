### allow docker sock
`chmod 666 /var/run/docker.sock`

### ignore the cert
/etc/docker/daemon.json
"insecure-registries":["10.44.20.71:30020", "10.44.20.72:30020"]
service docker restart

### delete an image
```
curl -sS -H 'Accept: application/vnd.docker.distribution.manifest.v2+json' \
-o /dev/null \
-w '%header{Docker-Content-Digest}' \
<domain-or-ip>:5000/v2/<repo>/manifests/<tag>
```

```
curl -sS -X DELETE <domain-or-ip>:5000/v2/<repo>/manifests/<digest>
```

```
registry garbage-collect /etc/docker/registry/config.yml
```

