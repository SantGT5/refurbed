<template>
  <article
    class="relative group bg-white rounded-xl shadow-md cursor-pointer transition-shadow hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
    role="button"
    tabindex="0"
    @click="$emit('click', product)"
    @keydown.enter="$emit('click', product)"
    @keydown.space.prevent="$emit('click', product)"
  >
    <div class="relative aspect-[4/3] bg-gray-100 overflow-hidden rounded-t-xl">
      <img
        v-if="product.image_url"
        :src="product.image_url"
        :alt="product.name"
        class="w-full h-full object-contain"
      />
      <div
        v-else
        class="w-full h-full flex items-center justify-center text-gray-400"
      >
        <i class="pi pi-image" style="font-size: 3rem"></i>
      </div>

      <span
        v-if="product.bestseller"
        class="absolute top-0 left-0 z-10 px-2.5 py-1 rounded-1 bg-green-800 text-white text-[9px] font-medium"
      >
        Bestseller
      </span>
    </div>

    <span
      v-if="product.discount_percent"
      class="absolute -top-3 -right-3 z-20 w-10 h-10 rounded-full bg-purple-600 text-white text-sm font-medium flex items-center justify-center rotate-12 shadow-md"
    >
      -{{ product.discount_percent }}%
    </span>

    <div class="p-4 flex flex-col gap-3 rounded-b-xl overflow-hidden">
      <h2 class="text-lg font-semibold text-indigo-900">
        {{ product.name }}
      </h2>

      <div v-if="product.colors?.length" class="flex items-center gap-1.5">
        <span
          v-for="(colorName, index) in product.colors"
          :key="colorName"
          class="w-4 h-4 rounded-full shrink-0 border border-gray-200"
          :class="[colorSwatchClass(colorName)]"
          :style="colorSwatchStyle(colorName)"
          :title="colorName"
        />
      </div>

      <p class="text-xl font-medium text-green-600 mt-auto">
        {{ formattedPrice }}
      </p>
    </div>
  </article>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  product: {
    type: Object,
    required: true,
    validator: (p) =>
      p &&
      typeof p.id !== "undefined" &&
      typeof p.name !== "undefined" &&
      typeof p.base_price !== "undefined",
  },
});

defineEmits(["click"]);

const knownColors = {
  white: "bg-gray-100",
  black: "bg-gray-900",
  blue: "bg-blue-600",
  red: "bg-red-600",
  green: "bg-green-600",
  pink: "bg-pink-400",
  gray: "bg-gray-500",
  silver: "bg-gray-300",
  gold: "bg-amber-400",
  yellow: "bg-yellow-400",
  purple: "bg-purple-600",
  midnight: "bg-slate-900",
  starlight: "bg-stone-200",
};

function colorSwatchClass(colorName) {
  const key = String(colorName).toLowerCase();
  return knownColors[key] || "";
}

function colorSwatchStyle(colorName) {
  const key = String(colorName).toLowerCase();
  if (knownColors[key]) return {};
  return { backgroundColor: key };
}

const formattedPrice = computed(() => {
  const price = props.product.discount_percent
    ? props.product.base_price * (1 - props.product.discount_percent / 100)
    : props.product.base_price;
  return new Intl.NumberFormat("de-DE", {
    style: "currency",
    currency: "EUR",
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(price);
});
</script>
