# MusicLib
Репозиторий [тестового задания](https://docs.google.com/document/d/1zs_jBYmvFwXnJo6zU6gePzY784SJGGFx-UXuKDFzRV0/edit?usp=sharing) для Effective Mobile.

Для контейнеризации базы данных PostgreSQL и самого сервера используется [Docker Compose](https://github.com/lynxbites/musiclib/blob/main/compose.yaml), оба контейнера используют Alpine Linux.

API полностью описано сваггером, сваггер хостится вместе с самим API, доступ к нему можно получить перейдя по http://localhost:8001/doc/index.html, документация swagger сгенерирована с помощью [swag](https://github.com/swaggo/swag). 

Спасибо что уделили внимание! c:
## Usage
Клонировать репозиторий:

    git clone https://github.com/lynxbites/musiclib

Запустить Docker Compose:

    make compose
## Libraries
[github.com/go-chi/chi](https://github.com/go-chi/chi) - Удобный и простой роутер.\
[github.com/jackc/pgx](https://github.com/jackc/pgx) - Нативный драйвер для PostgreSQL.\
[github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate) - Библиотека для миграции баз данных SQL.\
[github.com/swaggo/swag](https://github.com/swaggo/swag) - Генерация документации.\
[github.com/charmbracelet/log](https://github.com/charmbracelet/log) - Логгер.\
[github.com/joho/godotenv](https://github.com/joho/godotenv) - Загружает переменные среды из .env файла.

