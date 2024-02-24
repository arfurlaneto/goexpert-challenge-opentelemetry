# What is this?

This a service that returns temperatures (celsius, fahrenheit and kelvin) for a Brazilian CEP.

It is also one of the challenges of https://goexpert.fullcycle.com.br/pos-goexpert/.

# How to run it?

Run with docker compose:

```bash
docker compose up
```

And do a request, something like this:

```bash
curl --request POST --url 'http://localhost:8080' -H "Content-Type: application/json" -d '{"cep" : "69400970"}'
```
