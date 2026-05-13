#!/bin/bash

# Скрипт рекурсивного обхода файлов (исключая директории out и .git)
# Сохраняет названия и содержимое в выходной файл.
# Использование: ./script.sh [выходной_файл]
# По умолчанию выходной файл: combined_output.txt

OUTPUT_FILE="${1:-combined_output.txt}"

# Создаём выходной файл и получаем его абсолютный путь (чтобы надёжно исключить его)
touch "$OUTPUT_FILE" 2>/dev/null || { echo "Ошибка: не могу создать $OUTPUT_FILE" >&2; exit 1; }
OUTPUT_ABS=$(realpath "$OUTPUT_FILE")
> "$OUTPUT_FILE"  # Очищаем содержимое

# Поиск всех файлов, исключая:
# - директории с именем out (prune)
# - директории с именем .git (prune)
# - сам выходной файл (сравнение абсолютных путей)
find . \( -type d -name out -prune \) -o \
       \( -type d -name .git -prune \) -o \
       -type f -print0 | while IFS= read -r -d '' file; do
    FILE_ABS=$(realpath "$file")
    # Пропускаем выходной файл, если он случайно попал в обход
    [[ "$FILE_ABS" == "$OUTPUT_ABS" ]] && continue

    echo "=== ИМЯ: $file ===" >> "$OUTPUT_FILE"
    if ! cat "$file" >> "$OUTPUT_FILE" 2>/dev/null; then
        echo "Ошибка: не удалось прочитать $file" >> "$OUTPUT_FILE"
    fi
    echo "" >> "$OUTPUT_FILE"   # Разделитель между файлами
done

echo "Готово. Результат сохранён в $OUTPUT_FILE"