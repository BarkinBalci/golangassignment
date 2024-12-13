[![Go Test](https://github.com/BarkinBalci/golangassignment/actions/workflows/test.yml/badge.svg)](https://github.com/BarkinBalci/golangassignment/actions/workflows/test.yml)

# Golang Assignment

## Building Instructions

```bash
docker build . -t golangassignment
```

## Running Instructions

```bash
docker run -d -p 8080:8080 --env MONGODB_URI="" --env PORT=8080 --env DB_NAME="" golangassignment
```

## Examples in Prod Env

### MongoDB Endpoint

```bash
curl -X POST -H "Content-Type: application/json" -d '{
    "startDate": "2017-01-29",
    "endDate": "2023-02-02",
    "minCount": 2200,
    "maxCount": 3500
}' http://3.72.47.193:8080/mongo
```

### Memory Endpoint

```bash
curl -X POST -H "Content-Type: application/json" -d '{"key": "active-user", "value": "john-doe"}' http://3.72.47.193:8080/memory

curl http://3.72.47.193:8080/memory?key=active-user
```