gohit
-----

[![Join the chat at https://gitter.im/gohit/Lobby](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gohit/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Run curl commands from yaml files.

Download current version [0.1.0](https://github.com/fabiofalci/gohit/releases)

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
{
  "id": 21726337,
  "name": "sconsify",
  "full_name": "fabiofalci/sconsify",
  "network_count": 14,
  ....
  "subscribers_count": 24
}
```

### yaml file spec

```yaml

# list of files to be imported
files:
  - 'api.yaml'
  - 'api-security-token.yaml'

# list of headers
headers:
  - 'Accept: application/vnd.github.v3+json'

# curl url
url: https://api.github.com

# list of any other curl option
options:
  - '--compress'

# global variables
variables:
  name: value1
  date: value2


# endpoint definitions
endpoints:

  # endpoint name
  get_repo:

    # endpoint path
    path: /repos/{owner}/{repo}

    # http method
    method: GET

    # query string
    query: name={name}&date={date}

    # any other curl option
    options:
      - '--silent'

# request definitions
requests:

  # request name
  show_sconsify:

    # endpoint name
    endpoint: get_repo

    # local variables
    owner: fabiofalci
    repo: sconsify
```

### Environments

You can define one basic api file and then import it from different environment files:

`api.yaml`

```yaml
url: https://{env}.my-api.com

endpoints:
  get_something:
    path: /something/{id}
```


`api-staging.yaml`

```yaml
files:
  - 'api.yaml'

variables:
  env: dev

requests:
  get_something_123:
    endpoint: get_something
    id: 123
```

`api-uat.yaml`

```yaml
files:
  - 'api.yaml'

variables:
  env: uat

requests:
  get_something_456:
    endpoint: get_something
    id: 456
```
