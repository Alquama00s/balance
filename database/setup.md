
# Steps

1. build databse image
```sh
docker build -t balance/postgres:1 -f database/dockerfile .
```

2. start db container
```sh
docker run --name balance-db -p 5482:5432 -e POSTGRES_PASSWORD=dgfvtygfvt@2534HG -d balance/postgres:1

```