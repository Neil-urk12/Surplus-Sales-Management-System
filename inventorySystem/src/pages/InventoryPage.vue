<template>
  <q-page class="row items-center justify-evenly">
    <h1>Inventory Page</h1>
    <ProductCardGrid :products="inventoryProducts" />
  </q-page>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import ProductCardGrid, { type Product } from 'components/ProductCardGrid.vue'

const inventoryProducts = ref<Product[]>([
  {
    id: 1,
    name: 'Wireless Headphones',
    description: 'High-quality noise-canceling wireless headphones with 30-hour battery life.',
    price: 199.99,
    stock: 15,
    image: 'https://picsum.photos/400/300?random=1'
  },
  {
    id: 2,
    name: 'Wireless Headphones',
    description: 'High-quality noise-canceling wireless headphones with 30-hour battery life.',
    price: 199.99,
    stock: 15,
    image: 'https://picsum.photos/400/300?random=1'
  },
  {
    id: 3,
    name: 'Wireless Headphones',
    description: 'High-quality noise-canceling wireless headphones with 30-hour battery life.',
    price: 199.99,
    stock: 15,
    image: 'https://picsum.photos/400/300?random=1'
  },
  {
    id: 4,
    name: 'Wireless Headphones',
    description: 'High-quality noise-canceling wireless headphones with 30-hour battery life.',
    price: 199.99,
    stock: 0,
    image: 'https://picsum.photos/400/300?random=1'
  },


])
</script>

<style scoped>
.sr-only {
  position: absolute !important;
  width: 1px !important;
  height: 1px !important;
  padding: 0 !important;
  margin: -1px !important;
  overflow: hidden !important;
  clip: rect(0, 0, 0, 0) !important;
  border: 0 !important;
}

</style>

