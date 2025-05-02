<script setup lang="ts">
import { computed } from 'vue'
import placeholder from "src/assets/images/placeholder.png";

const props = defineProps<{
  modelValue: boolean
  image: string
  title: string
  price: number
  unit_color: string
  quantity?: number  // Add quantity prop
  details?: string
  status?: string
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
  <q-dialog v-model="isOpen" persistent class="iphone-se">
    <q-card style="width: 45rem; max-width: 90vw" class="q-pa-none">
      <q-card-section>
        <div class="row q-pt-md q-px-md">
          <div class="col-12 col-md-8">
<!-- Image area -->
            <q-img
              :src="image || placeholder"
              ratio="4/3"
              height="240px"
              :placeholder-src="placeholder"
              :error-src="placeholder"
              spinner-color="primary"
              spinner-size="82px"
              class="product-rounded"
            >
              <template v-slot:loading>
                <q-spinner-dots color="primary" size="40px" />
              </template>
              <template v-slot:error>
                <div class="absolute-full flex flex-center bg-negative text-white">
                  Cannot load image
                </div>
              </template>
              <q-btn
                flat
                dense
                class="absolute-bottom-right q-mb-sm q-mr-sm"

              >
                <q-icon name="description" class="text-white" />
              </q-btn>
            </q-img>
          </div>
          <div class="col-12 col-md-4">
                  <!-- Information section with name, details, and price -->
            <q-card-section class="q-pt-md q-pb-none q-mb-md">
              <div class="row justify-between items-center">
                <!-- Left side - Product info -->
                <div class="col-7">
                  <div class="text-h6 text-bold">{{ title }}</div>
                  <div class="text-body2">{{ details }}</div>
                  <div class="text-body2">{{ unit_color }}</div>
                  <div class="col-5 text-h6 text-bold q-mt-md">
                    ₱{{ price.toLocaleString() }}
                  </div>
                  <div class="text-subtitle2">Quantity: {{ quantity }}</div>
                  <!--need to add this -->
                  <div class="text-h5 text-bold text-positive q-mt-md">{{ status }}</div>
                </div>
              </div>
            </q-card-section>
            <!-- Action buttons -->
            <q-card-actions align="right" class="q-pa-none">
                <q-btn color="primary" label="ADD" style="width: 30%;" @click="add" />
              </q-card-actions>
          </div>
        </div>
      </q-card-section>
<!-- Close button in the top-right -->
      <q-btn flat dense round class="absolute-top-right z-top q-mt-xs q-mr-xs"  @click="close"> <q-icon name="close" class="text-primary"/></q-btn>
    </q-card>
  </q-dialog>
</template>


<style scoped>
.z-top {
  z-index: 1000;
}

@media (max-width: 750px){
  .iphone-se {
    color: aqua;
  }
}
</style>
