# Work Service

REST API сервис для управления студенческими работами, группами и интеграции с GitBucket.

## Описание

Сервис предоставляет API для:

- Работы с группами студентов
- Поиска и управления студентами через LDAP
- Интеграции с GitBucket (управление репозиториями, коммитами)
- Прокси-запросов к внешним сервисам
- Хранения фотографий студентов

## Технологический стек

- Go 1.25.1
- Gin Web Framework
- LDAP интеграция
- GitBucket API
- Docker & Docker Compose

## Требования

- Go 1.25.1 или выше (для локального запуска)
- Docker и Docker Compose (для запуска в контейнере)
- Доступ к LDAP серверу
- Доступ к GitBucket серверу

## Конфигурация

### Переменные окружения

Создайте файл `.env` на основе `.env.example`:

```bash
cp .env.example .env
```

Заполните следующие переменные:

```env
LDAP_URL=ldap://your-ldap-server:389
WEB_URL=https://your-web-url
GITBUCKET_URL=https://your-gitbucket-url
GITBUCKET_API_KEY=your-api-key
```

### Конфигурация сервера

Конфигурация сервера находится в `config/main.yml`:

```yaml
server:
  port: 8070
  maxHeaderBytes: 1
  readTimeout: 7s
  writeTimeout: 7s

app:
  test: false
```

## Запуск

### Вариант 1: Docker Compose с образом из Docker Hub (рекомендуется)

Проект настроен для использования готового образа из Docker Hub, поэтому локальная сборка не требуется.

1. Убедитесь, что создан файл `.env` с необходимыми переменными (см. раздел "Конфигурация")

2. Запустите сервис:

```bash
docker-compose up -d
```

Docker автоматически скачает образ `airsss/work-svc:latest` из Docker Hub при первом запуске.

Сервис будет доступен по адресу: `http://localhost:8070`

Для остановки сервиса:

```bash
docker-compose down
```

Для просмотра логов:

```bash
docker-compose logs -f work-svc
```

### Вариант 2: Docker run (запуск через Docker Hub образ)

Можно запустить только сервис без docker-compose, используя готовый образ из Docker Hub:

```bash
# Скачать образ (опционально, произойдет автоматически при запуске)
docker pull airsss/work-svc:latest

# Запустить контейнер
docker run -d \
  --name work-svc \
  -p 8070:8070 \
  -e LDAP_URL=ldap://your-ldap-server:389 \
  -e GITBUCKET_URL=https://your-gitbucket-url \
  -e GITBUCKET_API_KEY=your-api-key \
  -e WEB_URL=https://your-web-url \
  -v $(pwd)/photos:/app/photos \
  airsss/work-svc:latest
```

Сервис будет доступен по адресу: `http://localhost:8070`

Для остановки контейнера:

```bash
docker stop work-svc
docker rm work-svc
```

### Вариант 3: Локальный запуск

1. Установите зависимости:

```bash
go mod download
```

2. Запустите сервис:

```bash
go run cmd/main.go
```

### Вариант 4: Сборка и запуск

```bash
go build -o work-svc ./cmd/
./work-svc
```

## API Endpoints

### Health Check

**GET** `/api/ping`

Проверка работоспособности сервиса.

**Ответ:**

```json
{
  "message": "pong",
  "version": "v1"
}
```

---

### Поиск студентов

**POST** `/api/v1/search/students`

Поиск студентов в LDAP по запросу.

**Тело запроса:**

```json
{
  "query": "строка поиска"
}
```

**Ответ:**

```json
{
  "students": [
    {
      "id": "12345",
      "username": "ivanov_ii",
      "photoUrl": "http://localhost:8070/api/photos/ivanov_ii.jpg"
    }
  ],
  "total": 1
}
```

---

### Группы

**GET** `/api/v1/groups/it`

Получить список всех IT групп.

**Ответ:**

```json
{
  "groups": [
    {
      "name": "IT-101"
    },
    {
      "name": "IT-102"
    }
  ],
  "total": 2
}
```

**GET** `/api/v1/groups/:groupName/students`

Получить список студентов указанной группы.

**Параметры:**

- `groupName` - название группы (в URL)

**Ответ:**

```json
{
  "students": [
    {
      "id": "12345",
      "username": "ivanov_ii",
      "photoUrl": "http://localhost:8070/api/photos/ivanov_ii.jpg"
    }
  ],
  "total": 1
}
```

---

### Репозитории GitBucket

**GET** `/api/v1/repos/:owner/contents`

Получить список содержимого репозиториев пользователя (без дат модификации).

**Параметры:**

- `owner` - владелец репозитория (в URL)
- `path` - путь внутри репозитория (query параметр, необязательный)

**Пример:** `/api/v1/repos/student123/contents?path=src`

**Ответ:**

```json
{
  "items": [
    {
      "type": "dir",
      "name": "components",
      "path": "src/components",
      "download_url": "",
      "last_modified": "2024-10-16T10:30:00Z"
    }
  ]
}
```

**GET** `/api/v1/repos/:owner/:repo/contents`

Получить содержимое конкретного репозитория с датами последней модификации.

**Параметры:**

- `owner` - владелец репозитория (в URL)
- `repo` - название репозитория (в URL)
- `path` - путь внутри репозитория (query параметр, необязательный)

**Пример:** `/api/v1/repos/student123/project/contents?path=src`

**Ответ:**

```json
{
  "items": [
    {
      "type": "file",
      "name": "index.html",
      "path": "src/index.html",
      "download_url": "https://gitbucket/...",
      "last_modified": "2024-10-16T10:30:00Z"
    }
  ]
}
```

**GET** `/api/v1/repos/:owner/:repo/commits`

Получить список коммитов репозитория с пагинацией.

**Параметры:**

- `owner` - владелец репозитория (в URL)
- `repo` - название репозитория (в URL)
- `page` - номер страницы (query, по умолчанию: 1)
- `per_page` - количество коммитов на странице (query, по умолчанию: 30, максимум: 100)

**Пример:** `/api/v1/repos/student123/project/commits?page=1&per_page=10`

**Ответ:**

```json
{
  "count": 45,
  "page": 1,
  "per_page": 10,
  "has_next": true,
  "commits": [
    {
      "message": "Initial commit",
      "date": "2024-10-16T10:30:00Z",
      "files": [
        {
          "filename": "README.md"
        }
      ]
    }
  ]
}
```

---

### Прокси для файлов из репозитория

**GET** `/api/repos/:owner/:repo/branches/:ref/raw/*filepath`

Получить raw содержимое файла из репозитория.

**Параметры:**

- `owner` - владелец репозитория
- `repo` - название репозитория
- `ref` - ветка или тег
- `filepath` - путь к файлу

**Пример:** `/api/repos/student123/project/branches/main/raw/src/index.html`

**Ответ:** Содержимое файла с соответствующим Content-Type

**GET** `/api/repos/:owner/:repo/branches/:ref/html/*filepath`

Получить HTML файл с автоматической подстановкой базового URL для относительных путей.

**Параметры:**

- `owner` - владелец репозитория
- `repo` - название репозитория
- `ref` - ветка или тег
- `filepath` - путь к HTML файлу

**Пример:** `/api/repos/student123/project/branches/main/html/index.html`

**Ответ:** HTML контент с обработанными относительными путями (Content-Type: text/html)

---

### Статические файлы

**GET** `/api/photos/*filepath`

Получить фотографии студентов из директории `photos`.

**Пример:** `/api/photos/ivanov_ii.jpg`

## Лицензия

Проект распространяется по лицензии, указанной в файле LICENSE.
