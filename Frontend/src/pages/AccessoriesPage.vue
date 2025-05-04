<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import ProductCardModal from 'src/components/Global/ProductModal.vue'
import { useAccessoriesStore } from 'src/stores/accessories';
import type { AccessoryRow, NewAccessoryInput } from 'src/types/accessories';
import { getFirstFallbackImage } from 'src/config/defaultImages';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';

const store = useAccessoriesStore();
const showFilterDialog = ref(false);
const showAddDialog = ref(false);
const showEditDialog = ref(false);
const showDeleteDialog = ref(false);
const accessoryToDelete = ref<AccessoryRow | null>(null);
const showProductCardModal = ref(false);
const isDragging = ref(false);

// Track drag leave timeout for debouncing
let dragLeaveTimeout: number | null = null;
// Add a small buffer zone (in pixels) around the drop target
const DRAG_BOUNDARY_BUFFER = 50;

// Add at the top with other imports and constants
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif'] as const;
type AllowedMimeType = typeof ALLOWED_TYPES[number];

/**
 * Global drag event handlers are necessary to handle edge cases in drag-and-drop operations:
 * 1. When the user drags outside the browser window and releases the mouse
 * 2. When the drag operation ends outside the drop zone
 * 3. When the drag operation is cancelled with ESC key
 * 
 * Without these handlers, the isDragging state might remain true in these cases,
 * leading to the upload zone staying in a dragging visual state incorrectly.
 * Local dragend events within the drop zone alone are not sufficient to catch
 * all possible ways a drag operation can end.
 */
onMounted(() => {
  const handleGlobalDragEnd = () => {
    isDragging.value = false;
    if (dragLeaveTimeout) {
      clearTimeout(dragLeaveTimeout);
      dragLeaveTimeout = null;
    }
  };

  document.addEventListener('dragend', handleGlobalDragEnd);

  // Clean up on unmount to prevent memory leaks
  onUnmounted(() => {
    document.removeEventListener('dragend', handleGlobalDragEnd);
    if (dragLeaveTimeout) {
      clearTimeout(dragLeaveTimeout);
      dragLeaveTimeout = null;
    }
  });
});

/**
 * Handles the drag leave event with debouncing and a buffer zone.
 * This prevents flickering of the drop zone styling when:
 * 1. The drag briefly leaves and re-enters the zone
 * 2. The drag moves near the edges of the zone
 * 3. The user moves the mouse quickly within the zone
 * 
 * @param event - The drag event object
 */
function handleDragLeave(event: DragEvent) {
  // Clear any existing timeout to prevent premature state changes
  if (dragLeaveTimeout) {
    clearTimeout(dragLeaveTimeout);
    dragLeaveTimeout = null;
  }

  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
  const x = event.clientX;
  const y = event.clientY;

  // Check if the mouse is significantly outside the container's bounds
  // Add a buffer zone to make the interaction more forgiving
  const isOutsideBounds = 
    x <= rect.left - DRAG_BOUNDARY_BUFFER || 
    x >= rect.right + DRAG_BOUNDARY_BUFFER || 
    y <= rect.top - DRAG_BOUNDARY_BUFFER || 
    y >= rect.bottom + DRAG_BOUNDARY_BUFFER;

  if (isOutsideBounds) {
    // Debounce the state change to prevent flickering
    dragLeaveTimeout = window.setTimeout(() => {
      isDragging.value = false;
      dragLeaveTimeout = null;
    }, 100) as unknown as number; // 100ms debounce
  }
}

function handleDrop(event: DragEvent) {
  event.preventDefault();
  // Clear any pending drag leave timeout
  if (dragLeaveTimeout) {
    clearTimeout(dragLeaveTimeout);
    dragLeaveTimeout = null;
  }
  isDragging.value = false;

  if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
    const file = event.dataTransfer.files[0];
    void handleFile(file);
  }
}

const selected = ref<AccessoryRow>({
  name: '',
  id: 0,
  make: 'Generic',
  quantity: 0,
  price: 0,
  unit_color: 'Black',
  status: 'Out of Stock',
  image: '',
})

const newAccessory = ref<NewAccessoryInput>({
  name: '',
  make: '',
  quantity: 0,
  price: 0,
  unit_color: '',
  status: 'Out of Stock',
  image: 'https://loremflickr.com/600/400/accessory',
})

// Image validation
const imageUrlValid = ref(true);
const defaultImageUrl = getFirstFallbackImage('accessory');

// Available options from store
const { makes, colors, statuses } = store;

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

// Update status based on quantity
function updateStatus(quantity: number) {
  if (quantity === 0) {
    return 'Out of Stock';
  } else if (quantity <= 2) {
    return 'Low Stock';
  } else if (quantity <= 5) {
    return 'In Stock';
  } else {
    return 'Available';
  }
}

// Watch for quantity changes and update status
watch(() => newAccessory.value.quantity, (newQuantity) => {
  newAccessory.value.status = updateStatus(newQuantity);
});

const columns: QTableColumn[] = [
  { name: 'id', align: 'center', label: 'ID', field: 'id', sortable: true },
  {
    name: 'unitName',
    required: true,
    label: 'Unit Name',
    align: 'left',
    field: 'name',
    sortable: true
  },
  { name: 'make', label: 'Make', field: 'make' },
  { name: 'quantity', label: 'Quantity', field: 'quantity', sortable: true },
  { name: 'price', label: 'Price', field: 'price', sortable: true, format: (val: number) =>
        `₱ ${val.toLocaleString('en-PH', {
          minimumFractionDigits: 2,
          maximumFractionDigits: 2
        })}`
  },
  { name: 'status', label: 'Status', field: 'status' },
  { name: 'color', label: 'Color', field: 'unit_color' },
  {
    name: 'actions',
    label: 'Actions',
    field: 'actions',
    align: 'center',
    sortable: false
  }
];

/**
 * Validates and processes a base64 image, falling back to default if invalid.
 * This function handles the validation and sanitization of base64 encoded images,
 * updating the image state appropriately based on validation results.
 * 
 * @param base64Image - The base64 image string to validate
 * @param defaultImage - The fallback image URL to use if validation fails
 * @returns An object containing the validated image URL and validation status
 */
function validateAndProcessBase64Image(base64Image: string, defaultImage: string): { 
  imageUrl: string;
  isValid: boolean;
} {
  const validationResult = validateAndSanitizeBase64Image(base64Image);
  
  if (validationResult.isValid && validationResult.sanitizedData) {
    return {
      imageUrl: validationResult.sanitizedData,
      isValid: true
    };
  }

  // If validation fails, notify user and use default image
  operationNotifications.validation.warning('Invalid image data, using default image');
  return {
    imageUrl: defaultImage,
    isValid: true // We consider it valid since we're using a fallback
  };
}

/**
 * Handles image validation and state updates for both new and existing images.
 * This function centralizes the image processing logic to ensure consistent
 * behavior across different parts of the application.
 * 
 * @param imageUrl - The image URL or base64 string to validate
 * @param defaultUrl - The fallback URL to use if validation fails
 * @returns An object containing the processed image URL and validation state
 */
function handleImageValidation(imageUrl: string | undefined | null, defaultUrl: string): {
  processedUrl: string;
  isValid: boolean;
  previewUrl: string;
} {
  // Handle undefined or empty image URLs
  if (!imageUrl) {
    return {
      processedUrl: defaultUrl,
      isValid: true,
      previewUrl: defaultUrl
    };
  }

  // Handle base64 images
  if (imageUrl.startsWith('data:image/')) {
    const { imageUrl: validatedUrl, isValid } = validateAndProcessBase64Image(imageUrl, defaultUrl);
    return {
      processedUrl: validatedUrl,
      isValid,
      previewUrl: validatedUrl
    };
  }

  // For regular URLs, use as is with default fallback preview
  return {
    processedUrl: imageUrl,
    isValid: true,
    previewUrl: defaultUrl
  };
}

const onRowClick: QTableProps['onRowClick'] = (evt, row) => {
  // Check if the click originated from the action button or its menu
  const target = evt.target as HTMLElement;
  if (target.closest('.action-button') || target.closest('.action-menu')) {
    return;
  }
  
  // Update selected with a proper copy of the row data
  selected.value = { ...row as AccessoryRow };
  
  // Validate and set the image using the centralized handler
  const { processedUrl, isValid } = handleImageValidation(selected.value.image, defaultImageUrl);
  selected.value.image = processedUrl;
  imageUrlValid.value = isValid;
  
  showProductCardModal.value = true;
}

function addToCart() {
  if (selected.value) {
    console.log('added to cart for', selected.value.name);
  }
  showProductCardModal.value = false;
}

function openAddDialog() {
  newAccessory.value = {
    name: '',
    make: '',
    quantity: 0,
    price: 0,
    unit_color: '',
    status: 'Out of Stock',
    image: defaultImageUrl
  };
  
  // Use the centralized handler for consistent image state initialization
  const { previewUrl: newPreviewUrl, isValid } = handleImageValidation(defaultImageUrl, defaultImageUrl);
  previewUrl.value = newPreviewUrl;
  imageUrlValid.value = isValid;
  
  if (fileInput.value) {
    fileInput.value.value = ''; // Clear the file input
  }
  showAddDialog.value = true;
}

async function addNewAccessory() {
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

    // Execute the store action and await its completion
    const result = await store.addAccessory(newAccessory.value);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showAddDialog.value = false;
      clearImageInput(); // Clear the image input state
      operationNotifications.add.success(`accessory: ${newAccessory.value.name}`);
    }
  } catch (error) {
    console.error('Error adding accessory:', error);
    operationNotifications.add.error('accessory');
  }
}

// Image handling functions
const fileInput = ref<HTMLInputElement | null>(null);
const previewUrl = ref<string>(defaultImageUrl);
const isUploadingImage = ref(false);

/**
 * Clears all image-related state and resets the file input.
 * This is necessary to:
 * 1. Prevent the same file from being re-selected (browsers won't fire change event)
 * 2. Clear security-sensitive file references from memory
 * 3. Allow the same file to be uploaded again if needed
 * 4. Reset UI state for a clean user experience
 */
function clearImageInput() {
  // Clean up any existing blob URLs to prevent memory leaks
  if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(previewUrl.value);
  }
  
  // Reset all image-related state to defaults
  previewUrl.value = defaultImageUrl;
  newAccessory.value.image = defaultImageUrl;
  imageUrlValid.value = true;
  
  // Clear the file input value to ensure change events fire even if 
  // the same file is selected again. This is crucial for proper UX
  // when users want to re-upload the same file after an error.
  if (fileInput.value) {
    fileInput.value.value = '';
  }
  
  isUploadingImage.value = false;
}

// Function to handle edit accessory
function editAccessory(accessory: AccessoryRow) {
  selected.value = { ...accessory };
  newAccessory.value = {
    name: accessory.name,
    make: accessory.make,
    quantity: accessory.quantity,
    price: accessory.price,
    unit_color: accessory.unit_color,
    status: accessory.status,
    image: accessory.image
  };

  // Handle image validation and state updates using the centralized handler
  const { processedUrl, isValid, previewUrl: newPreviewUrl } = handleImageValidation(
    accessory.image,
    defaultImageUrl
  );
  
  newAccessory.value.image = processedUrl;
  imageUrlValid.value = isValid;
  previewUrl.value = newPreviewUrl;
  
  showEditDialog.value = true;
}

// Function to handle update accessory
async function updateAccessory() {
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

    if (!selected.value) {
      throw new Error('No accessory selected for update');
    }

    // Execute the store action and await its completion
    const result = await store.updateAccessory(selected.value.id, newAccessory.value);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showEditDialog.value = false;
      clearImageInput(); // Clear the image input state
      operationNotifications.update.success(`accessory: ${newAccessory.value.name}`);
    }
  } catch (error) {
    console.error('Error updating accessory:', error);
    operationNotifications.update.error('accessory');
  }
}

// Function to handle delete accessory
function deleteAccessory(accessory: AccessoryRow) {
  accessoryToDelete.value = accessory;
  showDeleteDialog.value = true;
}

// Function to confirm and execute delete
async function confirmDelete() {
  try {
    if (!accessoryToDelete.value) return;

    await store.deleteAccessory(accessoryToDelete.value.id);
    showDeleteDialog.value = false;
    accessoryToDelete.value = null;
    operationNotifications.delete.success('accessory');
  } catch (error) {
    console.error('Error deleting accessory:', error);
    operationNotifications.delete.error('accessory');
  }
}

/**
 * Processes and validates an uploaded image file.
 * Handles file type checking, conversion to base64, and validation.
 * 
 * @param file - The uploaded image file to process
 * @returns Promise resolving to the processed image data
 * @throws Error if file validation or processing fails
 */
async function processImageUpload(file: File): Promise<string> {
  // Validate file type
  if (!file.type.startsWith('image/')) {
    throw new Error('Please upload an image file');
  }

  // Convert file to base64
  const base64Data = await new Promise<string>((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });

  // Validate and sanitize the base64 image
  const validationResult = validateAndSanitizeBase64Image(base64Data);
  if (!validationResult.isValid || !validationResult.sanitizedData) {
    throw new Error('Invalid image data');
  }

  return validationResult.sanitizedData;
}

/**
 * Creates a preview URL for an image file.
 * Handles cleanup of previous preview URLs to prevent memory leaks.
 * 
 * @param file - The image file to preview
 * @returns The blob URL for the preview
 */
function createImagePreview(file: File): string {
  // Clean up any existing preview URL
  if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(previewUrl.value);
  }
  return URL.createObjectURL(file);
}

/**
 * Handles the complete image file upload process including validation,
 * preview generation, and state updates.
 * 
 * @param file - The uploaded image file
 */
async function handleFile(file: File) {
  try {
    isUploadingImage.value = true;

    // Generate preview immediately for better UX
    previewUrl.value = createImagePreview(file);

    // Process and validate the image
    const processedImage = await processImageUpload(file);
    
    // Update component state with the processed image
    newAccessory.value.image = processedImage;
    imageUrlValid.value = true;

    // Show success notification using the appropriate notification type
    operationNotifications.validation.warning('Image uploaded successfully');
  } catch (error) {
    console.error('Error handling file:', error);
    operationNotifications.validation.error(error instanceof Error ? error.message : 'Error processing image');
    clearImageInput();
  } finally {
    isUploadingImage.value = false;
  }
}

/**
 * Handles file selection from input element.
 * Validates file size and type before processing.
 * 
 * @param event - The file input change event
 */
async function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  
  if (!file) {
    operationNotifications.validation.warning('No file selected');
    return;
  }

  console.log('Selected file:', {
    name: file.name,
    type: file.type,
    size: file.size
  });

  // Validate file size
  if (file.size > MAX_FILE_SIZE) {
    const sizeMB = (file.size / (1024 * 1024)).toFixed(2);
    operationNotifications.validation.error(
      `File size (${sizeMB}MB) exceeds the 5MB limit`
    );
    return;
  }

  // Validate file type
  if (!ALLOWED_TYPES.includes(file.type as AllowedMimeType)) {
    operationNotifications.validation.error(
      `Invalid file type: ${file.type}. Allowed types are: JPEG, PNG, and GIF`
    );
    return;
  }

  await handleFile(file);
}

function removeImage() {
  clearImageInput();
}

function applyFilters() {
  showFilterDialog.value = false;
}

// Initialize data when component is mounted
onMounted(() => {
  void store.initializeAccessories();
});
</script>

<template>
  <div class="q-pa-md">
    <div class="q-pa-sm full-width">
      <div class="flex row q-my-sm">
        <div class="flex full-width col">
          <div class="flex col q-mr-sm">
            <q-input
              v-model="store.rawAccessorySearch"
              outlined
              dense
              placeholder="Search"
              class="full-width"
            >
              <template v-slot:prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>
          <div class="flex col">
            <q-btn
              outline
              icon="filter_list"
              label="Filters"
              @click="showFilterDialog = true"
            />
          </div>
        </div>

        <div class="flex row q-gutter-x-sm">
          <q-btn class="text-white bg-primary" unelevated @click="openAddDialog">
            <q-icon name="add" color="white" />
            Add
          </q-btn>
          <div class="flex row">
            <q-btn dense flat class="bg-primary text-white q-pa-sm">
              <q-icon name="download" color="white" />
              Download CSV
            </q-btn>
          </div>
        </div>
      </div>

      <!--ACCESSORIES TABLE-->
      <q-table
        class="my-sticky-column-table"
        flat
        bordered
        title="Accessories"
        :rows="store.filteredAccessoryRows"
        :columns="columns"
        row-key="id"
        :filter="store.accessorySearch"
        @row-click="onRowClick"
        :pagination="{ rowsPerPage: 5 }"
        :loading="store.isLoading"
      >
        <template v-slot:loading>
          <q-inner-loading showing color="primary">
            <q-spinner-gears size="50px" color="primary" />
          </q-inner-loading>
        </template>
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" auto-width :key="props.row.id">
            <q-btn
              flat
              round
              dense
              color="grey"
              icon="more_vert"
              class="action-button"
              :aria-label="'Actions for ' + props.row.name"
            >
              <q-menu class="action-menu" :aria-label="'Available actions for ' + props.row.name">
                <q-list style="min-width: 100px">
                  <q-item
                    clickable
                    v-close-popup
                    @click.stop="editAccessory(props.row)"
                    role="button"
                    :aria-label="'Edit ' + props.row.name"
                  >
                    <q-item-section>
                      <q-item-label>
                        <q-icon name="edit" size="xs" class="q-mr-sm" aria-hidden="true" />
                        Edit
                      </q-item-label>
                    </q-item-section>
                  </q-item>
                  <q-item
                    clickable
                    v-close-popup
                    @click.stop="deleteAccessory(props.row)"
                    role="button"
                    :aria-label="'Delete ' + props.row.name"
                    class="text-negative"
                  >
                    <q-item-section>
                      <q-item-label class="text-negative">
                        <q-icon name="delete" size="xs" class="q-mr-sm" aria-hidden="true" />
                        Delete
                      </q-item-label>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </q-td>
        </template>
      </q-table>

      <!-- Existing Accessory Modal -->
      <ProductCardModal
        v-model="showProductCardModal"
        :image="selected?.image || ''"
        :title="selected?.name || ''"
        :unit_color="selected?.unit_color || ''"
        :price="selected?.price || 0"
        :quantity="selected?.quantity || 0"
        :details="`${selected?.make}`"
        :status="selected?.status || ''"
        @add="addToCart"
      />

      <!-- Add Accessory Dialog -->
      <q-dialog
        v-model="showAddDialog"
        persistent
        @hide="clearImageInput"
      >
        <q-card style="min-width: 400px; max-width: 95vw">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6">New Accessory</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>

          <q-card-section>
            <q-form @submit.prevent="addNewAccessory" class="q-gutter-sm">
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
                    <div v-else class="row items-center justify-center">
                      <div class="col-8 text-center">
                        <img
                          :src="previewUrl"
                          class="preview-image"
                          :alt="newAccessory.name || 'Preview image'"
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
            <q-btn flat label="Cancel" v-close-popup />
            <q-btn
              unelevated
              color="primary"
              label="Add Accessory"
              @click="addNewAccessory"
              :disable="!newAccessory.name || !newAccessory.make || !newAccessory.unit_color || newAccessory.quantity < 0 || newAccessory.price <= 0"
            />
          </q-card-actions>
        </q-card>
      </q-dialog>

      <!-- Filter Dialog -->
      <q-dialog v-model="showFilterDialog">
        <q-card style="min-width: 350px">
          <q-card-section class="row items-center">
            <div class="text-h6">Filter Accessories</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>

          <q-card-section class="q-pt-none">
            <q-select
              v-model="store.filterMake"
              :options="makes"
              label="Make"
              clearable
              outlined
              class="q-mb-md"
            />

            <q-select
              v-model="store.filterColor"
              :options="colors"
              label="Color"
              clearable
              outlined
              class="q-mb-md"
            />

            <q-select
              v-model="store.filterStatus"
              :options="statuses"
              label="Status"
              clearable
              outlined
              class="q-mb-md"
            />
          </q-card-section>

          <q-card-actions align="right" class="text-primary">
            <q-btn flat label="Reset" color="negative" @click="store.resetFilters" />
            <q-btn flat label="Apply Filters" @click="applyFilters" />
          </q-card-actions>
        </q-card>
      </q-dialog>

      <!-- Edit Accessory Dialog -->
      <q-dialog
        v-model="showEditDialog"
        persistent
        @hide="clearImageInput"
      >
        <q-card style="min-width: 400px; max-width: 95vw">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6">Edit Accessory</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
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
                    <div v-else class="row items-center justify-center">
                      <div class="col-8 text-center">
                        <img
                          :src="previewUrl"
                          class="preview-image"
                          :alt="newAccessory.name || 'Preview image'"
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
            <q-btn flat label="Cancel" v-close-popup />
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

      <!-- Delete Confirmation Dialog -->
      <q-dialog v-model="showDeleteDialog" persistent>
        <q-card>
          <q-card-section class="row items-center">
            <q-avatar icon="warning" color="negative" text-color="white" />
            <span class="q-ml-sm text-h6">Delete Accessory</span>
          </q-card-section>

          <q-card-section>
            Are you sure you want to delete {{ accessoryToDelete?.name }}? This action cannot be undone.
          </q-card-section>

          <q-card-actions align="right">
            <q-btn flat label="Cancel" v-close-popup />
            <q-btn flat label="Delete" color="negative" @click="confirmDelete" />
          </q-card-actions>
        </q-card>
      </q-dialog>
    </div>
  </div>
</template>

<style lang="sass">
.my-sticky-column-table
  max-width: 100%

  thead tr:first-child th:nth-child(2)
    background-color: var(--sticky-column-bg)

  td:nth-child(2)
    background-color: var(--sticky-column-bg)

  th:nth-child(2),
  td:nth-child(2)
    position: sticky
    left: 0
    z-index: 1
    color: white

    .body--dark &
      color: black

.accessory-page-z-top
  z-index: 1000

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

.action-button
  position: relative
  z-index: 1

.action-menu
  z-index: 1001 !important
</style> 