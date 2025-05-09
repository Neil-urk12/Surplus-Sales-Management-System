<script setup lang="ts">
import { ref, watch, defineProps, defineEmits } from 'vue';
import { useQuasar } from 'quasar';
import { validateFileUpload, MAX_IMAGE_SIZE, DEFAULT_MAX_DIMENSION } from 'src/utils/imageValidation';

const props = defineProps({
  modelValue: {
    type: String,
    required: true
  },
  defaultImageUrl: {
    type: String,
    default: 'https://loremflickr.com/600/400/material'
  }
});

const emit = defineEmits(['update:modelValue']);

const $q = useQuasar();
const fileInput = ref<HTMLInputElement | null>(null);
const previewUrl = ref(props.modelValue || props.defaultImageUrl);
const isUploadingImage = ref(false);
const isDragging = ref(false);

// Constants for validation
const MAX_FILE_SIZE = MAX_IMAGE_SIZE;
const MAX_DIMENSION = DEFAULT_MAX_DIMENSION;

// Watch for prop changes
watch(() => props.modelValue, (newValue) => {
  if (newValue !== previewUrl.value) {
    previewUrl.value = newValue || props.defaultImageUrl;
  }
});

// Enhanced drag event handlers
function handleDragLeave(event: DragEvent) {
  // Check if the mouse left the container (not just moved between child elements)
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
  const x = event.clientX;
  const y = event.clientY;

  // Check if the mouse is outside the container's bounds
  if (x <= rect.left || x >= rect.right || y <= rect.top || y >= rect.bottom) {
    isDragging.value = false;
  }
}

// Handle drop function
function handleDrop(event: DragEvent) {
  event.preventDefault();
  isDragging.value = false;

  if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
    const file = event.dataTransfer.files[0];
    void handleFile(file);
  }
}

// Handle file function
async function handleFile(file: File) {
  try {
    isUploadingImage.value = true;
    console.log('Starting file validation for:', file.name);
    
    // Use the new validateFileUpload function
    const validation = await validateFileUpload(file, {
      maxFileSize: MAX_FILE_SIZE,
      maxDimension: MAX_DIMENSION
    });
    
    if (!validation.isValid) {
      console.error('File validation failed:', validation.error);
      $q.notify({
        type: 'negative',
        message: validation.error || 'Invalid file',
        position: 'top',
        timeout: 3000
      });
      return;
    }
    
    console.log('File validation passed');
    
    // If validation passed, we already have the sanitized data
    if (validation.sanitizedData) {
      previewUrl.value = validation.sanitizedData;
      emit('update:modelValue', validation.sanitizedData);
      
      $q.notify({
        type: 'positive',
        message: 'Image uploaded successfully',
        position: 'top',
        timeout: 2000
      });
    }
  } catch (error) {
    console.error('Error in handleFile:', error);
    previewUrl.value = props.defaultImageUrl;
    $q.notify({
      type: 'negative',
      message: 'An unexpected error occurred. Please try again.',
      position: 'top',
      timeout: 3000
    });
  } finally {
    isUploadingImage.value = false;
  }
}

// Remove image function
function removeImage(event?: Event) {
  if (event) {
    event.stopPropagation();
  }
  clearImageInput();
}

// Clear image input function
function clearImageInput() {
  if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(previewUrl.value);
  }
  previewUrl.value = props.defaultImageUrl;
  emit('update:modelValue', props.defaultImageUrl);
  
  if (fileInput.value) {
    fileInput.value.value = '';
  }
  isUploadingImage.value = false;
}

// Handle file select function
async function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files[0]) {
    const file = input.files[0];
    console.log('Selected file:', {
      name: file.name,
      type: file.type,
      size: file.size,
      lastModified: file.lastModified
    });

    try {
      await handleFile(file);
    } catch (error) {
      console.error('Error handling file:', error);
      $q.notify({
        type: 'negative',
        message: 'Error processing the image. Please try again.',
        position: 'top',
        timeout: 3000
      });
    }
  }
}

// Function to trigger file input
function triggerFileInput() {
  fileInput.value?.click();
}
</script>

<template>
  <div>
    <input
      ref="fileInput"
      type="file"
      accept="image/jpeg,image/png,image/gif"
      @change="handleFileSelect"
      class="hidden"
    />
    <div 
      class="upload-container q-pa-md" 
      :class="{ dragging: isDragging }"
      @click="triggerFileInput"
      @dragover.prevent="isDragging = true"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
    >
      <div v-if="isUploadingImage" class="text-center">
        <q-spinner color="primary" size="3em" />
        <div class="q-mt-sm">Uploading image...</div>
      </div>
      <div v-else-if="previewUrl" class="text-center">
        <img :src="previewUrl" class="preview-image" />
        <div class="q-mt-sm">
          <q-btn flat color="negative" size="sm" icon="delete" @click="removeImage">
            Remove Image
          </q-btn>
        </div>
      </div>
      <div v-else class="text-center">
        <q-icon name="upload_file" size="3em" color="grey-7" />
        <div class="text-subtitle1 q-mt-sm">Drop image here or click to upload</div>
        <div class="text-caption text-grey q-mt-xs">
          Supports: JPG, PNG, GIF (max 5MB)
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="sass" scoped>
.upload-container
  border: 2px dashed #ccc
  border-radius: 8px
  cursor: pointer
  transition: all 0.3s ease
  min-height: 200px
  display: flex
  align-items: center
  justify-content: center

  &:hover
    border-color: #00b4ff
    background: rgba(0, 180, 255, 0.05)

  &.dragging
    border-color: #00b4ff
    background: rgba(0, 180, 255, 0.1)

.preview-image
  width: 100%
  max-height: 180px
  object-fit: contain
  border-radius: 4px

.hidden
  display: none
</style> 