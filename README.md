## Описание

Инструкция по запуску как в контейнере.

## Запуск проекта в Docker

Перед запуском необходимо убедиться в наличии тблиц, указанных в ТЗ как: "уже созданные в бд"

Запуск осуществлется с помощью команды:
```
docker build -t tenderer . && docker run -e SERVER_ADDRESS=0.0.0.0:8080 -e POSTGRES_CONN=postgres://{username}:{password}@{host}:{5432}/{dbname}/?sslmode=disable -d -p 8080:8080 tenderer
```

P. S. если бд находитмся на локальном хосте, то в поле host необходимо вписать:
```
host.docker.internal
```
