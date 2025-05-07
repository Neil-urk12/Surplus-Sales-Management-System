<script setup lang="ts">
import { ref, computed, watch, onBeforeUnmount } from 'vue';
import { useQuasar } from 'quasar';
import type { PropType } from 'vue';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';
import type { MaterialRow } from 'src/stores/materials';
import type { MaterialCategory, MaterialSupplier } from 'src/types/materials';

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
    type: Object as PropType<MaterialRow | null>,
    default: null
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
  },
  statuses: {
    type: Array as PropType<string[]>,
    default: () => []
  }
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'hide'): void
  (e: 'update', updatedMaterial: MaterialRow): void
}>();

const $q = useQuasar();

// Constants for file validation
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif'] as const;
const MAX_DIMENSION = 4096; // Maximum image dimension in pixels
type AllowedMimeType = typeof ALLOWED_TYPES[number];

// --- Local State --- 
// Use a local ref to store editable data
const localMaterialData = ref<MaterialRow>({
  id: 0,
  name: '',
  category: '' as MaterialCategory,
  supplier: '' as MaterialSupplier,
  quantity: 0,
  status: 'Out of Stock',
  image: props.defaultImageUrl
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
  get: () => localMaterialData.value.name,
  set: (value: string) => {
    if (value) {
      localMaterialData.value.name = value.charAt(0).toUpperCase() + value.slice(1);
    } else {
      localMaterialData.value.name = value;
    }
  }
});

// --- Watchers --- 
watch(() => localMaterialData.value.quantity, (newQuantity) => {
  if (newQuantity === 0) {
    localMaterialData.value.status = 'Out of Stock';
  } else if (newQuantity <= 10) {
    localMaterialData.value.status = 'Low Stock';
  } else if (newQuantity <= 50) {
    localMaterialData.value.status = 'In Stock';
  } else {
    localMaterialData.value.status = 'Available';
  }
});

watch(() => localMaterialData.value.image, (newUrl: string) => {
  if (!newUrl || newUrl === props.defaultImageUrl) {
    previewUrl.value = props.defaultImageUrl;
    imageUrlValid.value = true;
    return;
  }
  if (newUrl.startsWith('data:image/')) {
    // Only do basic validation here, not full sanitization
    const validationResult = validateAndSanitizeBase64Image(newUrl);
    if (validationResult.isValid) {
      // Don't update the data yet, just show the preview
      previewUrl.value = newUrl;
      imageUrlValid.value = true;
    } else {
      $q.notify({ color: 'negative', message: validationResult.error || 'Invalid image data', position: 'top' });
      previewUrl.value = props.defaultImageUrl;
      imageUrlValid.value = false;
    }
    return;
  }

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

// Watch the incoming prop to update local state when the dialog opens/changes material
watch(() => props.materialData, (newMaterial) => {
  if (newMaterial) {
    // Reset local state based on the new material data
    const initialImage = newMaterial.image || props.defaultImageUrl;
    localMaterialData.value = {
      id: newMaterial.id,
      name: newMaterial.name,
      category: newMaterial.category,
      supplier: newMaterial.supplier,
      quantity: newMaterial.quantity,
      status: newMaterial.status,
      image: initialImage
    };

    // Directly validate the initial image
    if (initialImage.startsWith('data:image/')) {
      const validationResult = validateAndSanitizeBase64Image(initialImage);
      if (validationResult.isValid) {
        previewUrl.value = initialImage;
        imageUrlValid.value = true;
      } else {
        previewUrl.value = props.defaultImageUrl;
        imageUrlValid.value = false;
        $q.notify({
          type: 'negative',
          message: validationResult.error || 'Invalid image data',
          position: 'top'
        });
      }
    } else if (initialImage !== props.defaultImageUrl) {
      void validateImageUrl(initialImage);
    } else {
      previewUrl.value = props.defaultImageUrl;
      imageUrlValid.value = true;
    }
  } else {
    // Reset if materialData becomes null
    resetForm();
  }
}, { immediate: true }); // Immediate run on component mount

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
  // Reset local data to empty/default state
  localMaterialData.value = {
    id: 0,
    name: '',
    category: '' as MaterialCategory,
    supplier: '' as MaterialSupplier,
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

function onUpdate() {
  // Validate required fields
  if (!localMaterialData.value.name || !localMaterialData.value.category || !localMaterialData.value.supplier) {
    operationNotifications.validation.error('Please fill in all required fields');
    return;
  }

  // Validate image
  if (!imageUrlValid.value) {
    operationNotifications.validation.error('Please provide a valid image');
    return;
  }

  // Ensure image is set, fallback to default if necessary
  if (!localMaterialData.value.image || localMaterialData.value.image === '') {
    localMaterialData.value.image = props.defaultImageUrl;
  }

  // Full sanitization right before sending to backend
  if (localMaterialData.value.image.startsWith('data:image/')) {
    const sanitizationResult = validateAndSanitizeBase64Image(localMaterialData.value.image);
    if (!sanitizationResult.isValid) {
      operationNotifications.validation.error(sanitizationResult.error || 'Invalid image data');
      return;
    }
    localMaterialData.value.image = sanitizationResult.sanitizedData!;
  }

  emit('update', localMaterialData.value);
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
          localMaterialData.value.image = base64ValidationResult.sanitizedData!;
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
  // Reset preview and image data based on the potentially existing image in localMaterialData or default
  const initialImage = props.materialData?.image || props.defaultImageUrl;
  localMaterialData.value.image = initialImage;
  previewUrl.value = initialImage;
  imageUrlValid.value = true; // Assume original/default is valid initially
  if (fileInput.value) {
    fileInput.value.value = '';
  }
  isUploadingImage.value = false;
}

async function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files[0]) {
    const file = input.files[0];
    if (!ALLOWED_TYPES.includes(file.type as AllowedMimeType)) { clearImageInput(); return; }
    if (file.size > MAX_FILE_SIZE) { clearImageInput(); return; }
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
        <div class="text-h6">Edit Material</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>
      <q-card-section>
        <q-form @submit.prevent="onUpdate" class="q-gutter-sm">
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
                v-model="localMaterialData.category"
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
                v-model="localMaterialData.supplier"
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
            <div class="col-12 col-sm-6">
              <q-input
                v-model.number="localMaterialData.quantity"
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
            <div class="col-12 col-sm-6">
              <q-select
                v-model="localMaterialData.status"
                :options="statuses"
                label="Status"
                dense
                outlined
                required
                emit-value
                map-options
                lazy-rules
                :rules="[val => !!val || 'Status is required']"
              >
                <template v-slot:prepend>
                  <q-icon name="info" />
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
                      :alt="localMaterialData.name || 'Preview image'"
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
          label="Update Material"
          @click="onUpdate"
          :disable="disable || !localMaterialData.name || !localMaterialData.category || !localMaterialData.supplier || localMaterialData.quantity < 0 || !imageUrlValid || isUploadingImage"
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