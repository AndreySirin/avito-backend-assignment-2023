# avito-backend-assignment-2023
[Тестовое задание от Авито](./doc/task.md). Backend-trainee-assignment-2023. Сервис динамического сегментирования пользователей

## Необходимо для запуска
- [Taskfile](https://taskfile.dev/installation/)

## Уставновка 
```shell
git clone git@github.com:AndreySirin/avito-backend-assignment-2023.git
```
## Запуск проекта
1. Выполнить команду в терминале `task dev-tools:install` 
2. Далее `task up` 
3. удаление контейнеров `task down`

# Методы API
## Создание юзера
Реализованы CRUD методы для сущностей: users, segments. Ниже представлены описания только тех, которые требовались в описании тестового задания.
```shell
метод:POST
URL:http://localhost:8080/api/v1/users
Body:
{
  "full_name": "Иванов Иван",
  "gender": "male",
  "date_of_birth": "1995-05-01"
}
ответ:  "id": "7fe20999-8318-492e-aefc-1447f6012a27"
stutus 201
```
## Создание сегмента 
```shell
метод:POST
URL:http://localhost:8080/api/v1/segments
Body:
{
  "title": "Активные пользователи",
  "description": "Пользователи, совершившие покупки за последний месяц",
  "auto_user_prc": 0
}
ответ:  "id": "87825a4e-95a8-4953-831a-ef94d9df5598"
stutus 201
комментарий: Для данного метода так же реализовано автоматическое добавление пользователя в сегмент. Значение "auto_user_prc" это процент от общего числа существующий юзеров. Пользователи добавляются рандомно.
```
## Удаление сегмента 
```shell
метод:DELETE
URL:http://localhost:8080/api/v1/segments/{id}
Body:
{}
ответ:  "Sussefull delete"
stutus 200
```
## Метод добавления пользователя в сегмент
```shell
метод:POST
URL:http://localhost:8080/api/v1/subscription
Body:
{
    "id_user": "7fe20999-8318-492e-aefc-1447f6012a27",
    "name_segment": ["Активные пользователи"],
    "ttl": ["2025-06-01 10:00:00"],
    "auto_added": [false]
}
ответ:  "7fe20999-8318-492e-aefc-1447f6012a27"
stutus 200
комментарий: Так же реализована возможность добавлять пользователя сразу в несколько сегментов. Нужно лишь перечислить сегменты и их значения в соответствующих строках черех запятую.  
```
## Метод удаление пользователя из сегмента(ов)
```shell
метод:DELETE
URL:http://localhost:8080/api/v1/subscription/{UserID}
Body:
{
    "name_segment": ["Активные пользователи"],
    "ttl": ["2025-06-01 10:00:00"],
    "auto_added": [false]
}
ответ:  "7fe20999-8318-492e-aefc-1447f6012a27"
stutus 200
```
## Метод получения активных сегментов пользователя
```shell
метод:GET
URL:http://localhost:8080/api/v1/userSubscription/{UserID}
Body:
{}
ответ:  
{
  "id": "87825a4e-95a8-4953-831a-ef94d9df5598",
  "title": "Активные пользователи",
  "is_auto_add": true
}
stutus 200
```

## Запрос на историю добавления
```shell
метод:GET
URL:http://localhost:8080/api/v1/history
Body:
{
"time": "2023-05-01 15:04:05"
}
ответ:
{
  "user_id": "7fe20999-8318-492e-aefc-1447f6012a27",
  "title": "Активные пользователи",
  "create": "2025-07-04T17:01:33Z",
  "delete": null
    }
stutus 200
комментарий: В теле запроса мы отправляем время от которого и до time.now выводятся все добавление, удаления.  
```
## Добавление пользователя в сегмент на время
```shell
метод:GET
URL:http://localhost:8080/api/v1/CheckTTLSubscriptions
Body:
{}
ответ:
{}
stutus 200
комментарий: Данный запрос проверяет "ttl" сегмента. Если "ttl"<time.now, то пользователь автоматически удаляется из данного сегмента. 
```