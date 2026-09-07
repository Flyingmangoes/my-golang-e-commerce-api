<script setup lang="ts">
import ProductCard from './product_card.vue'
import type { Product } from '@/lib/models/Product'

defineProps<{
  products: Product[]
  emptyMessage?: string
}>()

const emit = defineEmits<{
  selectProduct: [product: Product]
}>()
</script>

<template>
  <div
    v-if="products.length > 0"
    class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-x-8 gap-y-14 lg:gap-x-20 lg:gap-y-16"
  >
    <ProductCard
      v-for="product in products"
      :key="product.productID"
      :product="product"
      @select="emit('selectProduct', $event)"
    />
  </div>

  <div v-else class="w-full py-24 text-center text-neutral-500">
    {{ emptyMessage ?? 'No products found.' }}
  </div>
</template>