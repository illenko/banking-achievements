### Producing transaction in Kafka topic:

```shell
curl --location 'http://localhost:8085/transactions' \
--header 'Content-Type: application/json' \
--data '[{
    "ID": "3b241101-e2bb-4255-8caf-4136c566a964",
    "Datetime": "2022-04-01T12:34:56Z",
    "Amount": 123.45,
    "Category": "coffee",
    "Country": "Ukraine"
}]'
```