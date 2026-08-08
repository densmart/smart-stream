# smart-stream TV Client (Android TV - Kotlin)

This is the Android TV client application for the smart-stream video streaming platform, written in Kotlin.

## Architecture Overview

Android TV application following modern Android architecture:

```
tv_client/
└── app/
    ├── src/main/
    │   ├── kotlin/com/smartstream/tvclient/
    │   │   ├── ui/              # UI layer
    │   │   │   ├── auth/        # Authentication screens
    │   │   │   ├── main/        # Main browsing UI
    │   │   │   ├── details/     # Content details
    │   │   │   ├── player/      # Video player
    │   │   │   └── settings/    # Settings screens
    │   │   ├── data/            # Data layer
    │   │   │   ├── api/         # API client
    │   │   │   ├── model/       # Data models
    │   │   │   └── repository/  # Repositories
    │   │   └── utils/           # Utility classes
    │   ├── res/                 # Resources
    │   │   ├── layout/          # XML layouts
    │   │   ├── drawable/        # Images and drawables
    │   │   ├── values/          # Strings, colors, styles
    │   │   └── xml/             # XML configs
    │   └── AndroidManifest.xml
    ├── build.gradle.kts         # App-level build config
    └── proguard-rules.pro       # ProGuard rules
```

## Key Technologies

- **Language**: Kotlin
- **Platform**: Android TV (Leanback library)
- **Video Player**: ExoPlayer (likely)
- **Networking**: Retrofit + OkHttp
- **Async**: Kotlin Coroutines + Flow
- **DI**: Hilt or Koin (likely)
- **Build**: Gradle (Kotlin DSL)

## Android TV UI Framework

### Leanback Components
The app uses Android TV Leanback library for:
- **BrowseFragment**: Main content browsing
- **DetailsFragment**: Content details and actions
- **PlaybackFragment**: Video playback controls
- **SearchFragment**: Content search (if implemented)

### 10-Foot Experience
- Optimized for viewing from distance
- D-pad navigation support
- Focus management
- Large, readable UI elements
- Simplified interaction model

## Development Workflow

### Building
```bash
cd tv_client
./gradlew assembleDebug
```

### Installing on TV
```bash
./gradlew installDebug
```

### Running Tests
```bash
./gradlew test
```

## Architecture Pattern

Likely follows MVVM or MVI:
- **View**: Activities/Fragments (ui/)
- **ViewModel**: ViewModels with LiveData/StateFlow
- **Model**: Repository pattern (data/)
- **Data Source**: API client + local cache

## Video Playback

### Player Features
- Multiple video format support
- Adaptive streaming
- Playback controls (play, pause, seek)
- Resume functionality
- Error handling and retry
- Buffering states

### Player Lifecycle
- Initialize in Activity/Fragment
- Handle configuration changes
- Release on destroy
- Save/restore playback position

## API Communication

### Authentication
- JWT tokens with extended TTL for TV clients
- Secure token storage (EncryptedSharedPreferences)
- Token refresh mechanism
- Handle authentication failures

### Network Layer
- Retrofit for API calls
- Coroutines for async operations
- Error handling and retry logic
- Response caching when appropriate

### Data Models
- Kotlin data classes
- JSON serialization (Gson or Moshi)
- Mapping between API models and domain models

## Navigation

### D-pad Navigation
- All interactive elements must be focusable
- Proper focus order defined
- Visual feedback for focused elements
- Handle navigation edge cases

### Screen Navigation
- Activity-based or single-Activity with Fragments
- Navigation component (if used)
- Back button handling
- Deep linking (if implemented)

## Performance Considerations

### Memory Management
- Release video player properly
- Cancel coroutines on lifecycle events
- Optimize image loading
- Monitor memory usage

### UI Performance
- Smooth animations (60 fps)
- Lazy loading of content
- Efficient RecyclerView usage
- Background processing for heavy tasks

## Build Configuration

### Gradle
- Kotlin DSL for build scripts
- Dependency management
- Build variants (debug, release)
- ProGuard/R8 for code shrinking

### Signing
- Debug keystore for development
- Release keystore for production
- Secure key management

## Testing Strategy

- Unit tests for ViewModels and repositories
- UI tests for critical user flows
- Test on actual TV devices
- Test different screen sizes

## Important Notes

- Follow Android TV design guidelines
- Optimize for 10-foot viewing distance
- Test with D-pad navigation extensively
- Handle player lifecycle correctly
- Implement proper error handling
- Use coroutines for all async operations
- Follow Kotlin coding conventions
- Keep UI responsive and smooth
- Monitor performance on actual TV hardware
- Secure token storage is critical
