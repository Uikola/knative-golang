## Начало Работы

Чтобы запустить приложение следуйте следующим шагам.

### Запуск

1. Клонируйте репозиторий
   ```sh
   git clone https://github.com/Uikola/knative-golang.git
   ```

2. Запустите билд докер файла
   ```sh
   docker build -t app -f Dockerfile .
   ```

3. Запустите готовый образ
   ```sh
   docker run -p 8000:8000 app:latest
   ```

4. Приложение готово к использованию!