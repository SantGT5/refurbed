<template>
  <div ref="rootRef" class="dropdown-wrapper relative w-full">
    <button
      type="button"
      :disabled="disabled"
      :name="name"
      class="dropdown-trigger flex items-center justify-between w-full rounded-lg border border-gray-300 bg-white py-2 px-4 text-left transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
      :class="{
        'opacity-60 cursor-not-allowed bg-gray-100': disabled,
        'ring-2 ring-blue-500 border-transparent': isOpen,
      }"
      aria-haspopup="listbox"
      :aria-expanded="isOpen"
      aria-labelledby="dropdown-label"
      @click="toggle"
      @keydown="onTriggerKeydown"
    >
      <span
        id="dropdown-label"
        class="flex-1 min-w-0 truncate"
        :class="selectedLabel ? 'text-gray-900' : 'text-gray-500'"
      >
        {{ selectedLabel || placeholder }}
      </span>
      <span
        class="flex-shrink-0 ml-2 text-gray-500 transition-transform"
        :class="{ 'rotate-180': isOpen }"
      >
        <slot name="right-icon">
          <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M19 9l-7 7-7-7"
            />
          </svg>
        </slot>
      </span>
    </button>

    <Transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <ul
        v-show="isOpen"
        ref="listRef"
        role="listbox"
        :aria-activedescendant="activeId"
        tabindex="-1"
        class="dropdown-list absolute z-10 mt-1 w-full max-h-60 overflow-auto rounded-lg border border-gray-200 bg-white py-1 shadow-lg focus:outline-none"
        @keydown="onListKeydown"
      >
        <li
          v-for="(option, index) in normalizedOptions"
          :id="`dropdown-option-${index}`"
          :key="getOptionValue(option)"
          role="option"
          :aria-selected="isSelected(option)"
          class="dropdown-option cursor-pointer py-2 px-4 text-gray-900 hover:bg-gray-100 focus:bg-gray-100 focus:outline-none"
          :class="{
            'bg-blue-50 text-blue-700 hover:bg-blue-100': isSelected(option),
          }"
          @click="select(option)"
          @mousemove="activeIndex = index"
        >
          {{ getOptionLabel(option) }}
        </li>
        <li
          v-if="normalizedOptions.length === 0"
          class="py-2 px-4 text-gray-500 text-sm"
        >
          No options
        </li>
      </ul>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, watch } from "vue";

defineOptions({
  inheritAttrs: false,
});

const props = defineProps({
  modelValue: {
    type: [String, Number, Boolean, Object],
    default: null,
  },
  name: {
    type: String,
    default: "",
  },
  placeholder: {
    type: String,
    default: "Select...",
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  options: {
    type: Array,
    default: () => [],
    validator: (val) =>
      val.every(
        (o) =>
          typeof o === "string" ||
          (typeof o === "object" && o !== null && ("label" in o || "value" in o))
      ),
  },
});

const emit = defineEmits(["update:modelValue"]);

const isOpen = ref(false);
const activeIndex = ref(-1);
const rootRef = ref(null);
const listRef = ref(null);
let clickOutsideHandler = null;

const normalizedOptions = computed(() =>
  props.options.map((o) =>
    typeof o === "string" ? { label: o, value: o } : { ...o, label: o.label ?? o.value, value: o.value }
  )
);

const selectedLabel = computed(() => {
  const opt = normalizedOptions.value.find((o) => o.value === props.modelValue);
  return opt ? opt.label : null;
});

const activeId = computed(() =>
  activeIndex.value >= 0 && activeIndex.value < normalizedOptions.value.length
    ? `dropdown-option-${activeIndex.value}`
    : undefined
);

function getOptionLabel(option) {
  return option?.label ?? option?.value ?? "";
}

function getOptionValue(option) {
  return option?.value;
}

function isSelected(option) {
  return props.modelValue === getOptionValue(option);
}

function toggle() {
  if (props.disabled) return;
  isOpen.value = !isOpen.value;
  if (isOpen.value) {
    activeIndex.value = normalizedOptions.value.findIndex((o) =>
      isSelected(o)
    );
    if (activeIndex.value < 0) activeIndex.value = 0;
  }
}

function select(option) {
  const value = getOptionValue(option);
  emit("update:modelValue", value);
  isOpen.value = false;
  activeIndex.value = -1;
}

function close() {
  isOpen.value = false;
  activeIndex.value = -1;
  if (clickOutsideHandler) {
    document.removeEventListener("click", clickOutsideHandler);
    clickOutsideHandler = null;
  }
}

function onTriggerKeydown(e) {
  if (e.key === "Enter" || e.key === " ") {
    e.preventDefault();
    toggle();
  }
  if (e.key === "ArrowDown" && !isOpen.value) {
    e.preventDefault();
    isOpen.value = true;
    activeIndex.value = 0;
  }
}

function onListKeydown(e) {
  const n = normalizedOptions.value.length;
  if (e.key === "Escape") {
    e.preventDefault();
    close();
    rootRef.value?.querySelector("button")?.focus();
    return;
  }
  if (e.key === "Enter" || e.key === " ") {
    e.preventDefault();
    if (activeIndex.value >= 0 && activeIndex.value < n) {
      select(normalizedOptions.value[activeIndex.value]);
    }
    return;
  }
  if (e.key === "ArrowDown") {
    e.preventDefault();
    activeIndex.value = Math.min(activeIndex.value + 1, n - 1);
    return;
  }
  if (e.key === "ArrowUp") {
    e.preventDefault();
    activeIndex.value = Math.max(activeIndex.value - 1, 0);
  }
}

watch(isOpen, (open) => {
  if (!open) return;
  clickOutsideHandler = (e) => {
    if (rootRef.value && !rootRef.value.contains(e.target)) {
      close();
    }
  };
  requestAnimationFrame(() =>
    document.addEventListener("click", clickOutsideHandler)
  );
});
</script>
