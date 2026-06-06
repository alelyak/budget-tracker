# Проверка версии Go
go version

# Скачивание зависимостей проекта
go mod download

# Сборка бинарника (обычная)
go build -o budget-api ./cmd/api

# Кросс-компиляция (сборка для Linux на ноутбуке)
GOOS=linux GOARCH=amd64 go build -o budget-api ./cmd/api

# Копирование всей папки бэкенда на сервер
scp -r ./finance-backend alelyak@192.168.248.163:/home/alelyak/budget-tracker/

# Копирование одного файла (бинарника)
scp ./budget-api alelyak@192.168.248.163:/home/alelyak/budget-tracker/finance-backend/

# Копирование с переименованием
scp ./budget-api.exe alelyak@192.168.248.163:/home/alelyak/budget-tracker/finance-backend/budget-api

# Создание папки
mkdir -p /home/alelyak/budget-tracker

# Просмотр содержимого папки
ls -la

# Перемещение/копирование файла
sudo cp /home/alelyak/budget-tracker/finance-backend/budget-api /opt/budget-tracker/

# Копирование .env файла
sudo cp /home/alelyak/budget-tracker/finance-backend/.env /opt/budget-tracker/

# Редактирование файла
nano /home/alelyak/budget-tracker/finance-backend/go.mod
nano /home/alelyak/budget-tracker/finance-backend/.env
sudo nano /etc/systemd/system/budget-api.service

# Просмотр содержимого файла
cat .env

# Создание папки для бинарника
sudo mkdir -p /opt/budget-tracker

# Перечитать unit-файлы (после создания/изменения сервиса)
sudo systemctl daemon-reload

# Добавить сервис в автозапуск
sudo systemctl enable budget-api

# Запустить сервис
sudo systemctl start budget-api

# Остановить сервис
sudo systemctl stop budget-api

# Проверить статус сервиса
sudo systemctl status budget-api

# Посмотреть логи сервиса
sudo journalctl -u budget-api -f
sudo journalctl -u budget-api -n 50

# Подключение к PostgreSQL от пользователя postgres
sudo -u postgres psql

# Просмотр списка баз данных (внутри psql)
\l

# Просмотр списка пользователей (внутри psql)
\du

# Смена пароля пользователя (внутри psql)
ALTER USER budget_user WITH PASSWORD 'новый_пароль';

# Выдача прав на схему public (внутри psql)
GRANT CREATE ON SCHEMA public TO budget_user;
GRANT USAGE ON SCHEMA public TO budget_user;

# Выход из psql
\q

# Прямой SQL-запрос из командной строки
psql -d budget_db -U budget_user -W -c "SELECT * FROM users;"

# GET-запрос (без авторизации)
curl http://localhost:8080/ping
curl http://localhost:8080/api/ping

# GET-запрос с авторизацией
curl http://localhost:8080/api/categories -H "Authorization: Bearer <TOKEN>"

# POST-запрос с JSON-телом (регистрация)
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"123456"}'

# POST-запрос (логин)
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

# Проверка, что порт слушается
curl http://localhost:8080

# Запуск бинарника вручную (терминал зависает)
./budget-api

# Остановка ручного запуска
Ctrl + C

# Просмотр запущенных процессов
ps aux | grep budget-api

# Просмотр портов
sudo ss -tlnp | grep 8080

# Проверка, доступен ли сервер
ping 192.168.248.163

# Проверка, открыт ли порт
nc -zv localhost 8080

# Последние строки логов сервиса
sudo journalctl -u budget-api -n 20

# Логи в реальном времени
sudo journalctl -u budget-api -f

go version                           # версия Go
go build -o имя ./cmd/путь           # сборка бинарника
scp -r папка user@ip:/путь/на/сервере # копирование на сервер
sudo systemctl status budget-api     # статус сервиса
sudo systemctl restart budget-api    # перезапуск сервиса
sudo journalctl -u budget-api -f     # логи сервиса
curl http://localhost:8080/api/ping  # проверка API
