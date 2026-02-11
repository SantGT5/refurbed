<template>
  <div class="container mx-auto p-8">
    <h1 class="text-3xl font-bold mb-4">Product Listing</h1>
    <Input
      size="md"
      name="search"
      variant="outline"
      :modelValue="form.search"
      placeholder="Search products..."
      @update:modelValue="(value) => (form.search = value)"
    >
      <template #left-icon>
        <i class="pi pi-search" style="font-size: 1.25rem"></i>
      </template>
    </Input>

    <div class="flex gap-4 mt-4">
      <SideBar title="Filters">
        <DropDown
          name="color"
          v-model="form.color"
          placeholder="Select color"
          :options="[
            { label: 'Blue', value: 'blue' },
            { label: 'Gray', value: 'gray' },
            { label: 'Pink', value: 'pink' },
            { label: 'Green', value: 'green' },
            { label: 'Silver', value: 'silver' },
          ]"
        />

        <Range
          title="Price Range"
          :min="0"
          :max="2000"
          :model-value="{
            min: Number(form.minPrice) || 0,
            max: Number(form.maxPrice) || 2000,
          }"
          @update:model-value="
            (v) => {
              form.minPrice = v.min;
              form.maxPrice = v.max;
            }
          "
        />

        <DropDown
          name="bestseller"
          v-model="form.bestseller"
          placeholder="Select bestseller"
          :options="[
            { label: 'Bestseller', value: true },
            { label: 'Not bestseller', value: false },
          ]"
        />
      </SideBar>

      <div
        class="flex-1 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
      >
        <template v-if="loadingProducts">
          <Loading class="col-span-full" />
        </template>

        <template v-else-if="errorProducts">
          <Alert
            title="Error"
            variant="error"
            :description="errorProducts?.data?.error || 'Something went wrong'"
          />
        </template>

        <template
          v-else-if="
            Array.isArray(Object.values(products || {})) &&
            Object.values(products || {})?.length
          "
        >
          <ProductCard
            v-for="product in products"
            :key="product.id"
            :product="product"
            @click="(p) => console.log('Clicked product', p)"
          />
        </template>

        <template v-else-if="!Object.values(products || {})?.length">
          <NoResult class="col-span-full" @clear-filters="resetForm" />
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from "vue";

import { useFetch } from "./composables";
import {
  buildQueryString,
  getQueryParams,
  updateUrlWithQueries,
} from "./utils/url";
import {
  SideBar,
  Input,
  DropDown,
  Range,
  ProductCard,
  NoResult,
  Loading,
  Alert,
} from "./components";

const { search, color, bestseller, minPrice, maxPrice } = getQueryParams(
  window.location.search,
);

const form = ref({ search, color, bestseller, minPrice, maxPrice });

const {
  data: products,
  error: errorProducts,
  execute: fetchProducts,
  loading: loadingProducts,
} = useFetch({
  url: buildQueryString({
    baseUrl: "/products",
    params: { color, search, minPrice, maxPrice, bestseller },
  }),
  onError: (error) => {
    console.error("Error fetching products:", error?.data?.error);
  },
});

let debounceTimeout = null;
watch(
  () => form.value,
  (newForm) => {
    if (debounceTimeout) {
      clearTimeout(debounceTimeout);
    }

    // Set loading to true immediately when filters change
    loadingProducts.value = true;

    debounceTimeout = setTimeout(() => {
      updateUrlWithQueries({
        params: {
          color: newForm.color,
          search: newForm.search,
          minPrice: newForm.minPrice,
          maxPrice: newForm.maxPrice,
          bestseller: newForm.bestseller,
        },
        replaceHistory: true,
      });

      fetchProducts({
        url: buildQueryString({
          baseUrl: "/products",
          params: {
            color: newForm.color,
            search: newForm.search,
            minPrice: newForm.minPrice,
            maxPrice: newForm.maxPrice,
            bestseller: newForm.bestseller,
          },
        }),
      });
    }, 500);
  },
  { deep: true, immediate: false },
);

const resetForm = async () => {
  form.value = {
    search: "",
    color: "",
    minPrice: 0,
    bestseller: "",
    maxPrice: 2000,
  };
};
</script>
