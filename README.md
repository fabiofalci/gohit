gohit
-----

Run curl commands from yaml files.

* Define your API:

```yaml
headers:
  - 'Accept: application/vnd.github.v3+json'

url: https://api.github.com

options:
  - '--compress'

endpoints:
  get_repo:
    path: /repos/{owner}/{repo}
  get_user:
    path: /users/{username}
    options:
      - '--silent'

requests:
  show_sconsify:
    endpoint: get_repo
    owner: fabiofalci
    repo: sconsify
  show_fabio:
    endpoint: get_user
    username: fabiofalci
```

* Show requests:

```
$ gohit -f github.yaml requests
Request show_fabio:
curl https://api.github.com/users/fabiofalci \
        -H 'Accept: application/vnd.github.v3+json' \
        --compress \
        --silent \
        -XGET

Request show_sconsify:
curl https://api.github.com/repos/fabiofalci/sconsify \
        -H 'Accept: application/vnd.github.v3+json' \
        --compress \
        -XGET

```

* Execute a request:

```
$ gohit -f github.yaml run show_sconsify
show_sconsify:
{
  "id": 21726337,
  "name": "sconsify",
  "full_name": "fabiofalci/sconsify",
  "network_count": 14,
  ....
  "subscribers_count": 24
}

```