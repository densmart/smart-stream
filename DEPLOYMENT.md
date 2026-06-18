# Smart Stream - Deployment Guide

Руководство по развертыванию приложения Smart Stream в Docker контейнерах.

## Предварительные требования

- Docker версии 20.10 или выше
- Docker Compose версии 2.0 или выше
- Доступ к PostgreSQL базе данных (внешней)

## Быстрый старт

### 1. Настройка переменных окружения

Скопируйте файлы с примерами конфигурации и заполните их:

```bash
# Backend configuration
cp backend/.env.example backend/.env

# Frontend configuration
cp frontend/.env.example frontend/.env
```

**Отредактируйте `backend/.env`:**

- **POSTGRES_DSN** - строку подключения к вашей PostgreSQL БД
  - **БД на той же машине**: используйте `host.docker.internal` (Mac/Windows)
  - **БД на отдельном сервере**: используйте IP адрес или hostname сервера
  - **Docker на Linux**: используйте `172.17.0.1` или `host.docker.internal` (Docker 20.10+)
  - ⚠️ **НЕ используйте** `localhost` - это не работает внутри контейнера!
- **JWT_SECRET** - секретный ключ для JWT (сгенерируйте командой: `openssl rand -base64 32`)
- **ADMIN_LOGIN** и **ADMIN_PASSWORD** - учетные данные администратора

Примеры POSTGRES_DSN:
```bash
# БД на той же машине (для разработки на Mac/Windows)
POSTGRES_DSN=postgres://user:password@host.docker.internal:5432/smart_stream?sslmode=disable

# БД на отдельном сервере
POSTGRES_DSN=postgres://user:password@192.168.1.100:5432/smart_stream?sslmode=disable
POSTGRES_DSN=postgres://user:password@db.example.com:5432/smart_stream?sslmode=disable
```

**Отредактируйте `frontend/.env`:**

- **VITE_API_URL** - URL backend API, доступный из браузера
  - ⚠️ **Важно:** URL встраивается в frontend **во время сборки** Docker образа
  - URL должен быть доступен из **браузера пользователя**, а не из контейнера
  - Для локальной разработки: `http://localhost:17777`
  - Для production: укажите внешний IP или домен (например `http://192.168.1.100:17777`)
  - При изменении нужно **пересобрать** frontend образ: `docker compose build frontend`
- **FRONTEND_PORT** - порт для доступа к frontend (по умолчанию 3000)

```bash
# Пример для production сервера
VITE_API_URL=http://192.168.1.100:17777
FRONTEND_PORT=3000

# После изменения VITE_API_URL пересоберите frontend
docker compose build frontend
docker compose up -d
```

### 2. Сборка контейнеров

```bash
docker compose build
```

Эта команда создаст образы:
- `smart-stream:latest` - backend (Go приложение)
- `smart-stream-frontend:latest` - frontend (Vue.js + Nginx)

Для сборки только одного сервиса:
```bash
docker compose build backend   # только backend
docker compose build frontend  # только frontend
```

### 3. Запуск приложения

```bash
docker compose up -d
```

Флаг `-d` запускает контейнеры в фоновом режиме.

### 4. Проверка состояния

Просмотр логов backend:

```bash
docker compose logs -f backend
```

Проверка запущенных контейнеров:

```bash
docker compose ps
```

## Доступные сервисы

После успешного запуска будут доступны следующие сервисы:

- **Frontend (Web UI)**: `http://your-server:${FRONTEND_PORT}` (по умолчанию 3000)
- **Web Admin API**: `http://your-server:${WEB_API_PORT}` (по умолчанию 17777)
- **TV Client API**: `http://your-server:${TV_APP_API_PORT}` (по умолчанию 18888)

### Конфигурация портов

**Backend** (`backend/.env`):
- `WEB_API_PORT` - порт для Web Admin API (по умолчанию 17777)
- `TV_APP_API_PORT` - порт для TV Client API (по умолчанию 18888)

**Frontend** (`frontend/.env`):
- `FRONTEND_PORT` - порт для frontend приложения (по умолчанию 3000)
- `VITE_API_URL` - URL backend API для frontend (встраивается в сборку)

## Управление контейнерами

### Остановка всех сервисов

```bash
docker compose down
```

### Перезапуск сервисов

```bash
docker compose restart
```

### Просмотр логов

```bash
# Все сервисы
docker compose logs -f

# Только backend
docker compose logs -f backend
```

### Обновление приложения

1. Остановите контейнеры:
```bash
docker compose down
```

2. Получите последние изменения:
```bash
git pull
```

3. Пересоберите образы:
```bash
docker compose build --no-cache
```

4. Запустите обновленные контейнеры:
```bash
docker compose up -d
```

## Хранилище данных

Приложение использует bind mounts для хранения медиа-файлов напрямую на хост-машине:

- **./storage/media** - видео и аудио файлы
- **./storage/posters** - изображения постеров

Файлы хранятся в директории `storage/` в корне проекта и напрямую доступны из операционной системы.

### Добавление медиа через FTP

Вы можете настроить FTP-сервер на хост-машине и загружать файлы напрямую в директории:
- `<путь-к-проекту>/storage/media/` - для видео и аудио
- `<путь-к-проекту>/storage/posters/` - для постеров

Файлы станут сразу доступны в контейнере по путям:
- `/app/storage/media/`
- `/app/storage/posters/`

### Резервное копирование

```bash
# Backup всей директории storage
tar czf storage-backup-$(date +%Y%m%d).tar.gz ./storage/

# Backup только media
tar czf media-backup-$(date +%Y%m%d).tar.gz ./storage/media/

# Backup только posters
tar czf posters-backup-$(date +%Y%m%d).tar.gz ./storage/posters/

# Или через rsync на удаленный сервер
rsync -avz ./storage/ user@backup-server:/backups/smart-stream/
```

### Восстановление

```bash
# Restore из tar архива
tar xzf storage-backup-20240618.tar.gz

# Restore через rsync
rsync -avz user@backup-server:/backups/smart-stream/ ./storage/
```

### Права доступа

Убедитесь что у Docker есть права на чтение/запись в директорию `storage/`:

```bash
# Установить корректные права
chmod -R 755 ./storage/
```

## Миграции базы данных

Миграции автоматически применяются при запуске backend контейнера.

Если нужно применить миграции вручную:

```bash
docker compose exec backend ./start_point migrate
```

## Troubleshooting

### Контейнер не запускается

Проверьте логи:
```bash
docker compose logs backend
```

### Проблемы с подключением к БД

1. Убедитесь что POSTGRES_DSN корректный
2. Проверьте доступность БД с сервера:
```bash
docker compose exec backend ping your-db-host
```

### Проблемы с портами

Если порты уже заняты, измените их в соответствующих .env файлах:

**Backend порты** (`backend/.env`):
```bash
WEB_API_PORT=18000  # Измените на свободный порт
TV_APP_API_PORT=18001  # Измените на свободный порт
```

**Frontend порт** (`frontend/.env`):
```bash
FRONTEND_PORT=8080  # Измените на свободный порт
```

После изменения пересоздайте контейнеры:
```bash
docker compose down
docker compose up -d
```

## Безопасность

1. **ВАЖНО**: Не коммитьте `.env` файлы в git!
   - `backend/.env` - содержит credentials для БД
   - `frontend/.env` - может содержать чувствительные URL
2. Используйте сильный пароль для ADMIN_PASSWORD
3. Генерируйте уникальный JWT_SECRET для каждого окружения
4. Настройте firewall для ограничения доступа к портам
5. Используйте SSL/TLS для production окружения (рекомендуется nginx reverse proxy)

## Архитектура Docker

```
┌─────────────────────────────────────────────────┐
│                   Host Machine                  │
│                                                 │
│  ┌───────────────┐          ┌────────────────┐ │
│  │   Frontend    │          │    Backend     │ │
│  │  (Vue + Nginx)│          │   (Go + API)   │ │
│  │               │          │                │ │
│  │  Port: 3000   │◄────────►│ Port: 17777   │ │
│  │               │   HTTP   │ Port: 18888   │ │
│  └───────────────┘          └────────────────┘ │
│         │                           │          │
│         │                           │          │
│         └───────────┬───────────────┘          │
│                     │                          │
│           smart-stream-network                 │
│                                                 │
│  ┌────────────────────┐                        │
│  │   Bind Mounts      │                        │
│  │  ./storage/media   │◄──────────────────────┤
│  │  ./storage/posters │                        │
│  └────────────────────┘                        │
│                                                 │
│                     ▼                           │
│         External PostgreSQL DB                  │
└─────────────────────────────────────────────────┘
```

## Следующие шаги

- [ ] Настройка Nginx reverse proxy с SSL (рекомендуется для production)
- [ ] Настройка автоматического backup для storage и БД
- [ ] Настройка мониторинга и алертов (Prometheus + Grafana)
- [ ] Настройка CI/CD для автоматического деплоя
