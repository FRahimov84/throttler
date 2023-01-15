# throttler

### Описание задачи

1. Есть внешний сервис, который имеет лимиты вызова API. Он может обработать максимум N запросов за K секунд. При превышении ограничения внешний сервис блокирует запросы на X минут.

2. Необходимо реализовать сервис throttler, который:
   1. Принимает запрос на обработку и возвращает id задачи.
   2. Асинхронно разбирает полученные запросы и с учетом лимитов API внешнего сервиса выполняет их.
   3. Позволяет по id задачи получить ее статус и результат обработки.

3. Реализация должна работать как один процесс, очередь необходимо реализовать без использования вспомогательных систем, таких, как rabbitmq. В качестве БД рекомендуется использовать postgresql. REST/grpc на ваш выбор.

4. Недопустимо терять полученные запросы, если отдали клиенту id задачи.

5. Дополнительный плюс: реализация сервиса с in memory хранилищем в дополнение к postgresql с возможностью через ENV настроить хранилище. Такой вариант, например, может быть развернут в тестовой среде, где потеря данных допустима в случае рестарта сервиса.

6. Дополнительный плюс: сервис должен иметь документацию по использованию и подниматься одной командой со всеми требуемыми зависимостямии.

### Быстрый запуск

Для запуска вам необходим докер (если у вас его нет, скачать его можно с официального сайта [Docker](https://www.docker.com/get-started)):

1. Для запуска приложения выполнить команду
    ```sh
    $ docker-compose up --build throttler-go
    ```
2. Swagger-v1 развернут по [этому](http://localhost:8080/swagger/index.html) адресу.
3. Для запуска с in memory хранилищем(у меня Redis), необходимо записать в переменную окружения ENABLE_REDIS=<strong>true</strong> и это можно сделать в compose.yaml.
4. Запускаем проект по дефоулту, пункт 1.
