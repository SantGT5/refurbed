<template>
  <button
    :type="type"
    :disabled="disabled"
    class="px-4 py-2 rounded-lg font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-60 disabled:cursor-not-allowed disabled:pointer-events-none"
    :class="buttonClasses"
    v-bind="$attrs"
  >
    <slot />
  </button>
</template>

<script setup>
import { computed } from "vue";

defineOptions({
  inheritAttrs: false,
});

const props = defineProps({
  disabled: {
    type: Boolean,
    default: false,
  },
  variant: {
    type: String,
    default: "primary",
    validator: (v) => ["primary", "secondary", "outline", "ghost"].includes(v),
  },
  type: {
    type: String,
    default: "button",
    validator: (v) => ["button", "submit", "reset"].includes(v),
  },
});

const buttonClasses = computed(() => {
  const base = "inline-flex items-center justify-center";
  const variants = {
    primary:
      "bg-blue-600 text-white hover:bg-blue-700 active:bg-blue-800 disabled:bg-blue-400",
    secondary:
      "bg-gray-200 text-gray-900 hover:bg-gray-300 active:bg-gray-400 disabled:bg-gray-100",
    outline:
      "border-2 border-gray-300 bg-transparent text-gray-700 hover:bg-gray-50 active:bg-gray-100 disabled:border-gray-200 disabled:text-gray-400",
    ghost:
      "bg-transparent text-gray-700 hover:bg-gray-100 active:bg-gray-200 disabled:text-gray-400",
  };
  return `${base} ${variants[props.variant]}`;
});
</script>
