# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Zashboard is a Wails desktop application that packages a Clash API-based web dashboard into a lightweight native desktop app. It supports Windows, macOS, and Linux.

**Tech Stack:**
- Backend: Go 1.25 + Wails v2
- Frontend: Vue 3 + TypeScript + Vite + Tailwind CSS + DaisyUI

## Development Commands

### Main Project (Wails)
```bash
wails dev          # Development mode with hot reload
wails build        # Production build
```

### Frontend (in /frontend directory)
```bash
pnpm install       # Install dependencies
pnpm run dev       # Vite dev server
pnpm run build     # Production build (full fonts)
pnpm run build:misans-only   # Build with MiSans font only (default)
pnpm run lint      # ESLint check and auto-fix
pnpm run format    # Prettier formatting
pnpm run type-check   # TypeScript type checking
```

Font build variants: `build:cdn-fonts`, `build:firasans-only`, `build:pingfang-only`, `build:sarasa-only`, `build:no-fonts`

## Architecture

### Directory Structure
```
Zashboard/
├── app.go, main.go       # Wails application entry
├── wails.json            # Wails config (frontend commands)
├── frontend/             # Vue 3 SPA (git submodule)
│   ├── src/
│   │   ├── api/          # Clash API integration
│   │   ├── store/        # Vue reactive state (ref/reactive, no Pinia)
│   │   ├── views/        # Page components
│   │   ├── components/   # Reusable components (organized by feature)
│   │   ├── composables/  # Vue composition hooks
│   │   ├── i18n/         # Internationalization
│   │   ├── helper/       # Utility functions
│   │   ├── types/        # TypeScript definitions
│   │   └── constant/     # Constants (EMOJIS, FONTS, ROUTE_NAME)
│   └── vite.config.ts    # Vite config with PWA and Git commit injection
└── build/                # Platform-specific build resources
```

### Key Patterns

**State Management:** Uses Vue 3 `ref`/`reactive` directly in `src/store/*.ts` files - no external state library.

**Routing:** Vue Router with hash mode. Main routes under `/` include: proxies, overview, connections, logs, rules, settings. Setup page at `/setup`.

**Components:** Single-file components using `<script setup>` syntax with TypeScript.

**API Layer:** `src/api/index.ts` handles Clash API calls. WebSocket connections use `reconnectingwebsocket` for real-time updates.

**Theming:** DaisyUI themes with custom background blur effects. Supports multiple font families and emoji sets (Twemoji/Noto).

## Code Style

- **Prettier:** No semicolons, single quotes, 100 char width, one attribute per line
- **ESLint:** Vue recommended rules with TypeScript
- **Git Hooks:** Husky runs lint-staged on pre-commit (ESLint + Prettier + sort-package-json)

## Important Files

| File | Purpose |
|------|---------|
| `frontend/src/store/settings.ts` | User preferences state |
| `frontend/src/store/proxies.ts` | Proxy selector logic |
| `frontend/src/api/index.ts` | Clash API integration |
| `frontend/src/App.vue` | Global theme/font/background management |
| `frontend/src/router/index.ts` | Route definitions and transitions |