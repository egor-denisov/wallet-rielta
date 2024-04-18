
# Wallet - тестовое задание от rielta

Перевод денег между клиентами банка со счёта на счёт. Нужно реализовать систему транзакций. 
- есть два клиента банка
- первый клиент делает отправку денег второму клиенту

Как происходит транзакция: 
Идет запрос на сервер от клиента, запрос попадает в очередь, которую должны обрабатывать воркеры и выполнять работу по переводу.

Важно: 
- Избежать потери данных при параллельном выполнении запросов
- Код должен соблюдать принципы SOLID,KISS,DRY 
- предусмотреть различные ситуации остановки\перезагрузки сервера: плановая остановка и перезапуск, падение
- история запросов не должна пропасть при перезапуске сервера
- Для реализации очереди можно использовать RabbitMQ

Что нужно реализовать : 
- БД на postgresql, где будет схема с клиентами и их балансами
- сервер, который проверяет все условия (например: хватает ли денег для совершения операции) и делает изменение баланса (на + или -)
- Вся инфраструктура должна легко подниматься через Docker, написать Dockerfile, docker-compose


## Реализация

Реализован HTTP сервис с 4 эндпоинтами. Документация swagger находится в директории /docs.

- Создание кошелька;
- Перевод средств с одного кошелька на другой;
- Получение историй входящих и исходящих транзакций;
- Получение текущего состояния кошелька.

## Как запустить?

Для запуска приложения в контейнере, необходимо выполнить команду:
```
make build
```

Для запуска вне контейнера:
```
make run
```


## Переменные окружения и конфигурация

Для корректной работы приложения необходимо указать переменные окружения в файле **.env** корневой директории. Ниже представлены сами переменные и краткое описание:

`DISABLE_SWAGGER_HTTP_HANDLER` - при существовании данной переменной Swagger не работает.

`GIN_MODE` - устанавливается в режим debug при необходимости отладки.

`POSTGRES_USER`, `POSTGRES_DB`, `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_PASSWORD` - параметры для инициализации базы данных Postgresql в docker-compose.

`PG_URL` - ссылка для подключения к Postgresql.

`RMQ_URL` - ссылка на очередь rabbitmq.

Также присутствует файл [config.yaml](https://github.com/egor-denisov/wallet-rielta/blob/main/config/config.yml) в котором указываются остальные данные (название и версия приложения, стандартный баланс и др.).

## Архитектура приложения

В папке internal/wallet - логика http сервиса. А в папке internal/walletWorker - обработка записей из очереди и работа с бд.
