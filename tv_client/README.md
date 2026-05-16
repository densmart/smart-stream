# Smart Stream TV Client

Android TV приложение для стриминга видео контента с сервера Smart Stream.

## 📋 Описание

Полнофункциональное клиентское приложение для Android TV приставок с кастомным UI, которое предоставляет:

### 🎯 Основные возможности:
- **Авторизация клиентов** с JWT токенами
- **Настройка сервера** - изменение IP и порта прямо из приложения
- **Просмотр контента** - медиафайлы и плейлисты в grid-интерфейсе
- **Иерархия плейлистов** - поддержка вложенных плейлистов (franchise → series → episodes)
- **Стриминг видео** - ExoPlayer с HTTP streaming и JWT аутентификацией
- **Сохранение прогресса** - автоматическое сохранение позиции воспроизведения
- **Превью контента** - боковая панель с деталями при фокусе
- **Автоматический логаут** - при истечении токена (401 response)

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
├── SmartStreamApp.kt                          # Application class
├── data/
│   ├── api/
│   │   ├── AuthInterceptor.kt                # JWT authentication interceptor
│   │   ├── UnauthorizedInterceptor.kt        # Auto-logout on 401 response
│   │   ├── ClientApiService.kt               # Retrofit API interface
│   │   └── RetrofitClient.kt                 # Retrofit singleton with dynamic base URL
│   ├── model/
│   │   ├── ApiResponse.kt                    # Unified API response wrapper
│   │   ├── Auth.kt                           # Auth request/response models
│   │   ├── Media.kt                          # Media model with helpers (poster, stream URL, formatting)
│   │   └── Playlist.kt                       # Playlist model with hierarchy support
│   └── repository/
│       ├── AuthRepository.kt                 # Authentication operations
│       ├── MediaRepository.kt                # Media operations (unassigned + playlist media)
│       └── PlaylistRepository.kt             # Playlist operations (list + children)
├── ui/
│   ├── auth/
│   │   ├── AuthActivity.kt                   # Sign-in screen with TV-optimized UI
│   │   └── ServerConfigDialog.kt             # Server settings dialog (deprecated)
│   ├── main/
│   │   ├── MainActivity.kt                   # Main launcher with auth check
│   │   ├── MainBrowseFragment.kt             # Original Leanback browse (deprecated)
│   │   ├── MainBrowseFragmentNew.kt          # Custom grid UI with tabs and preview
│   │   ├── CardPresenter.kt                  # Leanback card presenter
│   │   └── MediaCardAdapter.kt               # RecyclerView adapter for grid
│   ├── details/
│   │   ├── PlaylistDetailsActivity.kt        # Playlist details container
│   │   ├── PlaylistDetailsFragment.kt        # Playlist details with child playlists/media
│   │   ├── MediaDetailsActivity.kt           # Media details container
│   │   ├── MediaDetailsFragment.kt           # Media details with play action (deprecated)
│   │   ├── MediaDetailsFragmentSimple.kt     # Simplified media details
│   │   └── MediaDetailsDescriptionPresenter.kt # Leanback description presenter
│   ├── player/
│   │   └── PlayerActivity.kt                 # ExoPlayer-based video player with position saving
│   └── settings/
│       └── ServerSettingsDialogFragment.kt   # Server configuration dialog
└── utils/
    ├── Constants.kt                          # App constants (defaults, keys)
    └── SharedPrefsManager.kt                 # Storage for JWT, server config, playback positions
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
| POST | `/auth/sign-in/` | Авторизация клиента (login + password) |
| GET | `/playlists/` | Список плейлистов верхнего уровня |
| GET | `/playlists/:id/children/` | Дочерние плейлисты (иерархия) |
| GET | `/playlists/:id/media/` | Медиа в плейлисте |
| GET | `/media/` | Медиа без плейлистов (unassigned) |
| GET | `/media/:id/stream/` | Стриминг видео (HTTP Range requests) |
| GET | `/static/posters/:filename` | Постеры для медиа и плейлистов |

**Все эндпоинты (кроме `/auth/sign-in/`) требуют JWT токен в заголовке:**
```
Authorization: Bearer <token>
```

Подробнее в `../docs/API.md`.

## 💡 Примеры использования API

### Получение плейлистов:
```kotlin
// PlaylistRepository.kt
val result = playlistRepository.getPlaylists(limit = 100, offset = 0)
result.onSuccess { playlists ->
    // playlists: List<Playlist>
    playlists.forEach { playlist ->
        println("${playlist.name} - hasChildren: ${playlist.hasChildren}")
    }
}
```

### Получение дочерних плейлистов:
```kotlin
val result = playlistRepository.getPlaylistChildren(playlistId, limit = 100, offset = 0)
result.onSuccess { childPlaylists ->
    // childPlaylists: List<Playlist>
}
```

### Получение медиа в плейлисте:
```kotlin
val result = mediaRepository.getPlaylistMedia(playlistId, limit = 100, offset = 0)
result.onSuccess { apiResponse ->
    if (apiResponse.isSuccess()) {
        val mediaList = apiResponse.result // List<Media>
    }
}
```

### Получение unassigned медиа:
```kotlin
val result = mediaRepository.getMediaList(onlyUnassigned = true)
result.onSuccess { mediaList ->
    // mediaList: List<Media>
}
```

### Построение stream URL:
```kotlin
val media: Media = ...
val streamUrl = media.getStreamUrl(SharedPrefsManager.getBaseUrl())
// Результат: "http://192.168.0.17:18888/media/uuid/stream/"
```

## 🎬 User Flow (Как работает приложение)

### 1️⃣ Запуск и авторизация
```
MainActivity → проверка токена
  ├─ Если токен валиден → MainBrowseFragmentNew
  └─ Если нет токена → AuthActivity
      └─ Ввод login/password → сохранение JWT → MainActivity
```

### 2️⃣ Главный экран (MainBrowseFragmentNew)
```
Табы: [Media] [Playlists]  [Settings] [Search]

Grid (6 колонок):          Preview Panel (появляется при фокусе):
┌─────┬─────┬─────┐        ┌──────────────┐
│ [1] │ [2] │ [3] │        │   Poster     │
├─────┼─────┼─────┤   →    │   Title      │
│ [4] │ [5] │ [6] │        │   Duration   │
└─────┴─────┴─────┘        │   Size       │
                            └──────────────┘
```

**Media Tab:** Показывает unassigned медиа (не привязанные к плейлистам)
**Playlists Tab:** Показывает плейлисты верхнего уровня (parent_id = null)

### 3️⃣ Навигация по плейлистам
```
MainBrowseFragment → Клик на Playlist
  └─ PlaylistDetailsActivity/Fragment
      ├─ Если hasChildren = true:
      │   └─ Показывает дочерние плейлисты (GET /playlists/:id/children/)
      │       └─ Клик на child playlist → рекурсивно открывает PlaylistDetailsActivity
      └─ Если hasChildren = false:
          └─ Показывает медиа в плейлисте (GET /playlists/:id/media/)
              └─ Клик на Media → MediaDetailsActivity
```

### 4️⃣ Воспроизведение медиа
```
MediaDetailsActivity
  ├─ Показывает детали медиа (постер, название, длительность, размер)
  ├─ Проверяет сохраненную позицию (SharedPrefsManager)
  └─ Кнопка "Play" → PlayerActivity
      ├─ ExoPlayer с JWT токеном
      ├─ Стриминг с {BASE_URL}/media/{id}/stream/
      ├─ Сохранение позиции в onPause()
      └─ Очистка позиции при завершении (STATE_ENDED)
```

### 5️⃣ Автоматический логаут
```
Любой API запрос → 401 Unauthorized
  └─ UnauthorizedInterceptor
      ├─ Очищает JWT токен
      ├─ Очищает auth данные
      └─ Редирект на AuthActivity
```

## 🎨 Ресурсы

- **Темы:** `res/values/themes.xml` - Leanback тема для TV
- **Цвета:** `res/values/colors.xml` - Цветовая палитра
- **Строки:** `res/values/strings.xml` - Локализация
- **Размеры:** `res/values/dimens.xml` - Размеры для TV экранов
- **Layouts:** `res/layout/` - XML layouts для всех экранов
- **Drawable:** `res/drawable/` - Иконки и placeholder изображения

## ✅ Текущий статус проекта

Все основные функции полностью реализованы и работают:

### ✅ Выполнено:
- [x] **Авторизация** - JWT-based аутентификация с валидацией
- [x] **Настройка сервера** - Динамическая настройка IP и порта из UI
- [x] **Главный экран** - Кастомный grid UI с табами (Media/Playlists)
- [x] **Иерархия плейлистов** - Поддержка parent-child структуры (franchise → series → episodes)
- [x] **Детали плейлистов** - Отображение дочерних плейлистов или медиа внутри
- [x] **Детали медиа** - Информация о медиа с кнопкой воспроизведения
- [x] **Видео-плеер** - ExoPlayer с HTTP streaming и JWT аутентификацией
- [x] **Сохранение позиции** - Автоматическое сохранение/восстановление позиции воспроизведения
- [x] **Preview панель** - Боковая панель с деталями при фокусе на карточке
- [x] **Автологаут** - UnauthorizedInterceptor для обработки 401 responses
- [x] **Постеры** - Загрузка и отображение постеров через Glide
- [x] **TV навигация** - Оптимизированное управление D-pad'ом

### 🎯 Эталонная версия (v1.0)
Текущая версия считается стабильной и полностью функциональной для production use.

## 📝 Заметки для разработчика

### 🔑 Ключевые особенности реализации:

#### 1. **Архитектура UI**
- **MainBrowseFragmentNew** - основной экран с кастомным grid layout (6 колонок)
- **Динамическая смена колонок** - 6 без preview, 4 с preview панелью
- **Табы** - Media/Playlists с независимой загрузкой данных
- **Custom GridLayoutManager** - с оптимизированной навигацией D-pad (wrap-around navigation)

#### 2. **Иерархия плейлистов**
- **PlaylistDetailsFragment** проверяет флаг `hasChildren`
- Если `true` → загружает дочерние плейлисты через `/playlists/:id/children/`
- Если `false` → загружает медиа через `/playlists/:id/media/`
- **Рекурсивная навигация** - можно заходить в плейлисты любой глубины

#### 3. **Сохранение позиции воспроизведения**
```kotlin
// PlayerActivity.kt:onPause()
SharedPrefsManager.saveMediaPosition(mediaId, currentPosition)

// MediaDetailsActivity - восстановление
val savedPosition = SharedPrefsManager.getMediaPosition(mediaId)
intent.putExtra(EXTRA_START_POSITION, savedPosition)
```
- Позиция сохраняется в onPause() плеера
- Автоматически восстанавливается при следующем запуске
- Очищается при завершении видео (STATE_ENDED)

#### 4. **JWT аутентификация**
- **AuthInterceptor** - добавляет `Authorization: Bearer <token>` ко всем запросам
- **UnauthorizedInterceptor** - отлавливает 401, очищает токен, редиректит на логин
- **RetrofitClient** - динамически использует URL из SharedPrefsManager

#### 5. **ExoPlayer streaming**
```kotlin
// Создание DataSource с JWT токеном
val httpDataSourceFactory = DefaultHttpDataSource.Factory()
    .setDefaultRequestProperties(mapOf("Authorization" to "Bearer $jwtToken"))

// Поддержка HTTP Range requests для seeking
val dataSourceFactory = DefaultDataSource.Factory(this, httpDataSourceFactory)
player = ExoPlayer.Builder(this)
    .setMediaSourceFactory(DefaultMediaSourceFactory(dataSourceFactory))
    .build()
```

#### 6. **Модели данных**
- **Media.kt** - helper методы:
  - `getPosterUrl(baseUrl)` - построение URL постера
  - `getStreamUrl(baseUrl)` - построение stream URL
  - `getFormattedDuration()` - форматирование времени (HH:MM:SS)
  - `getFormattedSize()` - форматирование размера (KB/MB/GB)
- **Playlist.kt** - helper методы:
  - `getPosterUrl(baseUrl)` - построение URL постера
  - `getTypeDisplayName()` - человеко-читаемый тип

#### 7. **Настройки сервера**
- **ServerSettingsDialogFragment** - диалог для изменения IP и порта
- **SharedPrefsManager** хранит host и port отдельно
- **RetrofitClient пересоздается** при изменении настроек (требует перезапуск приложения)

### 🛠 Важные файлы:

| Файл | Назначение |
|------|-----------|
| `MainBrowseFragmentNew.kt` | Главный экран с grid и preview |
| `PlaylistDetailsFragment.kt` | Детали плейлиста с иерархией |
| `PlayerActivity.kt` | Видео-плеер с сохранением позиции |
| `SharedPrefsManager.kt` | Хранение JWT, настроек, позиций |
| `RetrofitClient.kt` | Динамический base URL |
| `ClientApiService.kt` | Все API endpoints |
| `Media.kt` / `Playlist.kt` | Модели с helper методами |

### 🐛 Debugging:

**HTTP логирование включено в `RetrofitClient.kt`:**
```kotlin
interceptor.level = HttpLoggingInterceptor.Level.BODY
```
Для production: `Level.NONE`

**Все ключевые действия логируются:**
- `PlayerActivity` - подробные логи плеера
- `PlaylistDetailsFragment` - загрузка плейлистов/медиа
- `MainBrowseFragment` - клики и навигация

**Просмотр логов:**
```bash
adb logcat | grep "SmartStream"
# или конкретные теги:
adb logcat | grep "PlayerActivity"
adb logcat | grep "PlaylistDetailsFrag"
```

## 🔮 Потенциальные улучшения

Текущая версия полностью функциональна, но можно добавить:

### UI/UX:
- [ ] **Поиск** - реализация функции поиска по названию медиа/плейлистов
- [ ] **Фильтрация** - фильтры по типу, дате, формату
- [ ] **Сортировка** - изменение порядка отображения (по дате, названию, длительности)
- [ ] **Настройки воспроизведения** - качество, субтитры (если поддерживается сервером)
- [ ] **История просмотров** - список последних просмотренных медиа
- [ ] **Избранное** - добавление медиа в избранное на клиенте

### Функциональность:
- [ ] **Offline режим** - кеширование для просмотра без сети
- [ ] **Background playback** - воспроизведение в фоне (Picture-in-Picture)
- [ ] **Плейлисты воспроизведения** - автоматический переход к следующему эпизоду
- [ ] **Subtitles** - поддержка субтитров (если есть на сервере)
- [ ] **Quality selection** - выбор качества стрима
- [ ] **Родительский контроль** - возрастные рейтинги

### Производительность:
- [ ] **Pagination** - ленивая подгрузка при скролле больших списков
- [ ] **Image caching** - улучшенное кеширование постеров
- [ ] **Preloading** - предзагрузка следующего эпизода

### Технические:
- [ ] **Unit тесты** - покрытие репозиториев и ViewModels
- [ ] **UI тесты** - Espresso тесты для критических флоу
- [ ] **CI/CD** - автоматическая сборка и тестирование
- [ ] **Crashlytics** - сбор крэш-репортов
- [ ] **Analytics** - аналитика использования

## ⚠️ Известные ограничения

1. **Перезапуск при смене сервера** - требуется перезапуск приложения после изменения настроек сервера (RetrofitClient инициализируется один раз)
2. **Без HLS/DASH** - используется простой HTTP streaming (не адаптивный bitrate)
3. **Без субтитров** - пока не реализована поддержка субтитров
4. **HTTP only** - нет поддержки HTTPS (можно добавить при необходимости)
5. **Portrait layout** - UI оптимизирован только для landscape (TV)
