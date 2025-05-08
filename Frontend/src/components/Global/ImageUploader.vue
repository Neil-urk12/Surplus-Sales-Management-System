<script setup lang="ts">
import { ref, watch, defineProps, defineEmits } from 'vue';
import { useQuasar } from 'quasar';
import { validateAndSanitizeBase64Image } from 'src/utils/imageValidation';

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
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif'] as const;
const MAX_DIMENSION = 4096; // Maximum image dimension in pixels

type AllowedMimeType = typeof ALLOWED_TYPES[number];

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
    const validation = await validateFile(file);
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

    // Create a temporary URL for preview
    console.log('Creating preview URL');
    const tempPreviewUrl = URL.createObjectURL(file);
    previewUrl.value = tempPreviewUrl;
    console.log('Preview URL set:', previewUrl.value);

    console.log('Starting FileReader');
    const reader = new FileReader();

    reader.onload = (e) => {
      console.log('FileReader loaded');
      if (e.target?.result) {
        const base64String = e.target.result as string;
        console.log('Processing base64 data');
        const base64ValidationResult = validateAndSanitizeBase64Image(base64String);

        if (!base64ValidationResult.isValid) {
          console.error('Base64 validation failed:', base64ValidationResult.error);
          $q.notify({
            type: 'negative',
            message: base64ValidationResult.error || 'Invalid image data',
            position: 'top',
            timeout: 3000
          });
          previewUrl.value = props.defaultImageUrl;
          return;
        }

        console.log('Base64 validation passed, updating image');
        emit('update:modelValue', base64ValidationResult.sanitizedData);

        $q.notify({
          type: 'positive',
          message: 'Image uploaded successfully',
          position: 'top',
          timeout: 2000
        });
      }
    };

    reader.onerror = (error) => {
      console.error('FileReader error:', error);
      previewUrl.value = props.defaultImageUrl;
      $q.notify({
        type: 'negative',
        message: 'Error reading file. Please try again.',
        position: 'top',
        timeout: 3000
      });
    };

    console.log('Starting file read');
    reader.readAsDataURL(file);
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

// Validate file function
async function validateFile(file: File): Promise<{ isValid: boolean; error?: string }> {
  try {
    console.log('Starting file validation:', {
      name: file.name,
      type: file.type,
      size: file.size
    });

    // Step 1: Basic file validation
    if (!file) {
      console.error('Validation failed: No file provided');
      return { isValid: false, error: 'No file provided.' };
    }

    // Step 2: Size validation
    if (file.size > MAX_FILE_SIZE) {
      const sizeMB = (file.size / (1024 * 1024)).toFixed(2);
      console.error(`Validation failed: File size ${sizeMB}MB exceeds limit`);
      return {
        isValid: false,
        error: `File size (${sizeMB}MB) exceeds 5MB limit. Please choose a smaller file.`
      };
    }

    // Step 3: Enhanced MIME type validation with file signature check
    const validMimeTypes = {
      'image/jpeg': [0xFF, 0xD8, 0xFF],
      'image/png': [0x89, 0x50, 0x4E, 0x47],
      'image/gif': [0x47, 0x49, 0x46, 0x38]
    };

    if (!Object.keys(validMimeTypes).includes(file.type)) {
      console.error(`Validation failed: Invalid file type ${file.type}`);
      return {
        isValid: false,
        error: `Invalid file type: ${file.type}. Please upload a JPG, PNG, or GIF image.`
      };
    }

    // Step 4: File signature validation
    const arrayBuffer = await file.slice(0, 4).arrayBuffer();
    const bytes = new Uint8Array(arrayBuffer);
    const expectedSignature = validMimeTypes[file.type as keyof typeof validMimeTypes];

    const isValidSignature = expectedSignature.every((byte, i) => byte === bytes[i]);
    if (!isValidSignature) {
      console.error('Validation failed: File signature mismatch');
      return {
        isValid: false,
        error: 'Invalid image format. The file content does not match its extension.'
      };
    }

    // Step 5: Validate image dimensions
    try {
      const dimensionValidation = await validateImageDimensions(file);
      if (!dimensionValidation.isValid) {
        console.error('Validation failed:', dimensionValidation.error);
        return dimensionValidation;
      }
    } catch (error) {
      console.error('Error validating image dimensions:', error);
      return {
        isValid: false,
        error: 'Error validating image dimensions. Please try again.'
      };
    }

    console.log('File validation passed successfully');
    return { isValid: true };
  } catch (error) {
    console.error('Unexpected error during file validation:', error);
    return {
      isValid: false,
      error: 'An unexpected error occurred during validation. Please try again.'
    };
  }
}

// Validate image dimensions function
function validateImageDimensions(file: File): Promise<{ isValid: boolean; error?: string }> {
  return new Promise((resolve) => {
    const img = new Image();
    const objectUrl = URL.createObjectURL(file);

    const cleanup = () => {
      URL.revokeObjectURL(objectUrl);
    };

    // Set a timeout to avoid hanging
    const timeout = setTimeout(() => {
      cleanup();
      console.error('Dimension validation timed out');
      resolve({
        isValid: false,
        error: 'Image validation timed out. Please try again.'
      });
    }, 10000); // 10 second timeout

    img.onload = () => {
      clearTimeout(timeout);
      cleanup();
      console.log('Image dimensions:', {
        width: img.width,
        height: img.height,
        maxAllowed: MAX_DIMENSION
      });

      if (img.width > MAX_DIMENSION || img.height > MAX_DIMENSION) {
        resolve({
          isValid: false,
          error: `Image dimensions (${img.width}x${img.height}) cannot exceed ${MAX_DIMENSION}x${MAX_DIMENSION} pixels.`
        });
      } else if (img.width === 0 || img.height === 0) {
        resolve({
          isValid: false,
          error: 'Invalid image dimensions.'
        });
      } else {
        resolve({ isValid: true });
      }
    };

    img.onerror = () => {
      clearTimeout(timeout);
      cleanup();
      resolve({
        isValid: false,
        error: 'Error loading image. Please ensure it is a valid image file.'
      });
    };

    img.src = objectUrl;
  });
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

    // Check file type
    if (!ALLOWED_TYPES.includes(file.type as AllowedMimeType)) {
      $q.notify({
        type: 'negative',
        message: `Invalid file type: ${file.type}. Allowed types are: JPEG, PNG, and GIF`,
        position: 'top',
        timeout: 3000
      });
      return;
    }

    // Check file size
    if (file.size > MAX_FILE_SIZE) {
      $q.notify({
        type: 'negative',
        message: `File size (${(file.size / 1024 / 1024).toFixed(2)}MB) exceeds the 5MB limit`,
        position: 'top',
        timeout: 3000
      });
      return;
    }

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