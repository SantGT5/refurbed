export const BASE_URL = "http://localhost:8080";

/**
 * @param {string} path - Path (e.g. '/users') or full URL
 * @returns {string} Full URL
 */
export function getFullUrl(path) {
  if (path.startsWith("http://") || path.startsWith("https://")) {
    return path;
  }

  const base = BASE_URL.replace(/\/$/, "");
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;

  return `${base}${normalizedPath}`;
}
