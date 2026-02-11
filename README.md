# Refurbed

A full-stack product catalog application with filtering capabilities, built with Vue.js frontend and Go backend.

## Overview

This project consists of two main components:

- **Frontend**: Vue.js application for browsing and filtering products
- **Backend**: Go API serving product data with filtering and caching

## Quick Start

### Backend

**Prerequisites:** Go 1.22+

```bash
cd backend
make start
```

The backend server runs on `http://localhost:8080`.

### Frontend

**Prerequisites:** Node.js 18+

```bash
cd frontend
npm run dev
```

The frontend server runs on `http://localhost:5173`.

## Main Features

- Product catalog with search functionality
- Filtering by color, price range, and bestseller status
- In-memory caching for improved performance
- Responsive UI built with Tailwind CSS

## API Endpoints

- `GET /products` - Product list with optional query filters
- `GET /health` - Health check

## Documentation

For detailed documentation, architecture details, production readiness improvements, and more information, please refer to:

- **[Backend README](backend/README.md)** - Backend architecture, API details, and production improvements
- **[Frontend README](frontend/README.md)** - Frontend architecture, components, and production improvements
