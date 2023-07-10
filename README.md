
# **Loan Mabagenebt** :moneybag:

## **About The Project**

This is a study project that simulates an api that manage loans. A previous knowledge about finances, loans and mathematics are necessary to fully understand this application

## **Getting Started**

### **Usage**

To test this application first you need to start the docker-compose

```bash
docker-compose up -d 
```

Then you can run the API

```bash
go run .
```

At this moment, this api expose the following routes: 

- **POST** `/funding-calculator/simulate` Simulates a Loan. It will only calculate the values and return
- **POST** `/funding-calculator/contract` Calculate and creates a new loan 
- **GET** `/funding-calculator/find` Find a loan by its ID
- **GET** `/funding-calculator/find-all` Find all loans