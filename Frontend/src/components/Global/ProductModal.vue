<script setup lang="ts">
import { computed }from 'vue'
import placeholder from "src/assets/images/placeholder.png";

const props = defineProps<{
  modelValue: boolean
  image: string
  title: string
  price?: number
  unit_color: string
  quantity?: number  // Add quantity prop
  details?: string
  status?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
  (e: 'addItem'): void
}>()

const isOpen = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v)
})

const validatedImage = computed(() => {
  try {
    if (!props.image) {
      console.log('No image provided, using placeholder');
      return placeholder;
    }

    // Handle base64 images
    if (props.image.startsWith('data:image/')) {
      console.log('Processing base64 image');
      return props.image; // Use the image as-is, let the img tag handle errors
    }

    // Handle relative paths from src/assets
    if (props.image.startsWith('src/assets/')) {
      try {
        // Try to import the image directly
        return new URL(`/${props.image}`, import.meta.url).href;
      } catch (error) {
        console.error('Error loading asset image:', error);
        return placeholder;
      }
    }

    // Handle URLs and blob URLs
    if (props.image.startsWith('http') || props.image.startsWith('blob:')) {
      return props.image;
    }

    // Handle absolute paths
    if (props.image.startsWith('/')) {
      return props.image;
    }

    // For any other case, try to use the image as is
    console.log('Attempting to use image as is:', props.image);
    return props.image;
  } catch (error) {
    console.error('Error in image validation:', error);
    return placeholder;
  }
});

const getStatusColor = computed(() => {
  switch (props.status?.toLowerCase()) {
    case 'available':
      return 'text-positive'
    case 'in stock':
      return 'text-positive'
    case 'low stock':
      return 'text-warning'
    case 'out of stock':
      return 'text-negative'
    default:
      return 'text-grey'
  }
})

function handleImageError(event: Event) {
  console.error('Image failed to load:', props.image);
  const imgElement = event.target as HTMLImageElement;
  imgElement.src = placeholder;
}

function close() {
  isOpen.value = false
}

</script>


<template>
  <q-dialog v-model="isOpen" dismissible >
    <q-card style="width: 45rem; max-width: 90vw" class="q-pa-none">
      <q-card-section>
        <div class="row q-pt-md q-px-md">
          <div class="col-12 col-md-8">
            <!-- Image area -->
            <q-img
              :src="validatedImage"
              :ratio="16/9"
              height="240px"
              :placeholder-src="placeholder"
              :error-src="placeholder"
              spinner-color="primary"
              spinner-size="82px"
              class="product-rounded"
              fit="contain"
              @error="handleImageError"
            >
              <template v-slot:loading>
                <div class="absolute-full flex flex-center bg-grey-2 text-grey-8">
                  <q-spinner-dots color="primary" size="40px" />
                </div>
              </template>
              <template v-slot:error>
                <div class="absolute-full flex flex-center bg-grey-2 text-grey-8">
                  <q-icon name="error" size="40px" color="grey" />
                </div>
              </template>
            </q-img>
          </div>
          <div class="col-12 col-md-4">
            <!-- Information section with name, details, and price -->
            <q-card-section class="q-pt-md q-pb-none q-mb-md">
              <div class="row justify-between items-center">
                <!-- Left side - Product info -->
                <div class="col-12">
                  <div class="text-h6 text-bold">{{ title }}</div>
                  <div class="text-body2">{{ details }}</div>
                  <div class="text-body2">{{ unit_color }}</div>
                  <div class="text-subtitle2">Quantity: {{ quantity }}</div>
                  <div class="text-h5 text-bold" :class="getStatusColor">{{ status }}</div>
                </div>
              </div>
            </q-card-section>
          </div>
        </div>
      </q-card-section>
      <!-- Close button in the top-right -->
      <q-btn flat dense round class="close-btn" @click="close">
        <q-icon name="close" class="close-icon" />
      </q-btn>
    </q-card>
  </q-dialog>
</template>


<style scoped>
.z-top {
  z-index: 1000;
}

.close-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
}

.close-icon {
  font-size: 1.5rem;
}

:deep(.q-dark) .close-icon {
  color: #fff !important;
}

body.body--light .close-icon {
  color: rgba(0, 0, 0, 0.7) !important;
}

.product-rounded {
  border-radius: 8px;
  overflow: hidden;
}

.image-container {
  background-color: var(--q-dark);
  border-radius: 8px;
  overflow: hidden;
}

.q-img {
  background: transparent;
}

:deep(.q-img__content > div) {
  background: transparent;
}
</style>
