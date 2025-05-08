<script setup lang="ts">
import { ref, computed, watch, onBeforeUnmount } from 'vue';
import type { PropType } from 'vue';
import type { NewMaterialInput, MaterialCategoryInput, MaterialSupplierInput, MaterialStatus } from 'src/types/materials';
import { validateFileUpload, MAX_IMAGE_SIZE, DEFAULT_MAX_DIMENSION } from '../../utils/imageValidation';

const props = defineProps({
  materialData: {
    type: Object as PropType<NewMaterialInput>,
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
  isProcessing: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits<{
  (e: 'update:materialData', value: NewMaterialInput): void;
  (e: 'submit'): void;
  (e: 'cancel'): void;
  (e: 'validation-error', message: string): void;
}>();

// Constants for file validation
const MAX_FILE_SIZE = MAX_IMAGE_SIZE;
const MAX_DIMENSION = DEFAULT_MAX_DIMENSION;

// --- State ---
const imageUrlValid = ref(true);
const fileInput = ref<HTMLInputElement | null>(null);
const previewUrl = ref<string>(props.materialData.image || props.defaultImageUrl);
const isUploadingImage = ref(false);
const isDragging = ref(false);
const localQuantity = ref(props.materialData.quantity ?? 0);
const isPreviewBlobUrl = ref(false);
let currentAbortController: AbortController | null = null;

// For debouncing
let dragLeaveTimer: number | null = null;
let quantityUpdateTimer: number | null = null;

// Watch for changes in the localQuantity
watch(localQuantity, (newVal) => {
  // Update the materialData quantity
  const numValue = typeof newVal === 'string' ? (Number.isFinite(Number(newVal)) ? Number(newVal) : 0) : newVal;
  const updatedData = { ...props.materialData, quantity: numValue };
  
  // Update status based on quantity
  if (numValue === 0) {
    updatedData.status = 'Out of Stock';
  } else if (numValue <= 10) {
    updatedData.status = 'Low Stock';
  } else {
    updatedData.status = 'In Stock';
  }
  
  emit('update:materialData', updatedData);
});

// Watch for external changes to materialData.quantity
watch(() => props.materialData.quantity, (newVal) => {
  localQuantity.value = newVal ?? 0;
}, { immediate: true });

// --- Computed --- 
const capitalizedName = computed({
  get: () => props.materialData.name,
  set: (value: string) => {
    if (value) {
      const updatedData = { 
        ...props.materialData, 
        name: value.charAt(0).toUpperCase() + value.slice(1) 
      };
      emit('update:materialData', updatedData);
    } else {
      const updatedData = { ...props.materialData, name: value };
      emit('update:materialData', updatedData);
    }
  }
});

// Computed properties for v-model binding
const category = computed({
  get: () => props.materialData.category,
  set: (value: MaterialCategoryInput) => updateField('category', value)
});

const supplier = computed({
  get: () => props.materialData.supplier,
  set: (value: MaterialSupplierInput) => updateField('supplier', value)
});

const status = computed({
  get: () => props.materialData.status,
  set: (value: MaterialStatus) => updateField('status', value)
});

// Computed for button disable condition
const isSubmitDisabled = computed(() => {
  return !props.materialData.name || 
         !props.materialData.category || 
         !props.materialData.supplier || 
         props.materialData.quantity < 0 || 
         !imageUrlValid.value || 
         isUploadingImage.value;
});

// Watch for changes in the material data
watch(() => props.materialData, () => {
  previewUrl.value = props.materialData.image || props.defaultImageUrl;
}, { deep: true });

// --- Lifecycle Hooks ---
onBeforeUnmount(() => {
  if (currentAbortController) {
    currentAbortController.abort();
    currentAbortController = null;
  }
  // Cleanup any object URLs to prevent memory leaks
  if (isPreviewBlobUrl.value && previewUrl.value && previewUrl.value.startsWith('blob:')) {
    try {
      URL.revokeObjectURL(previewUrl.value);
    } catch (error) {
      console.error('Error revoking object URL:', error);
    }
  }
  
  // Clear any pending timers
  if (dragLeaveTimer !== null) {
    clearTimeout(dragLeaveTimer);
    dragLeaveTimer = null;
  }
  
  if (quantityUpdateTimer !== null) {
    clearTimeout(quantityUpdateTimer);
    quantityUpdateTimer = null;
  }
});

// --- Functions ---
function handleQuantityInput(event: Event) {
  // Get the value directly from the input element
  const input = event.target as HTMLInputElement;
  const value = input.value;
  
  // Clear any existing debounce timer
  if (quantityUpdateTimer !== null) {
    clearTimeout(quantityUpdateTimer);
  }
  
  // Debounce the update
  quantityUpdateTimer = window.setTimeout(() => {
    // Update local quantity ref
    localQuantity.value = value === '' ? 0 : Number(value);
    quantityUpdateTimer = null;
  }, 300);
}

function handleQuantityBlur(event: Event) {
  // Get the value directly from the input element
  const input = event.target as HTMLInputElement;
  const value = input.value;
  
  // Clear any pending debounce timer
  if (quantityUpdateTimer !== null) {
    clearTimeout(quantityUpdateTimer);
    quantityUpdateTimer = null;
  }
  
  // Ensure the quantity is updated in the materialData
  const numValue = value === '' ? 0 : (Number.isFinite(Number(value)) ? Number(value) : 0);
  
  // Find the appropriate status
  let newStatus: MaterialStatus;
  if (numValue === 0) {
    newStatus = 'Out of Stock';
  } else if (numValue <= 10) {
    newStatus = 'Low Stock';
  } else {
    newStatus = 'In Stock';
  }
  
  // Update both quantity and status
  const updatedData = { 
    ...props.materialData, 
    quantity: numValue,
    status: newStatus
  };
  
  // Emit the update
  emit('update:materialData', updatedData);
}

function clearImageInput() {
  if (isPreviewBlobUrl.value && previewUrl.value && previewUrl.value.startsWith('blob:')) {
    try {
      URL.revokeObjectURL(previewUrl.value);
    } catch (error) {
      console.error('Error revoking object URL:', error);
    }
  }
  previewUrl.value = props.defaultImageUrl;
  isPreviewBlobUrl.value = false;
  
  const updatedData = { ...props.materialData, image: props.defaultImageUrl };
  emit('update:materialData', updatedData);
  
  imageUrlValid.value = true;
  if (fileInput.value) {
    fileInput.value.value = '';
  }
  isUploadingImage.value = false;
}

function handleSubmit() {
  if (!props.materialData.name) {
    emit('validation-error', 'Please enter a material name');
    return;
  }

  if (!props.materialData.category || !props.materialData.supplier) {
    emit('validation-error', 'Please fill in all required fields');
    return;
  }

  if (!imageUrlValid.value) {
    emit('validation-error', 'Please provide a valid image');
    return;
  }

  if (!props.materialData.image || props.materialData.image === '') {
    const updatedData = { ...props.materialData, image: props.defaultImageUrl };
    emit('update:materialData', updatedData);
  }

  emit('submit');
}

function handleCancel() {
  emit('cancel');
}

function updateField<T extends keyof NewMaterialInput>(field: T, value: NewMaterialInput[T]) {
  const updatedData = { ...props.materialData, [field]: value };
  emit('update:materialData', updatedData);
}

// --- Image Handling --- 
function validateAndProcessImage(file: File): Promise<{success: boolean, result?: string, error?: string}> {
  return validateFileUpload(file, {
    maxFileSize: MAX_FILE_SIZE,
    maxDimension: MAX_DIMENSION
  })
    .then(result => {
      if (result.isValid && result.sanitizedData) {
        return {
          success: true,
          result: result.sanitizedData
        };
      } else {
        return {
          success: false,
          error: result.error || 'Invalid image'
        };
      }
    })
    .catch(error => {
      console.error('Error processing image:', error);
      return {
        success: false,
        error: 'Error processing image'
      };
    });
}

function handleFile(file: File) {
  isUploadingImage.value = true;

  void validateAndProcessImage(file)
    .then(response => {
      if (response.success && response.result) {
        isPreviewBlobUrl.value = response.result.startsWith('blob:');
        const updatedData = { 
          ...props.materialData, 
          image: response.result 
        };
        emit('update:materialData', updatedData);
        previewUrl.value = response.result;
        imageUrlValid.value = true;
      } else {
        emit('validation-error', response.error || 'Invalid image');
        clearImageInput();
      }
    })
    .finally(() => {
      isUploadingImage.value = false;
    });
}

function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files[0]) {
    const file = input.files[0];
    void handleFile(file);
  }
  if (input) {
    input.value = '';
  }
}

function triggerFileInput() {
  fileInput.value?.click();
}

function removeImage() {
  clearImageInput();
}

// --- Drag & Drop --- 
function handleDragLeave(event: DragEvent) {
  // Improve drag leave detection by checking related target
  const related = event.relatedTarget as Node | null;
  const container = event.currentTarget as HTMLElement;
  
  // If the related target is a child of the container, don't set isDragging to false
  if (related && container.contains(related)) {
    return;
  }
  
  // Clear any existing timer
  if (dragLeaveTimer !== null) {
    clearTimeout(dragLeaveTimer);
  }
  
  // Set a new timer for debouncing
  dragLeaveTimer = window.setTimeout(() => {
    const rect = container.getBoundingClientRect();
    const x = event.clientX;
    const y = event.clientY;
    if (x <= rect.left || x >= rect.right || y <= rect.top || y >= rect.bottom) {
      isDragging.value = false;
    }
    dragLeaveTimer = null;
  }, 100); // 100ms debounce time
}

function handleDrop(event: DragEvent) {
  event.preventDefault();
  isDragging.value = false;
  if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
    const file = event.dataTransfer.files[0];
    void handleFile(file);
  }
}
</script>

<template>
  <q-form @submit.prevent="handleSubmit" class="q-gutter-sm">
    <q-input v-model="capitalizedName" label="Material Name" dense outlined required lazy-rules
      :rules="[val => !!val || 'Name is required']">
      <template v-slot:prepend>
        <q-icon name="inventory_2" />
      </template>
    </q-input>

    <div class="row q-col-gutter-sm">
      <div class="col-12 col-sm-6">
        <q-select 
          v-model="category" 
          :options="categories" 
          label="Category" 
          dense 
          outlined
          required 
          emit-value 
          map-options 
          placeholder="Select a category" 
          lazy-rules
          :rules="[val => !!val || 'Category is required']">
          <template v-slot:prepend>
            <q-icon name="category" />
          </template>
          <template v-slot:no-option>
            <q-item>
              <q-item-section class="text-grey">
                No results
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>

      <div class="col-12 col-sm-6">
        <q-select 
          v-model="supplier" 
          :options="suppliers" 
          label="Supplier" 
          dense 
          outlined
          required 
          emit-value 
          map-options 
          placeholder="Select a supplier" 
          lazy-rules
          :rules="[val => !!val || 'Supplier is required']">
          <template v-slot:prepend>
            <q-icon name="local_shipping" />
          </template>
          <template v-slot:no-option>
            <q-item>
              <q-item-section class="text-grey">
                No results
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>
    </div>

    <div class="row q-col-gutter-sm">
      <div class="col-12 col-sm-6">
        <q-input 
          v-model="localQuantity" 
          type="number" 
          min="0" 
          label="Quantity" 
          dense
          outlined 
          required
          @input="handleQuantityInput"
          @blur="handleQuantityBlur"
          :rules="[(val: any) => val !== null && val !== undefined && val >= 0 || 'Quantity must be positive']">
          <template v-slot:prepend>
            <q-icon name="numbers" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-sm-6">
        <q-input 
          v-model="status" 
          label="Status" 
          dense 
          outlined 
          readonly>
          <template v-slot:prepend>
            <q-icon name="info" />
          </template>
        </q-input>
      </div>
    </div>

    <div class="row q-col-gutter-sm">
      <div class="col-12">
        <div class="upload-container" :class="{ dragging: isDragging }"
          @dragover.prevent="isDragging = true" @dragleave="handleDragLeave"
          @drop.prevent="handleDrop" @click="triggerFileInput">
          <input ref="fileInput" type="file" accept="image/*" class="hidden"
            @change="handleFileSelect" />
          <div v-if="isUploadingImage" class="text-center">
            <q-spinner-dots color="primary" size="40px" />
            <div class="text-body1 q-mt-sm">Processing image...</div>
          </div>
          <div v-else-if="!previewUrl || previewUrl === defaultImageUrl" class="text-center">
            <q-icon name="cloud_upload" size="48px" color="primary" />
            <div class="text-body1 q-mt-sm">Drag/drop or click to select image</div>
            <div class="text-caption text-grey q-mt-sm">JPG, PNG, GIF | Max 5MB | Max {{ MAX_DIMENSION }}px
            </div>
          </div>
          <div v-else class="row items-center justify-center">
            <div class="col-8 text-center">
              <img :src="previewUrl" class="preview-image"
                :alt="props.materialData.name || 'Preview image'"
                @error="previewUrl = defaultImageUrl" />
            </div>
            <div class="col-4 text-center">
              <q-btn flat round color="negative" icon="close" @click.stop="removeImage">
                <q-tooltip>Remove Image</q-tooltip>
              </q-btn>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="row justify-end q-gutter-sm q-mt-md">
      <q-btn flat label="Cancel" @click="handleCancel" />
      <slot name="submitButton" 
            :disabled="isSubmitDisabled"
            :loading="isProcessing">
        <q-btn unelevated color="primary" label="Submit" type="submit" :loading="isProcessing"
          :disable="isSubmitDisabled" />
      </slot>
    </div>
  </q-form>
</template>

<style scoped lang="sass">
.upload-container
  position: relative
  border: 2px dashed #ccc
  border-radius: 8px
  cursor: pointer
  transition: all 0.3s ease
  min-height: 200px
  display: flex
  align-items: center
  justify-content: center
  background-color: transparent
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05)

  &:hover
    border-color: var(--q-primary)
    background-color: rgba(var(--q-primary-rgb), 0.04)
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1)
    transform: translateY(-1px)

  &.dragging
    border-color: var(--q-primary)
    background-color: rgba(var(--q-primary-rgb), 0.08)
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15)
    transform: translateY(-2px)

  .q-spinner-dots
    margin: 0 auto

.preview-image
  width: 100%
  max-height: 180px
  object-fit: contain
  border-radius: 4px
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1)
  transition: all 0.3s ease

  &:hover
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15)

.hidden
  display: none
</style> 