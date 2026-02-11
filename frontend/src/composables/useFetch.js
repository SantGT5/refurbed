import { isRef, onMounted, ref, watch } from "vue";
import { getFullUrl } from "../api/client";

/**
 * Generic fetch composable with loading, callbacks and optional fetch on render.
 *
 * @param {Object} options
 * @param {string} options.url - Path (e.g. '/users') or full URL; base is http://localhost:8080 for paths
 * @param {string} [options.method='GET'] - HTTP method (GET, POST, PUT, PATCH, DELETE, etc.)
 * @param {Object|string|null} [options.body=null] - Request body (object will be JSON.stringify'd)
 * @param {Object} [options.headers={}] - Additional headers (Content-Type: application/json added when body is object)
 * @param {boolean} [options.fetchOnRender=true] - If true, runs fetch when the component is mounted
 * @param {Array<any>} [options.refetchDeps=[]] - When any of these values change (e.g. refs or reactive deps), refetch is triggered
 * @param {function} [options.onSuccess] - Called with (data, response) on success
 * @param {function} [options.onError] - Called with (error, response) on error
 * @returns {{ data: Ref, error: Ref, loading: Ref, execute: (overrides?: Object) => Promise<void> }}
 */
export function useFetch(options = {}) {
  const {
    url = "",
    method = "GET",
    body = null,
    headers: customHeaders = {},
    fetchOnRender = false,
    refetchDeps = [],
    onSuccess,
    onError,
  } = options;

  const data = ref(null);
  const error = ref(null);
  const loading = ref(false);

  async function execute(overrides = {}) {
    const finalUrl = overrides.url ?? url;
    const finalMethod = (overrides.method ?? method).toUpperCase();
    const finalBody = overrides.body !== undefined ? overrides.body : body;
    const finalHeaders = { ...customHeaders, ...overrides.headers };

    if (!finalUrl) {
      error.value = new Error("useFetch: url is required");
      if (onError) onError(error.value, null);
      return;
    }

    loading.value = true;
    error.value = null;
    data.value = null;

    const fullUrl = getFullUrl(finalUrl);
    const requestHeaders = { ...finalHeaders };

    let requestBody = finalBody;
    if (finalBody != null && finalMethod !== "GET") {
      if (typeof finalBody === "object" && !(finalBody instanceof FormData)) {
        requestBody = JSON.stringify(finalBody);
        if (!requestHeaders["Content-Type"]) {
          requestHeaders["Content-Type"] = "application/json";
        }
      }
    }

    let response = null;
    try {
      response = await fetch(fullUrl, {
        method: finalMethod,
        headers: Object.keys(requestHeaders).length
          ? requestHeaders
          : undefined,
        body: finalMethod === "GET" ? undefined : requestBody,
      });

      const contentType = response.headers.get("content-type");
      const isJson = contentType && contentType.includes("application/json");
      const responseData = isJson
        ? await response.json()
        : await response.text();

      if (!response.ok) {
        const err = new Error(response.statusText || `HTTP ${response.status}`);
        err.status = response.status;
        err.data = JSON.parse(responseData);
        error.value = err;
        if (onError) onError(err, response);
        return;
      }

      data.value = responseData;
      if (onSuccess) onSuccess(responseData, response);
    } catch (err) {
      error.value = err;
      if (onError) onError(err, response);
    } finally {
      loading.value = false;
    }
  }

  if (fetchOnRender && url) {
    onMounted(() => execute());
  }

  if (Array.isArray(refetchDeps) && refetchDeps.length > 0 && url) {
    watch(
      () => refetchDeps.map((d) => (isRef(d) ? d.value : d)),
      () => execute(),
      { deep: true },
    );
  }

  return {
    data,
    error,
    loading,
    execute,
  };
}
