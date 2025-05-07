<script setup lang="ts">
import { ref, computed, watch, onBeforeUnmount } from 'vue';
import { useQuasar } from 'quasar';
import type { PropType } from 'vue';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';
import type { NewMaterialInput } from 'src/stores/materials';
import type { MaterialCategoryInput, MaterialSupplierInput } from 'src/types/materials';

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true,
  },
  disable: {
    type: Boolean,
    default: false
  },
  materialData: {
    type: Object as PropType<NewMaterialInput>,
    default: () => ({
      name: '',
      category: '',
      supplier: '',
      quantity: 0,
      status: 'Out of Stock',
      image: 'https://loremflickr.com/600/400/material'
    })
  },
  defaultImageUrl: {
    type: String,
    default: 'https://loremflickr.com/600/400/material'
  },
  categories: {
    type: Array as PropType<string[]>,
    default: () => []
  },
  suppliers: {
    type: Array as PropType<string[]>,
    default: () => []
  }
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'update:materialData', value: NewMaterialInput): void
  (e: 'hide'): void
  (e: 'add', newMaterial: NewMaterialInput): void
}>();

const $q = useQuasar();

// Constants for file validation
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif'] as const;
const MAX_DIMENSION = 4096; // Maximum image dimension in pixels
type AllowedMimeType = typeof ALLOWED_TYPES[number];

// --- State --- 
const newMaterial = computed({
  get: () => props.materialData,
  set: (value) => emit('update:materialData', value)
});

const imageUrlValid = ref(true);
const validatingImage = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);
const previewUrl = ref<string>(props.defaultImageUrl);
const isUploadingImage = ref(false);
const isDragging = ref(false);
let currentAbortController: AbortController | null = null;

// --- Computed --- 
const capitalizedName = computed({
  get: () => newMaterial.value.name,
  set: (value: string) => {
    if (value) {
      newMaterial.value.name = value.charAt(0).toUpperCase() + value.slice(1);
    } else {
      newMaterial.value.name = value;
    }
  }
});

// --- Watchers --- 
watch(() => newMaterial.value.quantity, (newQuantity) => {
  if (newQuantity === 0) {
    newMaterial.value.status = 'Out of Stock';
  } else if (newQuantity <= 10) {
    newMaterial.value.status = 'Low Stock';
  } else if (newQuantity <= 50) {
    newMaterial.value.status = 'In Stock';
  } else {
    newMaterial.value.status = 'Available';
  }
});

watch(() => newMaterial.value.image, (newUrl: string) => {
  if (!newUrl || newUrl === props.defaultImageUrl) {
    previewUrl.value = props.defaultImageUrl;
    imageUrlValid.value = true;
    return;
  }
  if (newUrl.startsWith('data:image/')) {
    const validationResult = validateAndSanitizeBase64Image(newUrl);
    if (validationResult.isValid) {
      newMaterial.value.image = validationResult.sanitizedData!;
      previewUrl.value = validationResult.sanitizedData!;
      imageUrlValid.value = true;
    } else {
      $q.notify({
        color: 'negative',
        message: validationResult.error || 'Invalid image data',
        position: 'top',
      });
      previewUrl.value = props.defaultImageUrl;
      imageUrlValid.value = false;
    }
    return;
  }

  // Validate HTTP/HTTPS URLs
  void validateImageUrl(newUrl).then(isValid => {
    if (isValid) {
      previewUrl.value = newUrl;
    } else {
      previewUrl.value = props.defaultImageUrl;
    }
  }).catch(error => {
    console.error('Error in image URL watcher:', error);
    previewUrl.value = props.defaultImageUrl;
    imageUrlValid.value = false;
  });
});

// Reset form when modal opens
watch(() => props.modelValue, (newValue) => {
  if (newValue) {
    resetForm();
  }
});

// --- Lifecycle Hooks ---
onBeforeUnmount(() => {
  // Abort any ongoing image URL validation when the component is unmounted
  if (currentAbortController) {
    currentAbortController.abort();
    currentAbortController = null;
  }
});

// --- Functions --- 
function resetForm() {
  newMaterial.value = {
    name: '',
    category: '' as MaterialCategoryInput,
    supplier: '' as MaterialSupplierInput,
    quantity: 0,
    status: 'Out of Stock',
    image: props.defaultImageUrl
  };
  clearImageInput();
}

function onHide() {
  emit('update:modelValue', false);
  emit('hide');
}

function onAdd() {
  // Validate required fields
  if (!newMaterial.value.name || !newMaterial.value.category || !newMaterial.value.supplier) {
    operationNotifications.validation.error('Please fill in all required fields');
    return;
  }

  // Validate image
  if (!imageUrlValid.value) {
    operationNotifications.validation.error('Please provide a valid image');
    return;
  }

  // Ensure image is set, fallback to default if necessary
  if (!newMaterial.value.image || newMaterial.value.image === '') {
    // Create a new object to avoid mutating the props directly
    const updatedMaterial = { ...newMaterial.value, image: props.defaultImageUrl };
    emit('update:materialData', updatedMaterial);
    emit('add', updatedMaterial);
    return;
  }

  // Full sanitization right before sending to backend
  if (newMaterial.value.image.startsWith('data:image/')) {
    const sanitizationResult = validateAndSanitizeBase64Image(newMaterial.value.image);
    if (!sanitizationResult.isValid) {
      operationNotifications.validation.error(sanitizationResult.error || 'Invalid image data');
      return;
    }
    // Create a new object to avoid mutating the props directly
    const updatedMaterial = { ...newMaterial.value, image: sanitizationResult.sanitizedData! };
    emit('update:materialData', updatedMaterial);
    emit('add', updatedMaterial);
    return;
  }
  
  emit('add', newMaterial.value);
}

async function validateImageUrl(url: string): Promise<boolean> {
  if (!url || !url.startsWith('http')) {
    imageUrlValid.value = false;
    return false;
  }

  validatingImage.value = true;

  if (currentAbortController) {
    currentAbortController.abort();
  }
  currentAbortController = new AbortController();
  const signal = currentAbortController.signal;

  try {
    const result = await new Promise<boolean>((resolve) => {
      const img = new Image();
      const cleanup = () => {
        img.onload = null;
        img.onerror = null;
        if (currentAbortController?.signal === signal) {
          currentAbortController = null;
        }
      };

      signal.addEventListener('abort', () => { cleanup(); resolve(false); });
      img.onload = () => { cleanup(); imageUrlValid.value = true; validatingImage.value = false; resolve(true); };
      img.onerror = () => {
        cleanup();
        imageUrlValid.value = false;
        validatingImage.value = false;
        resolve(false);
      };

      const timeoutId = setTimeout(() => {
        if (!signal.aborted) {
          currentAbortController?.abort();
          imageUrlValid.value = false;
          validatingImage.value = false;
          resolve(false);
        }
      }, 5000);

      signal.addEventListener('abort', () => { clearTimeout(timeoutId); });
      img.src = url;
    });
    return result;
  } catch (error) {
    console.error('Error validating image URL:', error);
    imageUrlValid.value = false;
    validatingImage.value = false;
    return false;
  } finally {
    if (validatingImage.value) {
      validatingImage.value = false;
    }
  }
}

async function validateFile(file: File): Promise<{ isValid: boolean; error?: string }> {
  try {
    console.log('Starting file validation:', { name: file.name, type: file.type, size: file.size });

    if (!file) return { isValid: false, error: 'No file provided.' };
    if (file.size > MAX_FILE_SIZE) {
      const sizeMB = (file.size / (1024 * 1024)).toFixed(2);
      return { isValid: false, error: `File size (${sizeMB}MB) exceeds 5MB limit.` };
    }

    const validMimeTypes = {
      'image/jpeg': [0xFF, 0xD8, 0xFF],
      'image/png': [0x89, 0x50, 0x4E, 0x47],
      'image/gif': [0x47, 0x49, 0x46, 0x38]
    };

    if (!Object.keys(validMimeTypes).includes(file.type)) {
      return { isValid: false, error: `Invalid file type: ${file.type}.` };
    }

    const arrayBuffer = await file.slice(0, 4).arrayBuffer();
    const bytes = new Uint8Array(arrayBuffer);
    const expectedSignature = validMimeTypes[file.type as keyof typeof validMimeTypes];
    const isValidSignature = expectedSignature.every((byte, i) => byte === bytes[i]);
    if (!isValidSignature) {
      return { isValid: false, error: 'File content does not match extension.' };
    }

    const dimensionValidation = await validateImageDimensions(file);
    if (!dimensionValidation.isValid) {
      return dimensionValidation;
    }

    return { isValid: true };
  } catch (error) {
    console.error('Unexpected error during file validation:', error);
    return { isValid: false, error: 'Validation error occurred.' };
  }
}

function validateImageDimensions(file: File): Promise<{ isValid: boolean; error?: string }> {
  return new Promise((resolve) => {
    const img = new Image();
    const objectUrl = URL.createObjectURL(file);
    const cleanup = () => URL.revokeObjectURL(objectUrl);
    let timeoutId: ReturnType<typeof setTimeout> | null = null;

    const resolveClean = (result: { isValid: boolean; error?: string }) => {
      if (timeoutId) clearTimeout(timeoutId);
      cleanup();
      resolve(result);
    };

    img.onload = () => {
      if (img.width > MAX_DIMENSION || img.height > MAX_DIMENSION) {
        resolveClean({ isValid: false, error: `Image dimensions exceed ${MAX_DIMENSION}px.` });
      } else if (img.width === 0 || img.height === 0) {
        resolveClean({ isValid: false, error: 'Invalid image dimensions.' });
      } else {
        resolveClean({ isValid: true });
      }
    };
    img.onerror = () => resolveClean({ isValid: false, error: 'Error loading image file.' });

    // Add timeout
    timeoutId = setTimeout(() => {
      console.error('Dimension validation timed out');
      resolveClean({ isValid: false, error: 'Image validation timed out. Please try again.' });
    }, 10000); // 10 second timeout

    img.src = objectUrl;
  });
}

async function handleFile(file: File) {
  try {
    isUploadingImage.value = true;
    previewUrl.value = ''; // Clear preview during processing

    const validation = await validateFile(file);
    if (!validation.isValid) {
      $q.notify({ type: 'negative', message: validation.error || 'Invalid file', position: 'top' });
      clearImageInput();
      return;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
      if (e.target?.result) {
        const base64String = e.target.result as string;
        const base64ValidationResult = validateAndSanitizeBase64Image(base64String);
        if (!base64ValidationResult.isValid) {
          $q.notify({ type: 'negative', message: base64ValidationResult.error || 'Invalid image data', position: 'top' });
          clearImageInput();
        } else {
          newMaterial.value.image = base64ValidationResult.sanitizedData!;
          imageUrlValid.value = true; // Watcher will update previewUrl
          $q.notify({ type: 'positive', message: 'Image uploaded successfully', position: 'top', timeout: 2000 });
        }
      } else {
        clearImageInput();
      }
      isUploadingImage.value = false;
    };
    reader.onerror = () => {
      $q.notify({ type: 'negative', message: 'Error reading file.', position: 'top' });
      clearImageInput();
      isUploadingImage.value = false;
    };
    reader.readAsDataURL(file);
  } catch (error) {
    console.error('Error in handleFile:', error);
    $q.notify({ type: 'negative', message: 'An unexpected error occurred.', position: 'top' });
    clearImageInput();
    isUploadingImage.value = false;
  }
}

function removeImage(event?: Event) {
  if (event) event.stopPropagation();
  clearImageInput();
}

function clearImageInput() {
  if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(previewUrl.value);
  }
  previewUrl.value = props.defaultImageUrl;
  newMaterial.value.image = props.defaultImageUrl;
  imageUrlValid.value = true;
  if (fileInput.value) {
    fileInput.value.value = '';
  }
  isUploadingImage.value = false;
}

async function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files[0]) {
    const file = input.files[0];
    if (!ALLOWED_TYPES.includes(file.type as AllowedMimeType)) {
      $q.notify({
        type: 'negative',
        message: `Invalid file type: ${file.type}. Allowed types are: JPEG, PNG, and GIF`,
        position: 'top',
        timeout: 3000
      });
      return;
    }
    await handleFile(file);
  }
  if (input) {
    input.value = '';
  }
}

function triggerFileInput() {
  fileInput.value?.click();
}

// --- Drag & Drop --- 
function handleDragLeave(event: DragEvent) {
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
  const x = event.clientX;
  const y = event.clientY;
  if (x <= rect.left || x >= rect.right || y <= rect.top || y >= rect.bottom) {
    isDragging.value = false;
  }
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
  <q-dialog :modelValue="modelValue" persistent @hide="onHide" @update:modelValue="$emit('update:modelValue', $event)">
    <q-card style="min-width: 400px; max-width: 95vw">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">New Material</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>
      <q-card-section>
        <q-form @submit.prevent="onAdd" class="q-gutter-sm">
          <q-input
            v-model="capitalizedName"
            label="Material Name"
            dense
            outlined
            required
            lazy-rules
            :rules="[val => !!val || 'Name is required']"
          >
            <template v-slot:prepend>
              <q-icon name="inventory_2" />
            </template>
          </q-input>

          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <q-select
                v-model="newMaterial.category"
                :options="categories"
                label="Category"
                dense
                outlined
                required
                emit-value
                map-options
                placeholder="Select a category"
                lazy-rules
                :rules="[val => !!val || 'Category is required']"
              >
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
                v-model="newMaterial.supplier"
                :options="suppliers"
                label="Supplier"
                dense
                outlined
                required
                emit-value
                map-options
                placeholder="Select a supplier"
                lazy-rules
                :rules="[val => !!val || 'Supplier is required']"
              >
                <template v-slot:prepend>
                  <q-icon name="business" />
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
            <div class="col-12">
              <q-input
                v-model.number="newMaterial.quantity"
                type="number"
                min="0"
                label="Quantity"
                dense
                outlined
                required
                lazy-rules
                :rules="[val => val !== null && val !== undefined && val >= 0 || 'Quantity must be positive']"
              >
                <template v-slot:prepend>
                  <q-icon name="numbers" />
                </template>
              </q-input>
            </div>
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-12">
              <q-input
                v-model="newMaterial.status"
                label="Status"
                dense
                outlined
                readonly
              >
                <template v-slot:prepend>
                  <q-icon name="info" />
                </template>
              </q-input>
            </div>
          </div>

          <!-- Enhanced Image Input Section -->
          <div class="row q-col-gutter-sm">
            <div class="col-12">
              <div
                class="upload-container q-pa-md"
                :class="{ 'dragging': isDragging }"
                @dragenter.prevent="isDragging = true"
                @dragover.prevent="isDragging = true"
                @dragleave.prevent="handleDragLeave"
                @drop.prevent="handleDrop"
                @dragend.prevent="isDragging = false"
                @click="triggerFileInput"
              >
                <input
                  type="file"
                  ref="fileInput"
                  accept="image/jpeg,image/png,image/gif"
                  class="hidden"
                  @change="handleFileSelect"
                >
                
                <div v-if="isUploadingImage" class="text-center">
                  <q-spinner-dots color="primary" size="40px" />
                  <div class="text-body1 q-mt-sm">
                    Processing image...
                  </div>
                </div>
                
                <div v-else-if="!previewUrl || previewUrl === defaultImageUrl" class="text-center">
                  <q-icon name="cloud_upload" size="48px" color="primary" />
                  <div class="text-body1 q-mt-sm">
                    Drag and drop an image here or click to select
                  </div>
                  <div class="text-caption text-grey q-mt-sm">
                    Supported formats: JPG, PNG, GIF
                    <br>
                    Maximum file size: 5MB
                    <br>
                    Maximum dimensions: 4096x4096px
                  </div>
                </div>
                
                <div v-else class="row items-center justify-center">
                  <div class="col-8 text-center">
                    <img
                      :src="previewUrl"
                      class="preview-image"
                      :alt="newMaterial.name || 'Preview image'"
                      @error="previewUrl = defaultImageUrl"
                    />
                  </div>
                  <div class="col-4 text-center">
                    <q-btn
                      flat
                      round
                      color="negative"
                      icon="close"
                      @click.stop="removeImage"
                    >
                      <q-tooltip>Remove Image</q-tooltip>
                    </q-btn>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </q-form>
      </q-card-section>
      <q-card-actions align="right" class="q-pa-md">
        <q-btn flat label="Cancel" @click="onHide" />
        <q-btn
          unelevated
          color="primary"
          label="Add Material"
          @click="onAdd"
          :disable="disable || !newMaterial.name || !newMaterial.category || !newMaterial.supplier || newMaterial.quantity < 0 || !imageUrlValid || isUploadingImage"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
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