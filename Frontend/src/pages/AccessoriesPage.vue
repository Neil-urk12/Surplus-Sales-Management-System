<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import ProductCardModal from 'src/components/Global/ProductModal.vue'
import { useAccessoriesStore } from 'src/stores/accessories';
import type { AccessoryRow, NewAccessoryInput } from 'src/types/accessories';
import { getDefaultImage } from 'src/config/defaultImages';
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

// Add global drag event handlers
onMounted(() => {
  const handleGlobalDragEnd = () => {
    isDragging.value = false;
  };

  document.addEventListener('dragend', handleGlobalDragEnd);

  // Clean up on unmount
  onUnmounted(() => {
    document.removeEventListener('dragend', handleGlobalDragEnd);
  });
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

function handleDrop(event: DragEvent) {
  event.preventDefault();
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
const defaultImageUrl = getDefaultImage('accessory');

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

const onRowClick: QTableProps['onRowClick'] = (evt, row) => {
  // Check if the click originated from the action button or its menu
  const target = evt.target as HTMLElement;
  if (target.closest('.action-button') || target.closest('.action-menu')) {
    return; // Do nothing if clicked on action button or its menu
  }
  
  // Update selected with a proper copy of the row data
  selected.value = { ...row as AccessoryRow };
  
  // Validate and set the image
  if (selected.value.image) {
    if (selected.value.image.startsWith('data:image/')) {
      // For base64 images, validate but preserve the original if it's a valid format
      const validationResult = validateAndSanitizeBase64Image(selected.value.image);
      if (validationResult.isValid && validationResult.sanitizedData) {
        selected.value.image = validationResult.sanitizedData;
      }
      // Even if validation fails, we'll let the ProductCardModal handle the fallback
    }
  } else {
    // If no image, use default
    selected.value.image = defaultImageUrl;
  }
  
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
  previewUrl.value = defaultImageUrl;
  imageUrlValid.value = true;
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

  // Handle the image preview for base64 images
  if (accessory.image.startsWith('data:image/')) {
    const validationResult = validateAndSanitizeBase64Image(accessory.image);
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