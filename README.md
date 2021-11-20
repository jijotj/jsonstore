# jsonstore

A simple HTTP service that stores and returns dynamic json configurations. 

Following are the endpoints implemented:

| Name   | Method      | URL
| ---    | ---         | ---
| List   | `GET`       | `/configs`
| Create | `POST`      | `/configs`
| Get    | `GET`       | `/configs/{name}`
| Update | `PUT/PATCH` | `/configs/{name}`
| Delete | `DELETE`    | `/configs/{name}`
| Query  | `GET`       | `/search?metadata.key=value`


### Query example:

```sh
curl http://config-service/search?metadata.monitoring.enabled=true
```

### Response:

```json
[
  {
    "name": "dc-1",
    "metadata": {
      "monitoring": {
        "enabled": "true"
      },
      "limits": {
        "cpu": {
          "enabled": "false",
          "value": "300m"
        }
      }
    }
  },
  {
    "name": "dc-2",
    "metadata": {
      "monitoring": {
        "enabled": "true"
      },
      "limits": {
        "cpu": {
          "enabled": "true",
          "value": "250m"
        }
      }
    }
  },
]
```

### Commands to set up and start the jsonstore http server:

- `make setup`: install dependencies in the local machine
- `make build`: builds the code and outputs a _jsonstore_ binary
- `make docker.build`: builds a docker image based out of alpine linux for the server
- `make test`: runs all tests and outputs the test coverage
- `make lint`: runs lint checks on the source code
- `make fix`: auto fixes lint issues in the source code
- `make mocks.regenerate`: auto generates mockery mocks used for unit testing

### Deploy to a kubernetes cluster

Run the following command to deploy the server to a kubernetes cluster:
```
helm template --name jsonstore --namespace <NAMESPACE_NAME> -f deploy/values/local.yaml deploy/ | kubectl apply -f -
```

### Example curls:

- Create a new config:
```
curl --request POST --data '{"name":"datacenter-1","metadata":{"monitoring":{"enabled":"true"},"limits":{"cpu":{"enabled":"false","value":"300m"}}}}' <SERVER_ADDRESS>/configs
``` 
- List all configs:
```
curl <SERVER_ADDRESS>/configs
```
- List a specific config:
```
curl <SERVER_ADDRESS>/configs/datacenter-1
```
- Search for configs:
```
curl <SERVER_ADDRESS>/configs/search?metadata.monitoring.enabled=true                                                                                                     
```
- Get application metrics:
```
curl <SERVER_ADDRESS>/metrics                                                                                                     
```
