<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import type { AccessoryRow, NewAccessoryInput } from 'src/types/accessories';
import { getDefaultImage } from 'src/config/defaultImages';
import { validateAndSanitizeBase64Image } from 'src/utils/imageValidation';
import { operationNotifications } from 'src/utils/notifications';

const props = defineProps<{
  modelValue: boolean;
  makes: string[];
  colors: string[];
  accessory: AccessoryRow;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
  (e: 'submit', id: number, value: NewAccessoryInput): void;
}>();

const isDragging = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);
const imageUrlValid = ref(true);
const isUploadingImage = ref(false);
const defaultImageUrl = getDefaultImage('accessory');
const previewUrl = ref<string>(defaultImageUrl);

const newAccessory = ref<NewAccessoryInput>({
  name: '',
  make: '',
  quantity: 0,
  price: 0,
  unit_color: '',
  status: 'Out of Stock',
  image: defaultImageUrl,
});

// Watch for dialog visibility changes and accessory changes
watch(() => [props.modelValue, props.accessory], ([newModelValue, newAccessory]) => {
  if (newModelValue && newAccessory) {
    resetForm();
  }
}, { immediate: true });

const capitalizedName = computed({
  get: () => newAccessory.value.name,
  set: (value: string) => {
    if (value) {
      newAccessory.value.name = value.charAt(0).toUpperCase() + value.slice(1);
    } else {
      newAccessory.value.name = value;
    }
  }
});

// Add watch for quantity changes
watch(() => newAccessory.value.quantity, (newQuantity) => {
  if (newQuantity === 0) {
    newAccessory.value.status = 'Out of Stock';
  } else if (newQuantity <= 2) {
    newAccessory.value.status = 'Low Stock';
  } else if (newQuantity <= 5) {
    newAccessory.value.status = 'In Stock';
  } else {
    newAccessory.value.status = 'Available';
  }
});

function resetForm() {
  // Initialize the form with accessory data
  newAccessory.value = {
    name: props.accessory.name,
    make: props.accessory.make,
    quantity: props.accessory.quantity,
    price: props.accessory.price,
    unit_color: props.accessory.unit_color,
    status: props.accessory.status,
    image: props.accessory.image
  };

  // Handle the image preview for base64 images
  if (props.accessory.image.startsWith('data:image/')) {
    const validationResult = validateAndSanitizeBase64Image(props.accessory.image);
    if (validationResult.isValid) {
      previewUrl.value = validationResult.sanitizedData!;
      newAccessory.value.image = validationResult.sanitizedData!;
      imageUrlValid.value = true;
    } else {
      previewUrl.value = defaultImageUrl;
      newAccessory.value.image = defaultImageUrl;
      imageUrlValid.value = true;
      operationNotifications.validation.warning('Invalid image data, using default image');
    }
  } else {
    // For any other case, use default image
    previewUrl.value = defaultImageUrl;
    newAccessory.value.image = defaultImageUrl;
    imageUrlValid.value = true;
  }
}

function closeDialog() {
  emit('update:modelValue', false);
  clearImageInput();
}

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

function handleDrop(event: DragEvent) {
  event.preventDefault();
  isDragging.value = false;

  if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
    const file = event.dataTransfer.files[0];
    void handleFile(file);
  }
}

// Function to handle file selection
async function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files[0]) {
    const file = input.files[0];
    await handleFile(file);
  }
}

// Function to handle file processing
async function handleFile(file: File) {
  try {
    isUploadingImage.value = true;

    // Validate file type
    if (!file.type.startsWith('image/')) {
      throw new Error('Please upload an image file');
    }

    // Create a blob URL for preview
    const blobUrl = URL.createObjectURL(file);
    previewUrl.value = blobUrl;

    // Convert file to base64
    const base64Data = await new Promise<string>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(reader.result as string);
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });

    // Validate and sanitize the base64 image
    const validationResult = validateAndSanitizeBase64Image(base64Data);
    if (!validationResult.isValid) {
      throw new Error('Invalid image data');
    }

    newAccessory.value.image = validationResult.sanitizedData!;
    imageUrlValid.value = true;
  } catch (error) {
    console.error('Error handling file:', error);
    operationNotifications.validation.error(error instanceof Error ? error.message : 'Error processing image');
    clearImageInput();
  } finally {
    isUploadingImage.value = false;
  }
}

function clearImageInput() {
  if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(previewUrl.value);
  }
  previewUrl.value = defaultImageUrl;
  newAccessory.value.image = defaultImageUrl;
  imageUrlValid.value = true;
  if (fileInput.value) {
    fileInput.value.value = '';
  }
  isUploadingImage.value = false;
}

function removeImage() {
  clearImageInput();
}

function updateAccessory() {
  try {
    // Validate required fields
    if (!newAccessory.value.make || !newAccessory.value.unit_color) {
      operationNotifications.validation.error('Please fill in all required fields');
      return;
    }

    // Validate image
    if (!imageUrlValid.value) {
      operationNotifications.validation.error('Please provide a valid image');
      return;
    }

    // If no image is uploaded, use default
    if (!newAccessory.value.image || newAccessory.value.image === '') {
      newAccessory.value.image = defaultImageUrl;
    }

    // Emit the submit event with the accessory id and updated data
    emit('submit', props.accessory.id, newAccessory.value);
    closeDialog();
  } catch (error) {
    console.error('Error processing accessory update:', error);
    operationNotifications.update.error('accessory');
  }
}
</script>

<template>
  <q-dialog
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
    @hide="clearImageInput"
  >
    <q-card style="min-width: 400px; max-width: 95vw">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">Edit Accessory</div>
        <q-space />
        <q-btn icon="close" flat round dense @click="closeDialog" />
      </q-card-section>

      <q-card-section>
        <q-form @submit.prevent="updateAccessory" class="q-gutter-sm">
          <q-input
            v-model="capitalizedName"
            label="Accessory Name"
            dense
            outlined
            required
            :rules="[val => !!val || 'Name is required']"
          >
            <template v-slot:prepend>
              <q-icon name="build" />
            </template>
          </q-input>

          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <q-select
                v-model="newAccessory.make"
                :options="makes"
                label="Make"
                dense
                outlined
                required
                emit-value
                map-options
                placeholder="Select a make"
                :rules="[val => !!val || 'Make is required']"
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

            <div class="col-12 col-sm-6">
              <q-select
                v-model="newAccessory.unit_color"
                :options="colors"
                label="Color"
                dense
                outlined
                required
                emit-value
                map-options
                placeholder="Select a color"
                :rules="[val => !!val || 'Color is required']"
              >
                <template v-slot:prepend>
                  <q-icon name="palette" />
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
                v-model.number="newAccessory.quantity"
                type="number"
                min="0"
                label="Quantity"
                dense
                outlined
                required
                :rules="[val => val >= 0 || 'Quantity must be positive']"
              >
                <template v-slot:prepend>
                  <q-icon name="numbers" />
                </template>
              </q-input>
            </div>

            <div class="col-12 col-sm-6">
              <q-input
                v-model.number="newAccessory.price"
                type="number"
                label="Price"
                dense
                outlined
                required
                :rules="[val => val > 0 || 'Price must be greater than 0']"
              >
                <template v-slot:prepend>
                  <q-icon name="attach_money" />
                </template>
              </q-input>
            </div>
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-12">
              <q-input
                v-model="newAccessory.status"
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

          <div class="row q-col-gutter-sm">
            <div class="col-12">
              <div
                class="upload-container"
                :class="{ dragging: isDragging }"
                @dragenter.prevent="isDragging = true"
                @dragover.prevent="isDragging = true"
                @dragleave.prevent="handleDragLeave"
                @drop.prevent="handleDrop"
                @click="fileInput?.click()"
              >
                <input
                  ref="fileInput"
                  type="file"
                  accept="image/*"
                  class="hidden"
                  @change="handleFileSelect"
                />
                <div v-if="isUploadingImage" class="text-center">
                  <q-spinner-dots color="primary" size="40px" />
                  <div class="text-subtitle2 q-mt-sm">Processing image...</div>
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
                  </div>
                </div>
                <div v-else class="row items-center justify-center">
                  <div class="col-8 text-center">
                    <img
                      :src="previewUrl"
                      class="preview-image"
                      :alt="newAccessory.name || 'Preview image'"
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
        <q-btn flat label="Cancel" @click="closeDialog" />
        <q-btn
          unelevated
          color="primary"
          label="Update Accessory"
          @click="updateAccessory"
          :disable="!newAccessory.name || !newAccessory.make || !newAccessory.unit_color || newAccessory.quantity < 0 || newAccessory.price <= 0"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style lang="sass" scoped>
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
