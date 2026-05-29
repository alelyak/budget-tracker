from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.keys import Keys
import time
import random

EMAIL = "test@test.ru"
PASSWORD = "123456"

class FinanceTrackerTest:
    def __init__(self):
        self.driver = webdriver.Chrome()
        self.driver.maximize_window()
        self.wait = WebDriverWait(self.driver, 20)
        self.results = []

    def test_login(self):
        print("\n" + "=" * 50)
        print("🔵 ТЕСТ 1: ВХОД В СИСТЕМУ")
        print("=" * 50)

        self.driver.get("http://localhost:5173/login")
        time.sleep(2)

        try:
            inputs = self.wait.until(
                EC.presence_of_all_elements_located((By.TAG_NAME, "input")))
            inputs[0].clear()
            inputs[0].send_keys(EMAIL)
            inputs[1].clear()
            inputs[1].send_keys(PASSWORD)
            print("✅ Логин и пароль введены")

            submit_btn = self.wait.until(EC.element_to_be_clickable(
                (By.XPATH, "//button[@type='submit']")))
            submit_btn.click()

            self.wait.until(EC.url_contains("dashboard"))
            print("✅ Вход выполнен успешно!")
            return True
        except Exception as e:
            print(f"❌ Ошибка входа: {e}")
            return False

    def test_go_to_categories(self):
        print("\n" + "=" * 50)
        print("🔵 ТЕСТ 2: ПЕРЕХОД В КАТЕГОРИИ")
        print("=" * 50)

        try:
            cat_link = self.wait.until(EC.element_to_be_clickable(
                (By.XPATH, "//a[contains(@href, '/categories')]")))
            cat_link.click()
            self.wait.until(EC.url_contains("/categories"))
            print("✅ Переход в раздел категорий выполнен")
            return True
        except Exception as e:
            print(f"❌ Ошибка перехода в категории: {e}")
            return False

    def test_create_category(self):
        print("\n" + "=" * 50)
        print("🔵 ТЕСТ 3: СОЗДАНИЕ КАТЕГОРИИ")
        print("=" * 50)

        category_name = f"ТестКат_{random.randint(100, 999)}"
        monthly_limit = 15000

        try:
            add_btn = self.wait.until(EC.element_to_be_clickable(
                (By.XPATH, "//button[contains(., 'Добавить категорию')]")))
            add_btn.click()
            print("✅ Открыто модальное окно")

            name_input = self.wait.until(EC.element_to_be_clickable(
                (By.XPATH, "//input[contains(@placeholder, 'Название') or contains(@class, 'TextInput-input')]")))
            name_input.click()
            name_input.clear()
            name_input.send_keys(category_name)
            print(f"✅ Название '{category_name}' введено")

            # Добавление месячного лимита
            try:
                limit_input = self.driver.find_element(
                    By.XPATH, "//input[@type='number']")
                limit_input.click()
                limit_input.clear()
                limit_input.send_keys(str(monthly_limit))
                print(f"✅ Месячный лимит {monthly_limit} ₽ установлен")
            except:
                print("⚠️ Поле для ввода лимита не найдено (пропускаем)")

            save_btn = self.wait.until(EC.element_to_be_clickable(
                (By.XPATH, "//button[contains(., 'Создать') or contains(., 'Сохранить')]")))
            save_btn.click()
            print("✅ Категория сохранена")

            time.sleep(1.5)
            self.driver.refresh()
            time.sleep(2)

            if category_name in self.driver.page_source:
                print(f"✅ Категория '{category_name}' успешно добавлена!")
                self.last_created_category = category_name
                return True
            return False

        except Exception as e:
            print(f"❌ Ошибка создания категории: {e}")
            return False

    def test_edit_category(self):
        print("\n" + "=" * 50)
        print("🔵 ТЕСТ 4: РЕДАКТИРОВАНИЕ КАТЕГОРИИ")
        print("=" * 50)

        if not hasattr(self, 'last_created_category'):
            print("⚠️ Нет созданной категории для редактирования")
            return False

        new_name = f"Отредактировано_{random.randint(100, 999)}"

        try:
            edit_xpath = f"//tr[contains(., '{self.last_created_category}')]//button"
            edit_btn = self.wait.until(
                EC.element_to_be_clickable((By.XPATH, edit_xpath)))
            edit_btn.click()
            print("✅ Открыто модальное редактирование")

            name_input = self.wait.until(EC.element_to_be_clickable(
                (By.XPATH, "//input[contains(@placeholder, 'Название') or contains(@class, 'TextInput-input')]")))
            name_input.click()
            name_input.clear()
            name_input.send_keys(new_name)
            print(f"✅ Новое название '{new_name}' введено")

            save_btn = self.wait.until(EC.element_to_be_clickable(
                (By.XPATH, "//button[contains(., 'Сохранить') or contains(., 'Обновить')]")))
            save_btn.click()
            print("✅ Изменения сохранены")

            time.sleep(2)
            self.driver.refresh()
            time.sleep(2)

            if new_name in self.driver.page_source:
                print(f"✅ Категория успешно переименована в '{new_name}'")
                self.last_created_category = new_name
                return True
            else:
                print("❌ Новое имя не найдено")
                return False

        except Exception as e:
            print(f"❌ Ошибка редактирования: {e}")
            self.driver.save_screenshot("edit_error.png")
            return False

    def test_logout(self):
        print("\n" + "=" * 50)
        print("🔵 ТЕСТ 6: ВЫХОД ИЗ АККАУНТА")
        print("=" * 50)

        try:
            self.driver.find_element(
                By.TAG_NAME, "body").send_keys(Keys.ESCAPE)
            time.sleep(0.5)

            logout_btn = self.wait.until(EC.element_to_be_clickable(
                (By.XPATH, "//button[contains(., 'Выйти')]")))
            logout_btn.click()
            print("✅ Выход выполнен")
            return True
        except Exception as e:
            print(f"❌ Ошибка выхода: {e}")
            return False

    def run_all_tests(self):
        print("\n" + "=" * 70)
        print("🧪 ЗАПУСК ТЕСТИРОВАНИЯ ФИНАНСОВОГО ТРЕКЕРА")
        print("=" * 70)

        try:
            self.driver.get("http://localhost:5173")
            print("✅ Приложение доступно")
        except:
            print("❌ Не удалось подключиться!")
            self.driver.quit()
            return

        self.results.append(self.test_login())

        if self.results[0]:
            self.results.append(self.test_go_to_categories())
            self.results.append(self.test_create_category())
            if self.results[-1]:
                self.results.append(self.test_edit_category())
            self.results.append(self.test_logout())

        print("\n" + "=" * 70)
        print("📊 ИТОГИ ТЕСТИРОВАНИЯ")
        print("=" * 70)

        passed = sum(1 for r in self.results if r)
        total = len(self.results)
        print(f"✅ Пройдено: {passed}/{total}")
        print(f"❌ Не пройдено: {total - passed}/{total}")

        if passed == total:
            print("\n🎉 ВСЕ ТЕСТЫ УСПЕШНО ЗАВЕРШЕНЫ!")
        else:
            print("\n⚠️ Некоторые тесты не пройдены")

        self.driver.quit()


if __name__ == "__main__":
    test = FinanceTrackerTest()
    test.run_all_tests()
