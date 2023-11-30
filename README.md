# pdf-parser
a bank statement pdf parser which takes pdfs from gmail using google cloud apis and parses the pdf to analyse on the transactions

## API Documentation

This document provides details on three APIs implemented in the application.

## 1. Get All Transactions

### Endpoint

```http
GET /transactions
```

### Description
Retrieves all transactions available in the system.

### Response
Status Code: 200 OK
Content-Type: application/json


```
{
    "transactions": [
        {
            "TxnDate": "2021-08-16T00:00:00Z",
            "TxnDesc": " Payment  -  Credit  Card ",
            "TxnAmount": 5400,
            "TxnBalance": 170400
        },
        {
            "TxnDate": "2021-08-21T00:00:00Z",
            "TxnDesc": " Payment  -  Insurance ",
            "TxnAmount": 3000,
            "TxnBalance": 167400
        },
        ...
    ]
}

```



## 2. Get Transactions by Date
### Endpoint
```http
GET /transactions_by_date
```
### Query Parameters
startDate (required): Date in the format "DD-MM-YYYY".
endDate (required): Date in the format "DD-MM-YYYY".
### Description
Retrieves transactions that occurred on a specific date.

### Example
```
http
GET /transactions_by_date?date=16-08-2021

```
### Response
Status Code: 200 OK
Content-Type: application/json

```
{
    "transactions": [
        {
            "TxnDate": "2021-08-16T00:00:00Z",
            "TxnDesc": " Payment  -  Credit  Card ",
            "TxnAmount": 5400,
            "TxnBalance": 170400
        },
        {
            "TxnDate": "2021-08-21T00:00:00Z",
            "TxnDesc": " Payment  -  Insurance ",
            "TxnAmount": 3000,
            "TxnBalance": 167400
        },
        ...
    ]
}
```

## 2. Get Balance by Date
### Endpoint
```http
GET /balance_by_date
```
### Query Parameters
date (required): Date in the format "DD-MM-YYYY".
### Description
Retrieves the balance of transactions up to a specific date.
### Example
```
http
GET /balance_by_date?date=16-08-2021

```
### Response
Status Code: 200 OK
Content-Type: application/json

```
{
    "balance": 591800
}
```