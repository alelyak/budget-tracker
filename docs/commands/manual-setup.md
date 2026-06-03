# Обновление системы
sudo apt update && sudo apt upgrade -y

# Установка программ
sudo apt install имя_пакета -y

# Посмотреть IP
ip a

# Проверить, что порт слушается
sudo netstat -tlnp

# Посмотреть логи службы
sudo journalctl -u имя_службы -f

# Перезагрузка
sudo reboot

# Выключение
sudo shutdown now
