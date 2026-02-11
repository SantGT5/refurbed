<template>
  <div class="input-container">
    <div
      class="input-wrapper flex items-center gap-2 rounded-lg border border-gray-300 bg-white transition-colors focus-within:ring-2 focus-within:ring-blue-500 focus-within:border-transparent"
      :class="{
        'opacity-60 cursor-not-allowed bg-gray-100': disabled,
        'w-full': size === 'full',
        'max-w-[300px]': size === 'sm',
        'max-w-[400px]': size === 'md',
        'max-w-[500px]': size === 'lg',
      }"
    >
      <span
        v-if="$slots['left-icon']"
        class="input-left-icon flex-shrink-0 text-gray-500 pl-3"
      >
        <slot name="left-icon" />
      </span>

      <input
        :value="modelValue"
        :name="name"
        :disabled="disabled"
        type="text"
        class="input-field flex-1 min-w-0 py-2 bg-transparent text-gray-900 placeholder-gray-500 focus:outline-none disabled:cursor-not-allowed border-0 rounded-lg"
        :class="hasLeftIcon ? 'pl-1 pr-4' : 'px-4'"
        v-bind="$attrs"
        @input="onInput"
      />
    </div>

    <p v-if="errorMessage" class="text-red-500 text-sm mt-1">
      {{ errorMessage }}
    </p>
  </div>
</template>

<script setup>
import { useSlots, computed } from "vue";

defineOptions({
  inheritAttrs: false,
});

const slots = useSlots();
const hasLeftIcon = computed(() => !!slots["left-icon"]);

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: "",
  },
  name: {
    type: String,
    required: true,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  size: {
    type: String,
    default: "w-full",
    validator: (v) => ["sm", "md", "lg", "full"].includes(v),
  },
  errorMessage: {
    type: String,
    default: "",
  },
});

const emit = defineEmits(["update:modelValue"]);

function onInput(event) {
  emit("update:modelValue", event.target.value);
}
</script>
