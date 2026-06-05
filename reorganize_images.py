"""
Скрипт для переноса картинок из docs/guides в img/<имя_папки> в корне проекта
и обновления путей в .md файлах внутри docs/guides.
Оригиналы картинок удаляются после успешного копирования.
Поддерживает:
- Картинки с любыми именами (image.png, image-1.png, screenshot 2024.png)
- Ссылки вида ![alt](image.png) в любом .md файле
"""

import re
import shutil
from pathlib import Path


def main():
    # Получаем путь к корню проекта (где лежит скрипт)
    root_dir = Path(__file__).parent.absolute()
    guides_dir = root_dir / "docs" / "guides"
    img_dir = root_dir / "img"

    # Проверяем, существует ли папка docs/guides
    if not guides_dir.exists():
        print(f"❌ Папка не найдена: {guides_dir}")
        print("   Убедитесь, что скрипт запускается из корня проекта,")
        print("   а ваши .md файлы лежат в docs/guides/")
        return

    # Запрашиваем имя папки у пользователя
    folder_name = input(
        "📁 Введите имя папки для картинок (например, manual-setup, manual-postgres): ").strip()
    if not folder_name:
        print("❌ Имя папки не может быть пустым!")
        return

    # Создаём целевую папку img/<folder_name>/
    target_dir = img_dir / folder_name
    target_dir.mkdir(parents=True, exist_ok=True)
    print(f"✅ Создана папка: {target_dir}")

    # Расширения картинок
    image_extensions = {".png", ".jpg", ".jpeg",
                        ".gif", ".svg", ".webp", ".bmp"}

    # Собираем все картинки из папки docs/guides (рекурсивно)
    images = []
    for ext in image_extensions:
        images.extend(guides_dir.rglob(f"*{ext}"))
        images.extend(guides_dir.rglob(f"*{ext.upper()}"))

    if not images:
        print("⚠️ В папке docs/guides не найдено ни одной картинки.")
        print("   Поиск в:", guides_dir)
        return

    print(f"\n📸 Найдено картинок: {len(images)}")

    # Копируем картинки в целевую папку
    copied_images = []
    for img_path in images:
        # Сохраняем оригинальное имя файла (включая пробелы, дефисы и т.д.)
        new_path = target_dir / img_path.name
        try:
            shutil.copy2(img_path, new_path)
            print(f"  📁 Скопирована: {img_path.name}")
            copied_images.append(img_path)
        except Exception as e:
            print(f"  ❌ Ошибка при копировании {img_path.name}: {e}")

    print(f"\n✅ Картинки скопированы в: {target_dir}")
    print("🔄 Обновление ссылок в .md файлах...")

    # Собираем все .md файлы в docs/guides
    md_files = list(guides_dir.rglob("*.md"))

    if not md_files:
        print("⚠️ .md файлы не найдены в docs/guides")
        return

    # Относительный путь от .md файла до папки img/ в корне
    # .md файл в docs/guides/... -> ../../img/...
    new_relative_path_prefix = f"../../img/{folder_name}"
    files_updated = 0
    total_replacements = 0

    for md_file in md_files:
        with open(md_file, 'r', encoding='utf-8') as f:
            content = f.read()

        original_content = content

        # Для каждой картинки обновляем пути в md файле
        for img_path in copied_images:
            img_name = img_path.name
            # Экранируем специальные символы в имени файла (например, пробелы, точки)
            escaped_name = re.escape(img_name)

            # Ищем Markdown ссылки вида:
            # ![alt](image.png)
            # ![alt](image-1.png)
            # ![alt](subfolder/image.png)
            # ![alt](./image.png)
            # НЕ трогаем пути, которые уже ведут в папку img (содержат '/img/')
            pattern = r'(!\[.*?\]\()(?!\.\./\.\./img/)([^)]*?)(' + \
                escaped_name + r')(\))'

            # Формируем новый путь
            replacement = rf'\g<1>{new_relative_path_prefix}/\g<3>\g<4>'

            new_content = re.sub(pattern, replacement, content)
            if new_content != content:
                total_replacements += 1
                content = new_content

        if content != original_content:
            with open(md_file, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"  📝 Обновлён: {md_file.relative_to(root_dir)}")
            files_updated += 1

    if files_updated == 0:
        print("  ⚠️ Ни один .md файл не был обновлён.")
    else:
        print(f"\n✅ Обновлено .md файлов: {files_updated}")
        print(f"✅ Всего замен путей: {total_replacements}")

    # Удаляем оригинальные картинки из docs/guides
    print("\n🗑️ Удаление оригинальных картинок из docs/guides...")
    deleted_count = 0
    for img_path in copied_images:
        try:
            img_path.unlink()
            print(f"  🗑️ Удалён: {img_path.relative_to(root_dir)}")
            deleted_count += 1
        except Exception as e:
            print(f"  ❌ Не удалось удалить {img_path.name}: {e}")

    # Удаляем пустые папки, которые могли остаться после удаления картинок
    print("\n🧹 Очистка пустых папок в docs/guides...")
    empty_folders_deleted = 0
    for folder in sorted(guides_dir.rglob("*"), reverse=True):
        if folder.is_dir() and not any(folder.iterdir()):
            try:
                folder.rmdir()
                print(
                    f"  🧹 Удалена пустая папка: {folder.relative_to(root_dir)}")
                empty_folders_deleted += 1
            except Exception:
                pass

    print("\n" + "="*50)
    print("🎉 ГОТОВО!")
    print("="*50)
    print(f"📸 Скопировано и удалено картинок: {deleted_count}")
    print(f"📝 Обновлено .md файлов: {files_updated}")
    print(f"🧹 Удалено пустых папок: {empty_folders_deleted}")
    print(f"📂 Картинки теперь в: {target_dir}")
    print("\n💡 Не забудьте закоммитить изменения:")
    print("   git add .")
    print("   git commit -m 'docs: move images to img/folder_name'")
    print("   git push")


if __name__ == "__main__":
    main()
