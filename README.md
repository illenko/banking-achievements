Implementation of achievements in online banking system.
==============================

### Technologies
- Go
- GoFr
- Kafka
- PostgreSQL
- React
- Docker

### Setup
1. Start PostgreSQL and Kafka: ```cd environment && docker-compose up -d```.
2. Start `transactions-producer` service: ```cd transactions-producer && go run github.com/illenko/transactions-producer```.
3. Start `achievements-service` service: ```cd achievements-service && go run github.com/illenko/achievements-service```.
4. Start `achievements-client` service: ```cd achievements-client && npm run dev```.
5. Open `http://localhost:5173` in your browser.

### Producing transaction in Kafka topic
```shell
curl --location 'http://localhost:8085/transactions' \
--header 'Content-Type: application/json' \
--data '[
    {
        "ID": "3b241101-e2bb-4255-8caf-4136c566a964",
        "Datetime": "2022-04-01T12:34:56Z",
        "Amount": 1.45,
        "Category": "coffee",
        "Country": "Ukraine"
    },
    {
        "ID": "4b241101-e2bb-4255-8caf-4136c566a965",
        "Datetime": "2022-04-02T12:34:56Z",
        "Amount": 24.56,
        "Category": "taxi",
        "Country": "USA"
    },
    {
        "ID": "6b241101-e2bb-4255-8caf-4136c566a967",
        "Datetime": "2022-04-04T12:34:56Z",
        "Amount": 56.78,
        "Category": "electronics",
        "Country": "France"
    },
    {
        "ID": "7b241101-e2bb-4255-8caf-4136c566a968",
        "Datetime": "2022-04-05T12:34:56Z",
        "Amount": 7.89,
        "Category": "restaurants",
        "Country": "Canada"
    }]'
```