**Добрый день >^^<**

Это тестовое задание на должность Junior Go Developer от Гребнева И.Д.  
Для запуска сервера достаточно использовать команду `go run .`  
При включении сервер проверяет, есть ли на устройстве бд "music_library". Если ее нет, то сервер подключится к служебной бд "postgres", создаст бд "music_library", заполнит ее с помощью миграций и подключится к ней.