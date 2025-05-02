<script setup lang="ts">
import { computed } from 'vue'
import placeholder from "src/assets/images/placeholder.png";

const props = defineProps<{
  modelValue: boolean
  image: string
  title: string
  price: number
  quantity?: number  // Add quantity prop
  details?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
  (e: 'add'): void
}>()

const isOpen = computed({
  get:  () => props.modelValue,
  set: v => emit('update:modelValue', v)
})

function close() {
  isOpen.value = false
}
function add() {
  emit('add')
  close()
}
</script>


<template>
  <q-dialog v-model="isOpen" persistent>
    <q-card style="width: 320px; max-width: 90vw" class="q-pa-none">
      <!-- Close button in the top-right -->
      <q-btn flat dense round class="absolute-top-right z-top q-mt-xs q-mr-xs"  @click="close"> <q-icon name="close" class="text-primary"/></q-btn>

      <!-- Image area -->
      <q-img
        :src="image || placeholder"
        ratio="4/3"
        height="240px"
        :placeholder-src="placeholder"
        :error-src="placeholder"
        spinner-color="primary"
        spinner-size="82px"
      >
        <template v-slot:loading>
          <q-spinner-dots color="primary" size="40px" />
        </template>
        <template v-slot:error>
          <div class="absolute-full flex flex-center bg-negative text-white">
            Cannot load image
          </div>
        </template>
      </q-img>

      <!-- Information section with name, details, and price -->
      <q-card-section class="q-pt-md q-pb-none">
        <div class="row justify-between items-center">
          <!-- Left side - Product info -->
          <div class="col-7">
            <div class="text-h6">{{ title }}</div>
            <div class="text-subtitle2">Quantity: {{ quantity }}</div>
            <div class="text-body2">{{ details }}</div>
          </div>

          <!-- Right side - Price -->
          <div class="col-5 text-right text-h6">
            ₱{{ price.toLocaleString() }}
          </div>
        </div>
      </q-card-section>

      <!-- Action buttons -->
      <q-card-actions align="right" class="q-pa-md">
        <q-btn color="primary" label="ADD" @click="add" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>


<style scoped>
.z-top {
  z-index: 1000;
}
</style>
