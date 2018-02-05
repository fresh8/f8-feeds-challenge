# Importer

Environment variables for configuration:

| Environment Variable        | Description           | Default  |
| ------------- |:-------------:| -----:|
| `STORE_ADDR` |  Address of the store | localhost:8001 |
| `STORE_SCHEME` | Scheme for store address | http |
| `FEED_ADDR` | Address of the feed | localhost:8000 |
| `FEED_SCHEME` | Scheme for feed address | http |

### Run tests
`make test-importer`

### Run the application
`make run-importer`

### Used dependencies:
* [envconfig](https://github.com/kelseyhightower/envconfig) - Used for parsing configuration variables. This lib is just so flexible it can fit into any application.
