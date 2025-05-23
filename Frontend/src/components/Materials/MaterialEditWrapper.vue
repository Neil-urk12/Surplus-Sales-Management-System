<script setup lang="ts">
import { ref, watch } from 'vue';
import type { PropType } from 'vue';
import type { MaterialRow, NewMaterialInput } from 'src/types/materials';
import EditMaterialDialog from './EditMaterialDialog.vue';

const props = defineProps({
  open: {
    type: Boolean,
    required: true,
  },
  materialData: {
    type: Object as PropType<MaterialRow | null>,
    required: true,
  },
  categories: {
    type: Array as PropType<string[]>,
    required: true,
  },
  suppliers: {
    type: Array as PropType<string[]>,
    required: true,
  },
  defaultImageUrl: {
    type: String,
    required: true,
  },
});

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void;
  (e: 'update-material', materialData: NewMaterialInput): void;
}>();

const dialogOpen = ref(props.open);

// Watch for changes in the open prop
watch(() => props.open, (newValue) => {
  console.log('MaterialEditWrapper - open prop changed to:', newValue);
  dialogOpen.value = newValue;
});

// Watch for changes in the dialogOpen ref
watch(() => dialogOpen.value, (newValue) => {
  console.log('MaterialEditWrapper - dialogOpen changed to:', newValue);
  emit('update:open', newValue);
});

function handleUpdateMaterial(materialData: NewMaterialInput) {
  console.log('MaterialEditWrapper - handleUpdateMaterial called');
  
  // Basic validation to ensure required fields are present
  if (!materialData || 
      !materialData.name || 
      !materialData.category || 
      !materialData.supplier || 
      materialData.quantity === undefined || 
      !materialData.image) {
    console.error('Invalid material data provided:', materialData);
    return;
  }
  
  emit('update-material', materialData);
  // Close the dialog immediately on successful material update
  dialogOpen.value = false;
}
</script>

<template>
  <EditMaterialDialog
    v-model="dialogOpen"
    :material-data="materialData"
    :categories="categories"
    :suppliers="suppliers"
    :default-image-url="defaultImageUrl"
    @update-material="handleUpdateMaterial"
  />
</template>