import os
from uuid import uuid4
from datetime import date
import pytest
import requests
import psycopg
from psycopg.conninfo import conninfo_to_dict
from selenium import webdriver
from selenium.webdriver.chrome.options import Options

API = 'http://localhost:8080/api'
FRONTEND = 'http://localhost:5173'
PASSWORD = 'TestPassword123!'


def call(method, path, token=None, **kwargs):
    headers = {'Authorization': f'Bearer {token}'} if token else {}
    return requests.request(method, API + path, headers=headers, timeout=10, **kwargs)


def expect(response, status):
    assert response.status_code == status, f'{response.status_code}: {response.text}'
    return response.json()


@pytest.fixture
def db():
    dsn = os.getenv('TEST_DB_DSN', '')
    if not dsn or conninfo_to_dict(dsn).get('dbname') != 'finance_test':
        pytest.fail('Set TEST_DB_DSN for the separate finance_test database')
    with psycopg.connect(dsn, autocommit=True) as conn:
        assert conn.execute('SELECT current_database()').fetchone()[0] == 'finance_test'
        yield conn


@pytest.fixture
def account_factory(db):
    emails = []
    def create():
        email = f'pytest-{uuid4().hex}@example.com'
        emails.append(email)
        result = expect(call('POST', '/register', json={'email': email, 'password': PASSWORD}), 201)
        assert db.execute('SELECT id FROM users WHERE email = %s', (email,)).fetchone() == (result['user']['id'],), 'Backend must use the same finance_test database'
        return result
    yield create
    for email in emails:
        db.execute('DELETE FROM users WHERE email = %s', (email,))


@pytest.fixture
def account(account_factory):
    return account_factory()


def category(account, name):
    categories = expect(call('GET', '/categories', account['token']), 200)
    return next(c for c in categories if c['name'] == name)


def transaction(account, name, amount, day=None):
    cat = category(account, name)
    return expect(call('POST', '/transactions', account['token'], json={
        'category_id': cat['id'], 'amount': amount,
        'date': (day or date.today()).isoformat(), 'note': 'pytest transaction'
    }), 201)


@pytest.fixture
def browser():
    options = Options()
    options.add_argument('--window-size=1440,1000')
    if os.getenv('HEADLESS', '1') == '1':
        options.add_argument('--headless=new')
    driver = webdriver.Chrome(options=options)
    driver.set_page_load_timeout(20)
    try:
        yield driver
    finally:
        driver.quit()
