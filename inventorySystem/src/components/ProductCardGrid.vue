<!-- components/ProductCardGrid.vue -->
<template>
  <div class="q-pa-md">
    <div class="row q-col-gutter-md">
      <div
        v-for="product in products"
        :key="product.id"
        class="col-12 col-sm-6 col-md-4"
      >
        <q-card class="my-card">
          <q-img
            :src="product.image"
            :ratio="16/9"
            basic
          >
            <div class="absolute-bottom text-subtitle2 text-center">
              {{ product.name }}
            </div>
          </q-img>

          <q-card-section>
            <div class="text-caption text-grey-8 q-mb-xs">
              {{ product.description }}
            </div>
            <div class="row items-center justify-between">
              <div class="text-h6 text-primary">
                ${{ product.price }}
              </div>
              <div
                class="text-caption"
                :class="product.stock > 0 ? 'text-green' : 'text-red'"
              >
                {{ product.stock > 0 ? 'In Stock' : 'Out of Stock' }}
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface Product {
  id: number
  name: string
  description: string
  price: number
  stock: number
  image: string
}

// Properly define and use props
const { products } = defineProps({
  products: {
    type: Array as () => Product[],
    required: true
  }
})
</script>

<style lang="scss" scoped>
.my-card {
  transition: transform 0.3s;
  &:hover {
    transform: translateY(-5px);
  }

  .text-caption {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}
</style>
