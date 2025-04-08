<template>
  <q-page class="row items-center justify-evenly">
    <h1>Inventory Page</h1>
    <ProductCardGrid :products="productList" />
  </q-page>
</template>
<script setup lang="ts">
import ProductCardGrid from 'components/ProductCardGrid.vue'
import type {Product} from '../types/products'
const productList: Product[] = [
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

]
</script>


