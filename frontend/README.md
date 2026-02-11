# Refurbed Frontend

## Prerequisites

**Prerequisites:** Node.js 18+

## Quick Start

```bash
npm run dev
```

The server listens on `http://localhost:5173`.

---

## What we would improve/change for a production-ready environment

### 1. Environment & Configuration

- **Environment variables** – Move `BASE_URL` from hardcoded `http://localhost:8080` in `src/api/client.js` to env vars (e.g. `VITE_API_BASE_URL`), with `.env`, `.env.development`, and `.env.production`.
- **API configuration** – Support different API origins per environment and avoid CORS issues by using proxy configuration in Vite for local dev.

### 2. Error Handling & Observability

- **Centralized error handling** – Standardize API error handling (e.g. toast notifications, retry logic, user-facing messages).
- **Logging & monitoring** – Integrate an error-tracking service (e.g. Sentry) and remove `console.log` from production code.
- **useFetch improvements** – In `useFetch.js`, avoid `JSON.parse(responseData)` on non-JSON responses; add request cancellation (e.g. AbortController) for unmounting/navigation.

### 3. Security

- **Input sanitization** – Ensure search and filter inputs are validated/sanitized before sending to the API or using in the UI.
- **Content Security Policy** – Add CSP headers in production build/deploy.
- **Dependency audit** – Run `npm audit` regularly and keep dependencies up to date.

### 4. Performance

- **Code splitting** – Use lazy loading with `defineAsyncComponent` and Vue Router for route-based code splitting.
- **Asset optimization** – Configure Vite for chunk splitting, tree shaking, and compression (e.g. gzip/brotli).
- **Caching** – Add HTTP caching headers for static assets; consider service worker / PWA for offline support where appropriate.

### 5. Testing

- **Unit tests** – Add Vitest (or Jest) for composables (`useFetch`), utils (`url.js`), and reusable components.
- **E2E tests** – Add Playwright or Cypress for critical flows (search, filters, product listing).
- **CI** – Run lint, type-check, and tests on every PR.

### 6. Code Quality

- **TypeScript** – Migrate to TypeScript for better type safety, API contracts, and editor support.
- **Linting** – Add ESLint with Vue and accessibility plugins (e.g. `eslint-plugin-vue`, `eslint-plugin-jsx-a11y`).
- **Formatting** – Add Prettier and enforce consistent style.
- **Husky + lint-staged** – Run lint and tests before commits.

### 7. Routing & State

- **Vue Router** – Introduce routing for multiple pages (e.g. product detail, filters in URL, deep links).

### 8. UX & Accessibility

- **Loading states** – Ensure consistent loading/skeleton UIs for async content.
- **Accessibility** – Add ARIA labels, keyboard support, focus management, and semantic HTML.
- **SEO** – For public pages, consider SSR/SSG (e.g. Nuxt) or pre-rendering if SEO matters.

### 9. DevOps & Deployment

- **Docker** – Provide a Dockerfile for reproducible production builds.
- **CI/CD** – Set up pipeline for build, test, and deploy (e.g. GitHub Actions).
- **Health checks** – Add a minimal health endpoint or static check for deployment validation.

### 10. Documentation & Maintenance

- **API documentation** – Document API endpoints and expected payloads (e.g. OpenAPI/Swagger).
- **Component documentation** – Use Storybook for component catalog and visual regression tests.

---

## Notes on architecture, decisions & other comments

### High-level architecture

The app follows a **single-page structure** with a flat folder layout:

- `api/` – HTTP client and base URL configuration
- `components/` – Reusable UI components (barrel export via `index.js`)
- `composables/` – Shared logic (e.g. `useFetch`) for data fetching
- `utils/` – Pure helpers (URL building, query parsing)

There is no router or global state store; the main view and form state live in `App.vue`.

### Data fetching

- **useFetch composable** – Generic data-fetching composable built on the native `fetch` API. It exposes `data`, `error`, `loading`, and an `execute` function. It supports `fetchOnRender`, `refetchDeps`, and callbacks (`onSuccess`, `onError`).
- **Manual execution** – In `App.vue`, `useFetch` is used without `fetchOnRender`; the `execute` function is called from a `watch` on the form state, so the fetch is driven by filter changes rather than mount.
- **API client** – `getFullUrl()` in `api/client.js` prepends a configurable base URL to relative paths and passes through absolute URLs unchanged.

### URL state & filters

- **URL as source of truth** – Filter values are read from the query string on load via `getQueryParams()`, and updates are written back with `updateUrlWithQueries()` using `history.replaceState` to avoid cluttering history.
- **Debouncing** – A 500ms debounce is applied to filter changes before updating the URL and refetching, reducing API calls while typing in search.

### UI & styling

- **Tailwind CSS** – Utility-first styling and responsive layout (e.g. grid breakpoints).
- **PrimeIcons** – Icon set for UI elements.
- **Component design** – Reusable components (`Input`, `DropDown`, `Range`, `ProductCard`, etc.) with props like `variant`, `size`, and slots for flexibility.

### Notable decisions & trade-offs

| Decision                               | Rationale                                                                    |
| -------------------------------------- | ---------------------------------------------------------------------------- |
| Native `fetch` over Axios              | Simpler, no extra HTTP client                                                |
| No Vue Router                          | Single page; filters are expressed in the URL instead of routes              |
| No Pinia/Vuex                          | Form state is local to `App.vue`; no cross-component shared state needed yet |
| Barrel exports (`components/index.js`) | Tidier imports and a single entry point for shared components                |
| Centralized URL utilities              | Query parsing and building are reused and easier to test                     |

### Other comments

- **Product data shape** – The product list handling checks `Object.values(products || {})` for arrays; the API may return an object keyed by id or a plain array, so the template logic is written to support both.
- **Filter options** – Color and bestseller options are currently defined inline in `App.vue`; they could be moved to constants or fetched from the API for easier maintenance.
