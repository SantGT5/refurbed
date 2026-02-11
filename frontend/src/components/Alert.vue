<template>
  <div role="alert" class="rounded-lg p-4 w-full h-fit" :class="alertClasses">
    <h4 v-if="title" class="font-semibold mb-1">{{ title }}</h4>
    <p v-if="description" :class="title ? 'text-sm' : ''">{{ description }}</p>
    <slot />
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  variant: {
    type: String,
    default: "info",
    validator: (v) => ["info", "warning", "error"].includes(v),
  },
  title: {
    type: String,
    default: "",
  },
  description: {
    type: String,
    default: "",
  },
});

const alertClasses = computed(() => {
  const base = "border";
  const variants = {
    info: "bg-blue-50 border-blue-200 text-blue-800",
    warning: "bg-amber-50 border-amber-200 text-amber-800",
    error: "bg-red-50 border-red-200 text-red-800",
  };
  return `${base} ${variants[props.variant]}`;
});
</script>
