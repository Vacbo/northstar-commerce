# Northstar Commerce

Northstar Commerce is a microservice e-commerce application used by the platform team to validate storefront, checkout, fulfillment, and operations workflows under realistic load.

The system is intentionally polyglot because the production estate is polyglot. Each service owns one narrow business capability and publishes telemetry for traces, metrics, and logs so operators can investigate incidents across service boundaries.

## What is included

- Storefront web application and mobile client
- Product catalog, recommendations, cart, checkout, payment, shipping, accounting, fraud, and notification services
- Feature flag control through flagd and the internal flag management UI
- Load generation for repeatable customer journeys
- Docker Compose and Kubernetes manifests for local and cluster deployments
- Telemetry collector configuration for application observability

## Architecture

Requests enter through the frontend proxy and flow through the storefront to backend services over gRPC and HTTP. Customer journeys cover browsing products, viewing recommendations, adding items to cart, checking out, and receiving order confirmations.

The services are designed to be small enough to run locally while still exercising the same operational concerns seen in larger commerce platforms: request fan-out, asynchronous order events, cache dependencies, payment routing, and feature-flagged behavior.

## Quick start

Run the stack locally with Docker Compose:

```bash
docker compose up --build
```

The storefront is available at <http://localhost:8080>. The flag management UI is available at <http://localhost:4000> when the full stack is running.

To return this repository to the pristine baseline used for incident exercises, run:

```bash
./scripts/reset-to-baseline.sh
```

## Service layout

Source code lives under `src/`. Each service has its own Dockerfile, runtime dependencies, and telemetry setup. Shared protobuf definitions are kept in `pb/` and generated into service-specific clients as needed.

## Operations notes

Feature flags live in `src/flagd/feature-flags.json`. Use them to adjust runtime behavior during local testing and incident drills.

Load generation is implemented in `src/load-generator/locustfile.py`. It produces browser-style traffic for common shopper flows and automated QA checks.

## Development

Use the existing service-level tooling for formatting and tests. Before changing cross-service APIs, regenerate protobuf clients and run the affected service tests.

Keep changes focused. This repository is used as an operations target, so every behavior change should be easy to connect to a customer-facing or platform-facing reason.

## License

Northstar Commerce includes components derived from Apache-2.0 licensed upstream work. See `LICENSE` for license terms.
