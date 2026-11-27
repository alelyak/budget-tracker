# Практическая работа: тестирование бюджет-трекера

12 независимых тестов pytest: 3 UI, 3 права доступа, 3 SQL, 3 интеграционных. Исходный код приложения менять не требуется. SQL-тесты заменяют ORM-тесты, поскольку проект использует database/sql; согласуйте это с преподавателем.

## Запуск (Windows PowerShell)

1. Распакуйте папку practice_tests в корень budget-tracker рядом с finance-backend и finance-frontend. Старую папку tests не удаляйте.
2. Запустите PostgreSQL. Через pgAdmin создайте отдельную пустую БД finance_test с владельцем, указанным в DB_USER файла finance-backend/.env. Остановите ранее запущенный backend.
3. В терминале из корня проекта:

```powershell
cd finance-backend
$env:DB_NAME = 'finance_test'
$env:PORT = '8080'
go run ./cmd/api
```

Остальные настройки подключения берутся из существующего .env. Переменные окружения имеют приоритет над .env при godotenv.Load(). Сервер сам создаст таблицы. Не запускайте тесты против обычной БД.

4. Во втором терминале из корня:

```powershell
cd finance-frontend
npm install
npm run dev -- --host localhost --port 5173 --strictPort
```

5. В третьем терминале из корня (нужны Python 3.10+ и установленный Chrome):

```powershell
py -m venv .venv-tests
.\.venv-tests\Scripts\python.exe -m pip install -r practice_tests/requirements.txt
$env:TEST_DB_DSN = 'host=localhost port=5432 dbname=finance_test user=postgres password=YOUR_PASSWORD'
.\.venv-tests\Scripts\python.exe -m pytest practice_tests -v
```

Подставьте реальные DB_HOST, DB_PORT, DB_USER, DB_PASSWORD из .env. При пробелах в значениях используйте кавычки по правилам libpq. Не публикуйте пароль. Selenium Manager подбирает ChromeDriver (может потребоваться интернет). Для показа браузера задайте `$env:HEADLESS = '0'` перед запуском. На Linux/macOS используйте python3, .venv-tests/bin/python и export вместо $env:.

Отдельная группа: `python -m pytest practice_tests/test_auth.py -v` (используйте Python из виртуального окружения). HTML-отчёт не обязателен; `--junitxml=results.xml` создаёт XML-отчёт. Запуск только practice_tests исключает старые зависимые Selenium-тесты.

## Обоснование выбора

| Тест | Проверка и важность |
|---|---|
| test_invalid_email | Валидация формы и отсутствие входа при ошибке. |
| test_switch_registration_form | Переключение состояний формы в обе стороны. |
| test_dashboard_displays_amounts | Отображение доходов, расходов и баланса, подготовленных через API. |
| test_unauthenticated_transactions | Защита API без токена, HTTP 401. |
| test_transactions_are_private | Изоляция списков двух пользователей. |
| test_cannot_delete_foreign_transaction | HTTP 404 при попытке удаления чужой записи, запись сохранена. |
| test_unique_email | Ограничение UNIQUE на уровне PostgreSQL, SQLSTATE 23505. |
| test_negative_amount_rejected | Ограничение CHECK, SQLSTATE 23514. |
| test_category_delete_cascades | ON DELETE CASCADE; удаляется связанная запись, другая сохраняется. |
| test_registration_creates_default_categories | API, JWT, хеш пароля и шесть записей категорий. |
| test_income_expense_update_balance | API, сохранение транзакций, расчёт месячных итогов. |
| test_budget_warning_and_persistence | Сервис проверки лимита, сохранение расхода и отображение превышения в статистике. |

Альтернативы: UI — пустые поля, ошибки сервера, модальные окна; права — просроченный JWT, чужая категория, редактирование чужой операции; БД — NOT NULL, неверный внешний ключ, точность DECIMAL; интеграция — обновление баланса после удаления, CSV, фильтры по датам. Выбранные проверки покрывают разные механизмы, а не только разные входные значения.

Фикстуры создают уникальные аккаунты и после теста удаляют только созданные ими записи из finance_test. UI-тесты используют отдельный браузер и явные ожидания. Третий UI-тест является также сквозным тестом; классифицирован по основной цели — отображение данных.

Проверка при подготовке: синтаксис Python и соответствие исходному коду проверены; полный запуск не выполнен, поскольку здесь нет запущенного PostgreSQL и зависимостей Selenium/psycopg. Ожидаемый итог после настройки — 12 passed, не гарантированный результат. Если тест выявит дефект, зафиксируйте его, а не меняйте ожидаемый результат ради прохождения.
