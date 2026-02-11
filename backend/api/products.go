package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"assignment-backend/cache"
)

var productsCache = cache.NewProductsCache(30 * time.Second)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Something went wrong"}`, http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	cacheKey := cache.BuildCacheKey(q)

	// Try to get from cache first (cached by query parameters)
	if cached, ok := productsCache.Get(cacheKey); ok {
		w.WriteHeader(http.StatusOK)
		w.Write(cached)
		return
	}

	// Cache miss - read from file
	productsBytes, err := os.ReadFile("data/products.json")
	if err != nil {
		http.Error(w, `{"error":"failed to read products"}`, http.StatusInternalServerError)
		return
	}

	var productsMap map[string]map[string]interface{}
	if err = json.Unmarshal(productsBytes, &productsMap); err != nil {
		http.Error(w, `{"error":"invalid products json"}`, http.StatusInternalServerError)
		return
	}

	// search
	search := strings.TrimSpace(q.Get("search"))
	searchLower := strings.ToLower(search)

	// color
	color := strings.TrimSpace(q.Get("color"))

	// bestseller
	bestsellerStr, hasBestseller := q["bestseller"]
	var bestseller bool
	if hasBestseller {
		bestseller = (len(bestsellerStr) > 0 && strings.ToLower(bestsellerStr[0]) == "true")
	}

	// minPrice
	minPriceStr, hasMinPrice := q["minPrice"]
	var minPrice float64
	if hasMinPrice {
		minPrice, _ = strconv.ParseFloat(minPriceStr[0], 64)
	}

	// maxPrice
	maxPriceStr, hasMaxPrice := q["maxPrice"]
	var maxPrice float64
	if hasMaxPrice {
		maxPrice, _ = strconv.ParseFloat(maxPriceStr[0], 64)
	}

	matches := func(p map[string]interface{}) bool {
		// SEARCH filter
		if search != "" {
			name, _ := p["name"].(string)
			if !strings.Contains(strings.ToLower(name), searchLower) {
				return false
			}
		}

		// COLOR filter
		if color != "" {
			rawColors, ok := p["colors"].([]interface{})
			if !ok {
				return false
			}
			found := false

			for _, c := range rawColors {
				cs, _ := c.(string)
				if cs == color {
					found = true
					break
				}
			}

			if !found {
				return false
			}
		}

		// BESTSELLER filter
		if hasBestseller {
			bs, ok := p["bestseller"].(bool)
			if !ok || bs != bestseller {
				return false
			}
		}

		// MIN_PRICE filter
		if hasMinPrice {
			mp, ok := p["base_price"].(float64)
			if !ok || mp < minPrice {
				return false
			}
		}

		// MAX_PRICE filter
		if hasMaxPrice {
			mp, ok := p["base_price"].(float64)
			if !ok || mp > maxPrice {
				return false
			}
		}

		return true
	}

	filtered := make(map[string]map[string]interface{})
	for id, product := range productsMap {
		if matches(product) {
			filtered[id] = product
		}
	}

	var out []byte
	if len(q) == 0 {
		out, err = json.Marshal(productsMap)
	} else {
		out, err = json.Marshal(filtered)
	}

	if err != nil {
		http.Error(w, `{"error":"failed to serialize response"}`, http.StatusInternalServerError)
		return
	}

	// Store filtered result in cache (keyed by query parameters)
	productsCache.Set(cacheKey, out)

	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
