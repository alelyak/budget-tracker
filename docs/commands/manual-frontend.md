# Проверка статуса бэкенда
sudo systemctl status budget-api

# Проверка, какие порты слушаются
sudo ss -tlnp | grep :80
sudo ss -tlnp | grep :8080

# Просмотр логов nginx в реальном времени
sudo tail -f /var/log/nginx/error.log

# Проверка работы API через curl
curl http://localhost:8080/api/ping
curl -X POST http://192.168.248.163/api/register -H "Content-Type: application/json" -d '{"email":"test@test.com","password":"123456","name":"Test"}'

# Установка nginx
sudo apt install -y nginx

# Включение автозапуска nginx
sudo systemctl enable nginx

# Запуск nginx
sudo systemctl start nginx

# Проверка статуса nginx
sudo systemctl status nginx

# Перезапуск nginx (после изменения конфига)
sudo systemctl restart nginx

# Проверка синтаксиса конфигурации nginx
sudo nginx -t

# Удаление дефолтного сайта
sudo rm /etc/nginx/sites-enabled/default

# Создание символической ссылки для активации сайта
sudo ln -s /etc/nginx/sites-available/budget-tracker /etc/nginx/sites-enabled/

# Редактирование конфига nginx
sudo nano /etc/nginx/sites-available/budget-tracker

# Переход в папку фронтенда
cd /home/alelyak/budget-tracker/finance-frontend

# Сборка React-приложения
npm run build

# Принудительное удаление папки dist (при ошибке прав)
sudo rm -rf /home/alelyak/budget-tracker/finance-frontend/dist

# Смена владельца папки (рекурсивно)
sudo chown -R alelyak:alelyak /home/alelyak/budget-tracker/finance-frontend

# Смена владельца для nginx (если нужно)
sudo chown -R www-data:www-data /home/alelyak/budget-tracker/finance-frontend/dist

# Выдача прав на выполнение для родительских папок
sudo chmod +x /home/alelyak
sudo chmod +x /home/alelyak/budget-tracker
sudo chmod +x /home/alelyak/budget-tracker/finance-frontend

# Выдача прав 755 для папки dist
sudo chmod -R 755 /home/alelyak/budget-tracker/finance-frontend/dist
