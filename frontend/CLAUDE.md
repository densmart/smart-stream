# smart-stream Frontend (Vue.js Admin Panel)

This is the admin panel for the smart-stream video streaming platform, built with Vue 3 and PrimeVue.

## Architecture Overview

Modern Vue 3 application with Composition API:

```
frontend/
├── src/
│   ├── components/       # Reusable UI components
│   ├── views/           # Page-level components
│   ├── layout/          # Layout components (AppLayout, AppTopbar, etc.)
│   ├── router/          # Vue Router configuration
│   ├── api/             # API service layer
│   ├── service/         # Business logic services
│   ├── types/           # TypeScript definitions
│   └── assets/          # Static assets (styles, images)
├── public/              # Public static files
└── index.html           # Entry HTML file
```

## Key Technologies

- **Framework**: Vue 3.4+ (Composition API)
- **UI Library**: PrimeVue 4.5+
- **Theming**: @primeuix/themes, Tailwind CSS 4.1+
- **State Management**: Pinia 3.0+
- **Routing**: Vue Router 4.4+
- **HTTP Client**: Axios 1.16+
- **Build Tool**: Vite 5.3+
- **Icons**: PrimeIcons 7.0+

## Project Structure

### Components
Reusable UI components that compose the application. Leverage PrimeVue components extensively.

### Views
Page-level components that correspond to routes. These orchestrate components and handle page-specific logic.

### Layout
Application shell components (header, sidebar, footer) that provide consistent structure.

### API Layer
Centralized API communication using Axios. All backend calls should go through this layer.

### Services
Business logic that doesn't fit into components or the API layer.

### State Management
Pinia stores for global application state. Organize by feature domain.

## API Documentation

**IMPORTANT**: When you need information about backend API endpoints, parameters, or response formats:

1. ✅ **ALWAYS check** `/Users/denyskazka/work/own/smart-stream/docs/WEB_API.md` first
2. ❌ **NEVER read** backend source code to understand API
3. ❌ **DO NOT guess** API parameters or formats

The `docs/WEB_API.md` file contains complete API documentation including:
- All endpoints with methods (GET, POST, PATCH, DELETE)
- Query parameters and their types
- Request/response formats
- Pagination parameters
- Search/filter parameters
- Error responses

If the documentation is incomplete or unclear, ask the user to update it via the backend expert, but do not read backend code yourself.

## Development Workflow

### Running Locally
```bash
cd frontend
npm install
npm run dev
```

### Building for Production
```bash
npm run build
```

### Linting
```bash
npm run lint
```

## Styling Approach

- **Primary**: Tailwind utility classes
- **Components**: PrimeVue themed components
- **Custom styles**: SASS when needed
- **Responsive**: Mobile-first approach

## API Integration

### Authentication
- JWT tokens for authentication
- Axios interceptors attach tokens to requests
- Automatic token refresh (if implemented)
- Redirect to login on 401 errors

### Error Handling
- Global error interceptor in Axios
- User-friendly error messages via PrimeVue Toast
- Consistent error handling across all API calls

### Loading States
- Show loading indicators during async operations
- Use PrimeVue ProgressSpinner or skeleton screens

## State Management with Pinia

Stores are organized by domain:
- Auth store: User authentication and permissions
- Content store: Video content management
- User store: User management
- etc.

## Routing

- Vue Router 4 with history mode
- Route guards for authentication
- Lazy-loaded route components
- Nested routes for complex layouts

## PrimeVue Integration

- Auto-import configured via unplugin-vue-components
- Theme customization in theme configuration
- Extensive use of:
  - DataTable for listings
  - Dialog for modals
  - Toast for notifications
  - Form components with validation
  - Menu and navigation components

## Code Style

- Use `<script setup>` syntax
- Composition API preferred
- TypeScript for type safety
- Composables for reusable logic
- Props validation
- Emits declaration

## Important Notes

- Keep components focused and single-responsibility
- Use PrimeVue components for consistency
- Centralize API calls in the api/ directory
- Handle all async operations with proper error handling
- Ensure responsive design for all screens
- Validate forms before submission
- Provide clear user feedback for all actions
- Follow Vue 3 and PrimeVue best practices
