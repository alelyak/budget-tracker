import pytest
from psycopg import errors
from conftest import category, transaction


def test_unique_email(db, account):
    with pytest.raises(errors.UniqueViolation):
        db.execute('INSERT INTO users(email, password_hash) VALUES (%s, %s)', (account['user']['email'], 'test-hash'))
    assert db.execute('SELECT count(*) FROM users WHERE email = %s', (account['user']['email'],)).fetchone()[0] == 1


def test_negative_amount_rejected(db, account):
    cat = category(account, 'Продукты')
    with pytest.raises(errors.CheckViolation):
        db.execute('INSERT INTO transactions(user_id, category_id, amount, date, note) VALUES (%s, %s, %s, CURRENT_DATE, %s)', (account['user']['id'], cat['id'], -100, 'invalid'))
    assert db.execute('SELECT count(*) FROM transactions WHERE user_id = %s', (account['user']['id'],)).fetchone()[0] == 0


def test_category_delete_cascades(db, account):
    cat = category(account, 'Продукты')
    tx = transaction(account, 'Продукты', 100)['id']
    other = transaction(account, 'Транспорт', 200)['id']
    assert db.execute('SELECT id FROM transactions WHERE id = %s', (tx,)).fetchone() == (tx,)
    assert db.execute('DELETE FROM categories WHERE id = %s', (cat['id'],)).rowcount == 1
    assert db.execute('SELECT id FROM transactions WHERE id = %s', (tx,)).fetchone() is None
    assert db.execute('SELECT id FROM transactions WHERE id = %s', (other,)).fetchone() == (other,)
