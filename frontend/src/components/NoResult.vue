<template>
  <div
    class="flex flex-col items-center justify-center min-h-[400px] px-6 rounded-lg"
    style="background-color: #f8f9fa"
  >
    <div
      class="flex items-center justify-center w-20 h-20 rounded-full mb-6"
      style="background-color: #eeeeee"
    >
      <span
        class="pi pi-box text-4xl"
        style="color: #888888"
        aria-hidden="true"
      ></span>
    </div>

    <h2 class="text-xl font-bold mb-2 text-center" style="color: #333333">
      No products found
    </h2>

    <p
      class="text-center text-sm max-w-[500px] mb-6"
      style="color: #666666; line-height: 1.5"
    >
      We couldn't find any products matching your current filters.<br />
      Try adjusting your search criteria.
    </p>

    <Button variant="secondary" @click="$emit('clear-filters')">
      Clear all filters
    </Button>
  </div>
</template>

<script setup>
import Button from "./Button.vue";

defineEmits(["clear-filters"]);
</script>
