<template>
  <div class="range-wrapper flex flex-col gap-3 w-full">
    <span v-if="title" class="range-title text-sm font-semibold text-gray-900">
      {{ title }}
    </span>

    <div
      ref="trackRef"
      class="range-track relative h-2 w-full rounded-full bg-gray-200 cursor-pointer select-none touch-none"
      @mousedown="onTrackMouseDown"
      @touchstart.passive="onTrackTouchStart"
    >
      <!-- Selected segment (between thumbs) -->
      <div
        class="range-fill absolute h-full rounded-full bg-gray-800 pointer-events-none"
        :style="fillStyle"
      />
      <!-- Min thumb -->
      <div
        ref="minThumbRef"
        role="slider"
        tabindex="0"
        aria-valuemin="minBound"
        aria-valuemax="localMax"
        :aria-valuenow="localMin"
        aria-label="Minimum value"
        class="range-thumb absolute top-1/2 -translate-y-1/2 w-5 h-5 rounded-full bg-white border-2 border-gray-800 shadow cursor-grab active:cursor-grabbing focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
        :style="{ left: `${minPercent}%`, marginLeft: '-10px' }"
        @mousedown="startDrag('min', $event)"
        @touchstart="startDrag('min', $event)"
        @keydown="onThumbKeydown('min', $event)"
      />
      <!-- Max thumb -->
      <div
        ref="maxThumbRef"
        role="slider"
        tabindex="0"
        aria-valuemin="localMin"
        aria-valuemax="maxBound"
        :aria-valuenow="localMax"
        aria-label="Maximum value"
        class="range-thumb absolute top-1/2 -translate-y-1/2 w-5 h-5 rounded-full bg-white border-2 border-gray-800 shadow cursor-grab active:cursor-grabbing focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
        :style="{ left: `${maxPercent}%`, marginLeft: '-10px' }"
        @mousedown="startDrag('max', $event)"
        @touchstart="startDrag('max', $event)"
        @keydown="onThumbKeydown('max', $event)"
      />
    </div>

    <div class="range-values flex justify-between text-sm">
      <div class="flex flex-col">
        <span class="text-gray-500">Minimum:</span>
        <span class="font-medium text-gray-900">{{
          formatValue(localMin)
        }}</span>
      </div>
      <div class="flex flex-col items-end">
        <span class="text-gray-500">Maximum:</span>
        <span class="font-medium text-gray-900">{{
          formatValue(localMax)
        }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onUnmounted } from "vue";

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({ min: 0, max: 100 }),
    validator: (v) =>
      v &&
      typeof v.min === "number" &&
      typeof v.max === "number" &&
      v.min <= v.max,
  },
  min: {
    type: Number,
    default: 0,
  },
  max: {
    type: Number,
    default: 500,
  },
  title: {
    type: String,
    default: "",
  },
  name: {
    type: String,
    default: undefined,
  },
  currency: {
    type: String,
    default: "€",
  },
  disabled: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["update:modelValue"]);

const trackRef = ref(null);
const minThumbRef = ref(null);
const maxThumbRef = ref(null);

const minBound = computed(() => props.min);
const maxBound = computed(() => props.max);

const localMin = ref(
  Math.round(
    Math.max(
      props.min,
      Math.min(props.max, props.modelValue?.min ?? props.min),
    ),
  ),
);
const localMax = ref(
  Math.round(
    Math.max(
      props.min,
      Math.min(props.max, props.modelValue?.max ?? props.max),
    ),
  ),
);

watch(
  () => props.modelValue,
  (v) => {
    if (v?.min != null) localMin.value = Math.round(clamp(v.min));
    if (v?.max != null) localMax.value = Math.round(clamp(v.max));
  },
  { deep: true },
);

watch(
  [localMin, localMax],
  () => {
    if (localMin.value > localMax.value) {
      const mid = Math.round((localMin.value + localMax.value) / 2);
      localMin.value = mid;
      localMax.value = mid;
    }
    emit("update:modelValue", {
      min: Math.round(localMin.value),
      max: Math.round(localMax.value),
    });
  },
  { immediate: true },
);

function clamp(v) {
  return Math.max(minBound.value, Math.min(maxBound.value, v));
}

const range = computed(() => maxBound.value - minBound.value || 1);

const minPercent = computed(
  () => ((localMin.value - minBound.value) / range.value) * 100,
);
const maxPercent = computed(
  () => ((localMax.value - minBound.value) / range.value) * 100,
);

const fillStyle = computed(() => ({
  left: `${minPercent.value}%`,
  width: `${maxPercent.value - minPercent.value}%`,
}));

function valueFromPercent(percent) {
  const p = Math.max(0, Math.min(100, percent));
  return Math.round(minBound.value + (p / 100) * range.value);
}

function percentFromClientX(clientX) {
  const el = trackRef.value;
  if (!el) return 0;
  const rect = el.getBoundingClientRect();
  const x = clientX - rect.left;
  return (x / rect.width) * 100;
}

let dragging = null;

function startDrag(which, event) {
  if (props.disabled) return;
  event.preventDefault();
  dragging = which;
  if (event.type === "mousedown") {
    window.addEventListener("mousemove", onMove);
    window.addEventListener("mouseup", stopDrag);
  } else if (event.type === "touchstart") {
    window.addEventListener("touchmove", onMoveTouch, { passive: false });
    window.addEventListener("touchend", stopDrag);
  }
}

function onMove(event) {
  if (dragging == null) return;
  const clientX = event.clientX;
  const percent = percentFromClientX(clientX);
  const value = valueFromPercent(percent);
  if (dragging === "min") {
    localMin.value = Math.min(value, localMax.value);
  } else {
    localMax.value = Math.max(value, localMin.value);
  }
}

function onMoveTouch(event) {
  if (event.cancelable) event.preventDefault();
  if (dragging == null || !event.touches[0]) return;
  const clientX = event.touches[0].clientX;
  const percent = percentFromClientX(clientX);
  const value = valueFromPercent(percent);
  if (dragging === "min") {
    localMin.value = Math.min(value, localMax.value);
  } else {
    localMax.value = Math.max(value, localMin.value);
  }
}

function stopDrag() {
  dragging = null;
  window.removeEventListener("mousemove", onMove);
  window.removeEventListener("mouseup", stopDrag);
  window.removeEventListener("touchmove", onMoveTouch);
  window.removeEventListener("touchend", stopDrag);
}

function onTrackMouseDown(event) {
  if (props.disabled || event.target?.closest?.(".range-thumb")) return;
  const percent = percentFromClientX(event.clientX);
  const value = valueFromPercent(percent);
  const distToMin = Math.abs(value - localMin.value);
  const distToMax = Math.abs(value - localMax.value);
  dragging = distToMin <= distToMax ? "min" : "max";
  if (dragging === "min") localMin.value = Math.min(value, localMax.value);
  else localMax.value = Math.max(value, localMin.value);
  window.addEventListener("mousemove", onMove);
  window.addEventListener("mouseup", stopDrag);
}

function onTrackTouchStart(event) {
  if (props.disabled || event.target?.closest?.(".range-thumb")) return;
  const touch = event.touches[0];
  const percent = percentFromClientX(touch.clientX);
  const value = valueFromPercent(percent);
  const distToMin = Math.abs(value - localMin.value);
  const distToMax = Math.abs(value - localMax.value);
  dragging = distToMin <= distToMax ? "min" : "max";
  if (dragging === "min") localMin.value = Math.min(value, localMax.value);
  else localMax.value = Math.max(value, localMin.value);
  window.addEventListener("touchmove", onMoveTouch, { passive: false });
  window.addEventListener("touchend", stopDrag);
}

function onThumbKeydown(which, event) {
  if (props.disabled) return;
  const step = range.value / 100;
  let delta = 0;
  switch (event.key) {
    case "ArrowRight":
    case "ArrowUp":
      delta = step;
      break;
    case "ArrowLeft":
    case "ArrowDown":
      delta = -step;
      break;
    case "Home":
      if (which === "min") localMin.value = minBound.value;
      else localMax.value = minBound.value;
      event.preventDefault();
      return;
    case "End":
      if (which === "min") localMin.value = maxBound.value;
      else localMax.value = maxBound.value;
      event.preventDefault();
      return;
    default:
      return;
  }
  event.preventDefault();
  if (which === "min") {
    localMin.value = Math.round(clamp(localMin.value + delta));
    if (localMin.value > localMax.value) localMax.value = localMin.value;
  } else {
    localMax.value = Math.round(clamp(localMax.value + delta));
    if (localMax.value < localMin.value) localMin.value = localMax.value;
  }
}

function formatValue(value) {
  const num = Math.round(value);
  return props.currency ? `${props.currency}${num}` : String(num);
}

onUnmounted(stopDrag);
</script>
