from datetime import date
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from conftest import FRONTEND, PASSWORD, transaction


def button(browser, text):
    return WebDriverWait(browser, 15).until(EC.element_to_be_clickable((By.XPATH, f"//button[normalize-space(.)='{text}']")))


def test_invalid_email(browser):
    browser.get(FRONTEND + '/login')
    browser.find_element(By.CSS_SELECTOR, 'input[placeholder="user@example.com"]').send_keys('invalid-email')
    browser.find_element(By.CSS_SELECTOR, 'input[type=password]').send_keys(PASSWORD)
    button(browser, 'Войти').click()
    WebDriverWait(browser, 15).until(EC.visibility_of_element_located((By.XPATH, "//*[normalize-space(text())='Некорректный email']")))
    assert browser.current_url.endswith('/login')
    assert browser.execute_script("return localStorage.getItem('token')") is None


def test_switch_registration_form(browser):
    browser.get(FRONTEND + '/login')
    button(browser, 'Нет аккаунта? Зарегистрироваться').click()
    WebDriverWait(browser, 15).until(EC.visibility_of_element_located((By.XPATH, "//h4[normalize-space(.)='Создайте аккаунт']")))
    assert button(browser, 'Зарегистрироваться').is_displayed()
    button(browser, 'Уже есть аккаунт? Войти').click()
    WebDriverWait(browser, 15).until(EC.visibility_of_element_located((By.XPATH, "//h4[normalize-space(.)='Войдите в аккаунт']")))
    assert button(browser, 'Войти').is_displayed()


def test_dashboard_displays_amounts(browser, account):
    today = date.today()
    transaction(account, 'Зарплата', 800, today)
    transaction(account, 'Продукты', 300, today)
    browser.get(FRONTEND + '/login')
    browser.find_element(By.CSS_SELECTOR, 'input[placeholder="user@example.com"]').send_keys(account['user']['email'])
    browser.find_element(By.CSS_SELECTOR, 'input[type=password]').send_keys(PASSWORD)
    button(browser, 'Войти').click()
    WebDriverWait(browser, 15).until(EC.url_contains('/dashboard'))
    for label, expected in [('Баланс', '500 ₽'), ('Доходы (месяц)', '+800 ₽'), ('Расходы (месяц)', '-300 ₽')]:
        locator = (By.XPATH, f"//p[normalize-space(.)='{label}']/following-sibling::p[1]")
        element = WebDriverWait(browser, 15).until(EC.visibility_of_element_located(locator))
        assert element.text.replace('\u00a0', ' ') == expected
