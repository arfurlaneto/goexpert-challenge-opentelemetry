# What is this?

This is a service that returns temperatures (celsius, fahrenheit and kelvin) and city name for a Brazilian CEP.

There are two involved services and both are instrumented with OpenTelemetry.

It is also one of the challenges of https://goexpert.fullcycle.com.br/pos-goexpert/.

# How to run it?

Run with docker compose (ports 8080 and 9411 must be free to bind):

```bash
docker compose up
```

And do a request, something like this:

```bash
curl --request POST --url 'http://localhost:8080' -H "Content-Type: application/json" -d '{"cep" : "69400970"}'
```

You can open Zipkin UI at `http://localhost:9411/` to check out application tracing.