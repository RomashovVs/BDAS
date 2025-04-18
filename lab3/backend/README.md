# САБД - Lab 3
## О задании
Name: Проектирование и реализация клиент-серверного приложения, 
взаимодействующего по HTTPS протоколу с использованием ключей шифрования для SSL/TLS (Two- way TLS)

Author: Andrianov Artemii

Group: #5140904/30202

## Цель
Необходимо написать две программы (клиент и сервер) с использованием библиотеки Spring Boot, 
которые взаимодействуют по HTTPS протоколу с использованием ключей шифрования для SSL/TLS.

## Решение
### About
Так как в Golang нет библиотеки Spring Boot, она будет заменена на стандартный пакет `net/http`.
Для создания ключей будет использован стандартный пакет `crypto`.

### Generate cers
For generate server and client cert/key pairs and CA cert use `make cert`.

### Build
For build service use `make build`. `bin/server` and `bin/client` will be created.

### Tests
For run tests use `make test`.


### Run lab
For run server use:
```bash 
    ./bin/server -config config/config-server.yaml
```

For run client use:
```bash 
    ./bin/server -config config/config-client.yaml
```
