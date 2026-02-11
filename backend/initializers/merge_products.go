package initializers

import (
	"encoding/json"
	"log"
	"os"
)

func MergeProducts() {
	if _, err := os.Stat("data/products.json"); err == nil {
		return
	}

	metadata, err := os.ReadFile("data/metadata.json")
	if err != nil {
		log.Fatal(err)
	}

	details, err := os.ReadFile("data/details.json")
	if err != nil {
		log.Fatal(err)
	}

	merged := make(map[string]map[string]interface{})

	var metadataItems []map[string]interface{}
	var detailsItems []map[string]interface{}

	json.Unmarshal(metadata, &metadataItems)
	json.Unmarshal(details, &detailsItems)

	mergeItem := func(item map[string]interface{}) {
		id := item["id"].(string)

		// create map if not exists
		if _, ok := merged[id]; !ok {
			merged[id] = make(map[string]interface{})
		}

		// copy fields
		for k, v := range item {
			merged[id][k] = v
		}
	}

	// merge metadata
	for _, metadataItem := range metadataItems {
		mergeItem(metadataItem)
	}

	// merge details
	for _, detailItem := range detailsItems {
		mergeItem(detailItem)
	}

	mergedJSON, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// write to products.json
	os.WriteFile("data/products.json", mergedJSON, 0644)
}
