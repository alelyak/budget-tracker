from datetime import date
from conftest import call, expect, transaction, PASSWORD


def test_registration_creates_default_categories(account_factory, db):
    user = account_factory()
    uid = user['user']['id']
    assert expect(call('GET', '/ping', user['token']), 200)['user_id'] == uid
    row = db.execute('SELECT email, password_hash FROM users WHERE id = %s', (uid,)).fetchone()
    assert row[0] == user['user']['email']
    assert row[1] != PASSWORD and row[1].startswith(('$2a$', '$2b$', '$2y$'))
    expected = {('Продукты', 'expense', 15000), ('Транспорт', 'expense', 5000), ('Кафе', 'expense', 3000), ('Развлечения', 'expense', 4000), ('Зарплата', 'income', 0), ('Фриланс', 'income', 0)}
    rows = db.execute('SELECT name, type, monthly_limit FROM categories WHERE user_id = %s', (uid,)).fetchall()
    assert len(rows) == 6 and set(rows) == expected
    api_rows = expect(call('GET', '/categories', user['token']), 200)
    assert len(api_rows) == 6
    assert {(r['name'], r['type'], r['monthly_limit']) for r in api_rows} == expected


def test_income_expense_update_balance(account, db):
    day = date.today()
    income = transaction(account, 'Зарплата', 10000, day)['id']
    expense = transaction(account, 'Продукты', 2500, day)['id']
    rows = db.execute('SELECT id, amount FROM transactions WHERE user_id = %s', (account['user']['id'],)).fetchall()
    assert dict(rows) == {income: 10000, expense: 2500}
    stats = expect(call('GET', '/dashboard/stats', account['token'], params={'month': day.strftime('%Y-%m')}), 200)
    assert (stats['total_income'], stats['total_expense'], stats['balance']) == (10000, 2500, 7500)


def test_budget_warning_and_persistence(account, db):
    day = date.today()
    first = transaction(account, 'Транспорт', 4000, day)
    assert 'budget_warning' not in first
    second = transaction(account, 'Транспорт', 2000, day)
    warning = second['budget_warning']
    assert warning['category_name'] == 'Транспорт'
    assert (warning['limit'], warning['spent'], warning['exceed_amount']) == (5000, 6000, 1000)
    assert db.execute('SELECT amount FROM transactions WHERE id = %s AND user_id = %s', (second['id'], account['user']['id'])).fetchone() == (2000,)
    stats = expect(call('GET', '/dashboard/stats', account['token'], params={'month': day.strftime('%Y-%m')}), 200)
    budget = next(c for c in stats['category_budgets'] if c['name'] == 'Транспорт')
    assert budget['is_exceeded'] is True
    assert budget['spent'] == 6000
