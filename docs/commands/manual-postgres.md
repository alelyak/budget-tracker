# Обновить список доступных пакетов (обязательно перед установкой)
sudo apt update

# Установить пакет с общими файлами PostgreSQL
sudo apt install -y postgresql-common

# Запустить скрипт автоматического добавления официального репозитория
sudo /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh

# Установить пакет с дополнительными модулями
sudo apt install -y postgresql-16 postgresql-contrib-16

# Полностью удалить пакет (с конфигами)
sudo apt purge -y postgresql-16

# Поиск пакета по имени
apt search postgresql

# Проверить статус сервера PostgreSQL
sudo systemctl status postgresql

# Посмотреть список кластеров БД
pg_lsclusters

# Переключиться на системного пользователя postgres
sudo -u postgres

# Запустить psql от имени пользователя postgres
sudo -u postgres psql

# Подключиться к конкретной БД от конкретного пользователя с запросом пароля
sudo -u postgres psql -d budget_db -U budget_user -W

# Выйти из psql
\q

-- Создать базу данных
CREATE DATABASE budget_db;

-- Создать пользователя с паролем
CREATE USER budget_user WITH PASSWORD 'ваш_пароль';

-- Дать все права на БД пользователю
GRANT ALL PRIVILEGES ON DATABASE budget_db TO budget_user;

-- Посмотреть список баз данных
\l

-- Посмотреть список пользователей
\du

-- Подключиться к другой БД (внутри psql)
\c budget_db budget_user

-- Показать текущую базу данных
SELECT current_database();

-- Показать текущего пользователя
SELECT current_user();

# Найти файл pg_hba.conf
sudo find /etc/postgresql -name "pg_hba.conf"

# Отредактировать файл в nano
sudo nano /etc/postgresql/16/main/pg_hba.conf

# Перезагрузить PostgreSQL после изменения конфига
sudo systemctl restart postgresql

# Посмотреть статус сервиса
sudo systemctl status postgresql

# Перезапустить сервис
sudo systemctl restart postgresql

# Добавить в автозапуск
sudo systemctl enable postgresql

# Убрать из автозапуска
sudo systemctl disable postgresql

# Посмотреть все сервисы
sudo systemctl list-units --type=service

# Показать, какие порты слушает система (netstat)
sudo ss -tlnp

# Показать версию PostgreSQL
sudo -u postgres psql -c "SELECT version();"
