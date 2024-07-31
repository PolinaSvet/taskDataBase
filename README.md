
# DataBases + Go

<div align="center">
	<img src="https://i.ibb.co/hFfTZh9/1.jpg">
</div>


## Структура программы:



## Требования к системе:

1) Сервер GoNews, который мы разрабатываем, должен предоставлять REST API, позволяющий выполнять следующие операции:
- Получение списка всех статей из БД,
- Добавление статьи в БД,
- Обновление статьи в БД,
- Удаление статьи из БД.

2) Сервер должен предоставлять данные в ответ на запросы по протоколу HTTP.
3) Сервер должен использовать характерную для REST API схему запросов:
- Запросы должны приходить на URL, соответствующий коллекции ресурсов. Например, коллекция статей.
- Для обозначения действий над коллекцией должны использоваться методы протокола HTTP: POST для создания ресурса, DELETE для удаления, PUT для обновления и GET для получения данных.
4) Сервер должен хранить всю информацию в базе данных.
5) Сервер должен предоставить как минимум две реализации хранилища данных: одну для реляционной СУБД и одну для документной.
- Объекты статьи должны содержать следующую информацию:
- Идентификатор,
- Имя автора,
- Заголовок,
- Текст,
- Время создания.

Для решения задачи от вас требуется следующее:

1) Разработать схему БД PostgreSQL в форме SQL-запроса. Запрос должен быть помещён в файл schema.sql в корневой каталог проекта.
2) По аналогии с пакетом "memdb" разработать пакет "postgres" для поддержки базы данных под управлением СУБД PostgreSQL.
3) По аналогии с пакетом "memdb" разработать пакет "mongo" для поддержки базы данных под управлением СУБД MongoDB.

Все выше перечисленные требования к системе учтены в проекте.
 

## Revision
 

- 1: PostgreSQL: add tables authors, posts and functions for working with it


## Usage:

**1.Enter this command to start the program:**

**go run server.go -typebd pg -loadbd true**

1) typebd: This parameter is responsible for selecting the database.
- pg - PostgreSQL
- mem - memdb(map)
- mongo - MongoDB

2) loadbd: This parameter determines whether to preload the database from a file or not.
- true - preload the database from a file
- false - not

**go run server.go**

defualt value (-typebd pg -loadbd false)


**2.Open the web browser and go to:**

```sh

http://127.0.0.1:8080/ or  localhost:8080

```

## Authors:

@PolinaSvet

**!!! It is for test now !!!**