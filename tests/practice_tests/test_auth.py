from conftest import call, expect, transaction


def test_unauthenticated_transactions():
    expect(call('GET', '/transactions'), 401)


def test_transactions_are_private(account_factory):
    alice, bob = account_factory(), account_factory()
    a = transaction(alice, 'Продукты', 100)['id']
    b = transaction(bob, 'Продукты', 200)['id']
    for user, own, other in [(alice, a, b), (bob, b, a)]:
        rows = expect(call('GET', '/transactions', user['token']), 200)
        assert {row['id'] for row in rows} == {own}
        assert other not in {row['id'] for row in rows}
        assert all(row['user_id'] == user['user']['id'] for row in rows)


def test_cannot_delete_foreign_transaction(account_factory, db):
    owner, stranger = account_factory(), account_factory()
    tx = transaction(owner, 'Продукты', 100)['id']
    expect(call('DELETE', f'/transactions/{tx}', stranger['token']), 404)
    assert db.execute('SELECT user_id, amount FROM transactions WHERE id = %s', (tx,)).fetchone() == (owner['user']['id'], 100)
