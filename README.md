# Бизнес-логика
## Предложения
Предложения создаются пользователями:
- от лица организации\
authorType = Organization

- от лица пользователя\
authorType = User

AuthorId - либо UUID пользователя, либо UUID организации.\
За предложением всегда стоит какая-то организация (по заданию предложения создаются пользователями от организаций).\
Тогда при обращении /bids/my будем отдавать предложения, authorId которых совпадает с id пользователя username, если authorType = User, и те предложения, authorId которых совпадает с id организации, в которой username является ответственным.
