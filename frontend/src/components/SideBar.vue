<template>
  <aside
    class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 w-fit h-fit"
    role="complementary"
  >
    <h2 v-if="title" class="text-lg font-semibold text-gray-900 mb-4">
      {{ title }}
    </h2>
    <div class="flex flex-col gap-6">
      <slot />
    </div>
  </aside>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: "",
  },
});
</script>
