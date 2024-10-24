# Quote service

Done by Pavlo Khvalygin

## Specifications
#### Database

- ElasticSearch 8.11

#### Libraries

- Uber FX
- GoFiber
- net/http
- ElasticSearch SDK
- joho/godotenv and kelseyhightower/envconfig
- github.com/stretchr/testify

## API Reference

#### 1. Get random quote
 `GET /quote/random`
 

#### 2. Like quote
 `POST /quote/:id/like`

 
#### 3. Find same quotes
 `GET /quote/:id/same`