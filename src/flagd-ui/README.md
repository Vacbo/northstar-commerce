# Flagd-ui

This application provides a user interface for configuring the feature
flags of the flagd service.

This is a [Phoenix](https://www.phoenixframework.org/) project.

## Running the application

The application can be run with the rest of the Northstar Commerce stack using
the Docker Compose or Make targets from the repository root.

## Local development

* Run `mix setup` to install and setup dependencies
* Create a `data` folder: `mkdir data`.
* Copy [../flagd/feature-flags.json](../flagd/feature-flags.json) to `./data/feature-flags.json`
  * `cp ../flagd/feature-flags.json ./data/feature-flags.json`
* Start Phoenix endpoint with `mix phx.server` or inside IEx with `iex -S mix phx.server`

Now you can visit `localhost:4000` from your browser.

## Programmatic use through the API

This service exposes a REST API to ease its usage in a programmatic way for
power users.

You can read the current configuration using this HTTP call:

```json
$ curl localhost:8080/feature/api/read | jq

{
  "flags": {
    "ad_delivery_retry_sampling": {
      "defaultVariant": "off",
      "description": "Fail ad service",
      "state": "ENABLED",
      "variants": {
        "off": false,
        "on": true
      }
    },
    "ad_personalization_burst": {
      "defaultVariant": "off",
      "description": "Triggers high cpu load in the ad service",
      "state": "ENABLED",
      "variants": {
        "off": false,
        "on": true
      }
    },
    "ad_memory_reclamation": {
      "defaultVariant": "off",
      "description": "Triggers full manual garbage collections in the ad service",
      "state": "ENABLED",
      "variants": {
        "off": false,
        "on": true
      }
    },
    ...
  }
}
```

You can also write a new settings file by sending a new configuration inside
the `data` field of a POST request body.

Bear in mind that _all_ the data will be rewritten by this write operation.

```sh
$ curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"data": {"$schema":"https://flagd.dev/schema/v0/flags.json","flags":{"ad_delivery_retry_sampling":{"defaultVariant":"on","description":"Fail ad service","state":"ENABLED","variants":{"off":false,"on":true}}...' \
  http://localhost:8080/feature/api/write
```

In addition to the `/read` and `/write` endpoint, these file-oriented endpoints
are available for operational tooling:

* `/read-file` (`GET`)
* `/write-to-file` (`POST`)
