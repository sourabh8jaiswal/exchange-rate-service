# Exchange Rate Service

A simple Go service to fetch, convert, and cache currency exchange rates using [exchangerate.host](https://exchangerate.host/).

## Features

- Convert currency amounts for a given date
- Get the latest exchange rate between two currencies
- Fetch historical rates for a date range
- Caches rates to reduce API calls

## API Endpoints

### 1. Convert Currency

```
GET /convert?from=USD&to=INR&amount=100&date=2024-05-14
```

**Query Parameters:**
- `from` (required): Source currency code (e.g., USD)
- `to` (required): Target currency code (e.g., INR)
- `amount` (optional): Amount to convert (default: 1)
- `date` (required): Date in `YYYY-MM-DD` format

### 2. Latest Rate

```
GET /latest?from=USD&to=INR
```

**Query Parameters:**
- `from` (required): Source currency code
- `to` (required): Target currency code

### 3. Historical Rates

```
GET /historical?from=USD&to=INR&startDate=2024-05-01&endDate=2024-05-10
```

**Query Parameters:**
- `from` (required): Source currency code
- `to` (required): Target currency code
- `startDate` (required): Start date in `YYYY-MM-DD`
- `endDate` (required): End date in `YYYY-MM-DD`

## Running Locally

### Prerequisites

- [Go](https://golang.org/) 1.19+
- [Docker](https://www.docker.com/) (optional)

### Build and Run (Go)

```sh
go mod tidy
go build -o main .
./main
```

### Build and Run (Docker)

```sh
docker build -t exchange-rate-service .
docker run -p 5050:5050 exchange-rate-service
```

The service will be available at [http://localhost:5050](http://localhost:5050).

## Project Structure

```
handlers/   # HTTP handlers
services/   # Business logic and API integration
cache/      # Caching logic
enums/      # Supported currencies
models/     # Data models
```


## Configuration

- The service uses [exchangerate.host](https://exchangerate.host/) for exchange rates.
- API key and endpoint are set in `services/exchange.go`.
