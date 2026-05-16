# Smart Stream TV Client

Android TV приложение для стриминга видео контента с сервера Smart Stream.

## 📋 Описание

Клиентское приложение для Android TV приставок, которое позволяет:
- Авторизоваться в системе
- Просматривать библиотеку медиа-контента
- Просматривать плейлисты
- Воспроизводить видео в режиме стриминга

## 🛠 Технологии

- **Язык:** Kotlin
- **Минимальная версия Android:** API 21 (Android 5.0)
- **Целевая версия:** API 34 (Android 14)
- **UI Framework:** Leanback Library (для TV интерфейса)
- **Видео-плеер:** ExoPlayer 2.19.1
- **HTTP клиент:** Retrofit 2.9.0 + OkHttp 4.11.0
- **Изображения:** Glide 4.16.0
- **Архитектура:** Repository Pattern + ViewModel + LiveData

## 📁 Структура проекта

```
app/src/main/kotlin/com/smartstream/tvclient/
├── SmartStreamApp.kt                 # Application class
├── data/
│   ├── api/
│   │   ├── AuthInterceptor.kt       # JWT authentication interceptor
│   │   ├── ClientApiService.kt      # Retrofit API interface
│   │   └── RetrofitClient.kt        # Retrofit singleton
│   ├── model/
│   │   ├── ApiResponse.kt           # Unified API response wrapper
│   │   ├── Auth.kt                  # Auth request/response models
│   │   ├── Media.kt                 # Media model
│   │   └── Playlist.kt              # Playlist model
│   └── repository/
│       ├── AuthRepository.kt        # Authentication operations
│       ├── MediaRepository.kt       # Media operations
│       └── PlaylistRepository.kt    # Playlist operations
├── ui/
│   ├── auth/
│   │   └── AuthActivity.kt          # Sign-in screen
│   ├── main/
│   │   └── MainActivity.kt          # Main launcher activity
│   └── player/
│       └── PlayerActivity.kt        # Video player screen
└── utils/
    ├── Constants.kt                  # App constants
    └── SharedPrefsManager.kt         # JWT token storage
```

## ⚙️ Настройка проекта

### 1. Клонирование репозитория

Проект уже инициализирован в папке `tv_client/`.

### 2. Открытие в Android Studio

1. Откройте Android Studio
2. File → Open → выберите папку `/path/to/smart-stream/tv_client/`
3. Дождитесь синхронизации Gradle

### 3. Настройка API URL

Настройка сервера выполняется непосредственно в приложении:

1. Запустите приложение на TV
2. На экране авторизации нажмите кнопку настроек (шестеренка) в правом верхнем углу
3. Введите IP адрес и порт вашего сервера
4. Нажмите "Сохранить"

**По умолчанию:** `10.0.2.2:18888` (алиас для localhost на Android эмуляторе)

### 4. Сборка проекта

```bash
# В терминале Android Studio или в корне проекта
./gradlew build
```

### 5. Запуск на эмуляторе TV

1. В Android Studio: Tools → Device Manager
2. Create Device → TV → выберите профиль (например, 1080p)
3. Выберите System Image (API 29+ рекомендуется)
4. Run → Run 'app'

### 6. Установка на реальное устройство

```bash
# Подключите TV приставку по USB или ADB over WiFi
adb connect <TV_IP_ADDRESS>

# Установите APK
./gradlew installDebug

# Или соберите APK вручную
./gradlew assembleDebug
# APK будет в: app/build/outputs/apk/debug/app-debug.apk
```

## 🔐 Авторизация

Для авторизации необходимо:
1. Создать клиента через веб-интерфейс администратора (frontend)
2. Использовать `login` и `password` клиента в TV приложении

**Пример:**
- Login: `client1`
- Password: `password123`

JWT токен автоматически сохраняется в SharedPreferences и добавляется ко всем запросам через `AuthInterceptor`.

## 📡 API Endpoints (ClientAPI)

Приложение использует следующие эндпоинты:

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/auth/sign-in/` | Авторизация клиента |
| GET | `/playlists/` | Список плейлистов |
| GET | `/playlists/:id/media/` | Медиа в плейлисте |
| GET | `/media/` | Медиа без плейлистов |
| GET | `/media/:id/stream/` | Стриминг видео |

Подробнее в `../docs/API.md`.

## 🎨 Ресурсы

- **Темы:** `res/values/themes.xml` - Leanback тема для TV
- **Цвета:** `res/values/colors.xml` - Цветовая палитра
- **Строки:** `res/values/strings.xml` - Локализация
- **Размеры:** `res/values/dimens.xml` - Размеры для TV экранов

## 🚀 Следующие шаги

Согласно DEVELOPMENT_PLAN.md, следующие этапы:

### ✅ Выполнено:
- [x] 3.1 Выбор технологии
- [x] 3.2 Инициализация проекта
- [x] 3.3.1 UI для TV - Base UI (экраны авторизации, главный экран, детали медиа/плейлистов)
- [x] 3.3.2 UI для TV - Настройки сервера (диалог конфигурации IP и порта)
- [x] 3.4 Авторизация - Реализована логика авторизации с валидацией
- [x] 3.5 Отображение списков - Реализованы списки медиа и плейлистов с Leanback
- [x] 3.6 Видео-плеер - ExoPlayer с HTTP streaming и JWT аутентификацией

### 📋 TODO:
- [ ] 3.7 Тестирование - Тестирование на реальном устройстве

## 📝 Заметки для разработчика

### Важные файлы для редактирования:

1. **Constants.kt** - обновите BASE_URL перед запуском
2. **AuthActivity.kt** - реализуйте UI и логику авторизации
3. **MainActivity.kt** - добавьте BrowseFragment для отображения контента
4. **PlayerActivity.kt** - интегрируйте ExoPlayer для воспроизведения

### Зависимости:

Все необходимые зависимости уже добавлены в `app/build.gradle.kts`:
- Leanback Library - для TV UI
- ExoPlayer - для видео
- Retrofit + OkHttp - для API
- Glide - для загрузки изображений
- Coroutines - для асинхронности

### Debugging:

HTTP логирование включено в `RetrofitClient.kt`:
```kotlin
interceptor.level = HttpLoggingInterceptor.Level.BODY
```

Для production режима измените на `Level.NONE`.

## 📄 Лицензия

Проект создан исключительно для образовательных целей.
