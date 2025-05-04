<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import ProductCardModal from 'src/components/Global/ProductModal.vue'
import { useQuasar } from 'quasar';
import { useMaterialsStore } from 'src/stores/materials';
import type { MaterialRow, NewMaterialInput } from 'src/stores/materials';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';
import AddMaterialDialog from 'src/components/AddMaterialDialog.vue';
import EditMaterialDialog from 'src/components/EditMaterialDialog.vue';
import FilterMaterialDialog from 'src/components/FilterMaterialDialog.vue';
import ConfirmDialog from 'src/components/ConfirmDialog.vue';
import { getFirstFallbackImage } from 'src/config/defaultImages';

const $q = useQuasar();
const store = useMaterialsStore();
const showFilterDialog = ref(false);
const selectedMaterial = ref<MaterialRow>({
  name: '',
  id: 0,
  category: 'Building',
  supplier: 'Steel Co.',
  quantity: 0,
  status: 'Out of Stock',
  image: ''
})

const newMaterial = ref<NewMaterialInput>({
  name: '',
  category: '',
  supplier: '',
  quantity: 0,
  status: 'Out of Stock',
  image: 'https://loremflickr.com/600/400/material'
})

// Image validation
const imageUrlValid = ref(true);
const validatingImage = ref(false);
const defaultImageUrl = getFirstFallbackImage('material');

// Available options from store
const { categories, suppliers, statuses } = store;

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

const materialColumns: QTableColumn[] = [
  { name: 'id', align: 'center', label: 'ID', field: 'id', sortable: true },
  {
    name: 'materialName',
    required: true,
    label: 'Material Name',
    align: 'left',
    field: 'name',
    sortable: true
  },
  { name: 'category', label: 'Category', field: 'category' },
  { name: 'supplier', label: 'Supplier', field: 'supplier' },
  { name: 'quantity', label: 'Quantity', field: 'quantity', sortable: true },
  { name: 'status', label: 'Status', field: 'status' },
  {
    name: 'actions',
    label: 'Actions',
    field: 'actions',
    align: 'center',
    sortable: false
  }
];

const showMaterial = ref(false)
const showAddDialog = ref(false)

const onMaterialRowClick: QTableProps['onRowClick'] = (evt, row) => {
  // Check if the click originated from the action button or its menu
  const target = evt.target as HTMLElement;
  if (target.closest('.action-button') || target.closest('.action-menu')) {
    return; // Do nothing if clicked on action button or its menu
  }
  selectedMaterial.value = row as MaterialRow
  showMaterial.value = true
}

function addMaterialToCart() {
  console.log('added material to cart', selectedMaterial.value.name)
  showMaterial.value = false
}

function openAddDialog() {
  newMaterial.value = {
    name: '',
    category: '',
    supplier: '',
    quantity: 0,
    status: 'Out of Stock',
    image: defaultImageUrl
  }
  imageUrlValid.value = true;
  showAddDialog.value = true
}

async function addNewMaterial() {
  try {
    // Validate image URL before proceeding
    if (!imageUrlValid.value) {
      operationNotifications.validation.error('Please provide a valid image URL');
      return;
    }

    // If image URL is empty, use default
    if (!newMaterial.value.image) {
      newMaterial.value.image = defaultImageUrl;
    }

    // Execute the store action and await its completion
    const result = await store.addMaterial(newMaterial.value);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showAddDialog.value = false;
      operationNotifications.add.success(`material: ${newMaterial.value.name}`);
    }
  } catch (error) {
    console.error('Error adding material:', error);
    operationNotifications.add.error('material');
  }
}

function applyFilters() {
  showFilterDialog.value = false;
  operationNotifications.filters.success();
}

// Add watch for quantity changes
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

// Function to validate if URL is a valid image
let currentAbortController: AbortController | null = null;

async function validateImageUrl(url: string): Promise<boolean> {
  if (!url) {
    imageUrlValid.value = false;
    return false;
  }

  if (!url.startsWith('http')) {
    imageUrlValid.value = false;
    return false;
  }

  validatingImage.value = true;

  // Abort any existing validation
  if (currentAbortController) {
    currentAbortController.abort();
  }

  // Create new abort controller for this validation
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

      // Handle abort signal
      signal.addEventListener('abort', () => {
        cleanup();
        resolve(false);
      });

      img.onload = () => {
        cleanup();
        imageUrlValid.value = true;
        validatingImage.value = false;
        resolve(true);
      };

      img.onerror = () => {
        cleanup();
        imageUrlValid.value = false;
        validatingImage.value = false;
        resolve(false);
      };

      // Set a timeout to avoid hanging
      const timeoutId = setTimeout(() => {
        if (!signal.aborted) {
          currentAbortController?.abort();
          imageUrlValid.value = false;
          validatingImage.value = false;
          resolve(false);
        }
      }, 5000);

      // Clean up timeout if aborted
      signal.addEventListener('abort', () => {
        clearTimeout(timeoutId);
      });

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

// Modify the watch for image URL changes to handle base64 validation
watch(() => newMaterial.value.image, async (newUrl: string) => {
  if (!newUrl || newUrl === defaultImageUrl) {
    imageUrlValid.value = true; // Default image or empty should be valid
    return;
  }
  try {
    if (newUrl.startsWith('data:image/')) {
      const validationResult = validateAndSanitizeBase64Image(newUrl);
      if (validationResult.isValid) {
        newMaterial.value.image = validationResult.sanitizedData!;
        imageUrlValid.value = true;
      } else {
        $q.notify({
          color: 'negative',
          message: validationResult.error || 'Invalid image data',
          position: 'top',
        });
        imageUrlValid.value = false;
      }
    } else {
      await validateImageUrl(newUrl);
    }
  } catch (error) {
    console.error('Error in image URL watcher:', error);
    imageUrlValid.value = false;
  }
});

// Add new refs for file handling
const fileInput = ref<HTMLInputElement | null>(null);
const previewUrl = ref('');
const isUploadingImage = ref(false);
const isDragging = ref(false);

// Add these constants at the top of the script
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif'] as const;
const MAX_DIMENSION = 4096; // Maximum image dimension in pixels

type AllowedMimeType = typeof ALLOWED_TYPES[number];

// Enhanced drag event handlers
function handleDragLeave(event: DragEvent) {
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
  const x = event.clientX;
  const y = event.clientY;

  if (x <= rect.left || x >= rect.right || y <= rect.top || y >= rect.bottom) {
    isDragging.value = false;
  }
}

// Update handleDrop function
function handleDrop(event: DragEvent) {
  event.preventDefault();
  isDragging.value = false;

  if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
    const file = event.dataTransfer.files[0];
    void handleFile(file);
  }
}

// Add clearImageInput function definition
function clearImageInput() {
  if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(previewUrl.value);
  }
  previewUrl.value = defaultImageUrl;
  newMaterial.value.image = defaultImageUrl;
  imageUrlValid.value = true;
  if (fileInput.value) {
    fileInput.value.value = '';
  }
  isUploadingImage.value = false;
}

// Update validateImageFile reference to use validateImageUrl
async function handleFile(file: File) {
  try {
    isUploadingImage.value = true;

    // Create a temporary URL for preview
    const tempPreviewUrl = URL.createObjectURL(file);
    previewUrl.value = tempPreviewUrl;

    const reader = new FileReader();
    reader.onload = async (e) => {
      if (e.target?.result) {
        const base64String = e.target.result as string;
        const base64ValidationResult = validateAndSanitizeBase64Image(base64String);

        if (!base64ValidationResult.isValid) {
          $q.notify({
            type: 'negative',
            message: base64ValidationResult.error || 'Invalid image data',
            position: 'top',
            timeout: 3000
          });
          clearImageInput();
          return;
        }

        newMaterial.value.image = base64ValidationResult.sanitizedData!;
        imageUrlValid.value = true;

        $q.notify({
          type: 'positive',
          message: 'Image uploaded successfully',
          position: 'top',
          timeout: 2000
        });
      }
    };

    reader.onerror = () => {
      clearImageInput();
      $q.notify({
        type: 'negative',
        message: 'Error reading file',
        position: 'top',
        timeout: 3000
      });
    };

    reader.readAsDataURL(file);
  } catch (error) {
    console.error('Error handling file:', error);
    clearImageInput();
    $q.notify({
      type: 'negative',
      message: 'An unexpected error occurred',
      position: 'top',
      timeout: 3000
    });
  } finally {
    isUploadingImage.value = false;
  }
}

// Update removeImage function to handle both Event and MouseEvent
function removeImage(event?: Event | MouseEvent) {
  if (event) {
    event.stopPropagation();
  }
  clearImageInput();
}

// Update handleFileSelect function
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

// Function to handle edit material
async function editMaterial(material: MaterialRow) {
  selectedMaterial.value = { ...material };
  newMaterial.value = {
    name: material.name,
    category: material.category,
    supplier: material.supplier,
    quantity: material.quantity,
    status: material.status,
    image: material.image
  };

  // Validate the image URL before setting preview
  if (material.image.startsWith('data:image/')) {
    const validationResult = validateAndSanitizeBase64Image(material.image);
    if (validationResult.isValid) {
      previewUrl.value = validationResult.sanitizedData!;
      newMaterial.value.image = validationResult.sanitizedData!;
      imageUrlValid.value = true;
    } else {
      // If invalid, use default image
      previewUrl.value = defaultImageUrl;
      newMaterial.value.image = defaultImageUrl;
      imageUrlValid.value = true;
      operationNotifications.validation.warning('Invalid image data, using default image');
    }
  } else {
    try {
      const isValid = await validateImageUrl(material.image);
      if (isValid) {
        previewUrl.value = material.image;
        imageUrlValid.value = true;
      } else {
        // If invalid, use default image
        previewUrl.value = defaultImageUrl;
        newMaterial.value.image = defaultImageUrl;
        imageUrlValid.value = true;
        operationNotifications.validation.warning('Invalid image URL, using default image');
      }
    } catch (error) {
      console.error('Error validating image URL:', error);
      previewUrl.value = defaultImageUrl;
      newMaterial.value.image = defaultImageUrl;
      imageUrlValid.value = true;
      operationNotifications.validation.warning('Error validating image, using default image');
    }
  }
  showEditDialog.value = true;
}

// Function to handle update material
async function updateMaterial() {
  try {
    // Validate image URL before proceeding
    if (!imageUrlValid.value) {
      operationNotifications.validation.error('Please provide a valid image URL');
      return;
    }

    // If image URL is empty, use default
    if (!newMaterial.value.image) {
      newMaterial.value.image = defaultImageUrl;
    }

    // Execute the store action and await its completion
    const result = await store.updateMaterial(selectedMaterial.value.id, newMaterial.value);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showEditDialog.value = false;
      clearImageInput();
      operationNotifications.update.success(`material: ${newMaterial.value.name}`);
    }
  } catch (error) {
    console.error('Error updating material:', error);
    operationNotifications.update.error('material');
  }
}

// Add new ref for delete dialog
const showDeleteDialog = ref(false);
const materialToDelete = ref<MaterialRow | null>(null);

// Function to handle delete material
function deleteMaterial(material: MaterialRow) {
  materialToDelete.value = material;
  showDeleteDialog.value = true;
}

// Function to confirm and execute delete
async function confirmDelete() {
  try {
    if (!materialToDelete.value) return;

    await store.deleteMaterial(materialToDelete.value.id);
    showDeleteDialog.value = false;
    materialToDelete.value = null;
    operationNotifications.delete.success('material');
  } catch (error) {
    console.error('Error deleting material:', error);
    operationNotifications.delete.error('material');
  }
}

// Add ref for edit dialog
const showEditDialog = ref(false);

// Update onMounted hook
onMounted(async () => {
  await store.initializeMaterials();
});
</script>

<template>
  <q-page class="flex inventory-page-padding">
    <div class="q-pa-sm full-width">
      <!-- Materials Section -->
      <div class="q-mt-sm">
        <div class="flex row q-my-sm">
          <div class="flex full-width col">
            <div class="flex col q-mr-sm">
              <q-input v-model="store.rawMaterialSearch" outlined dense placeholder="Search" class="full-width">
                <template v-slot:prepend>
                  <q-icon name="search" />
                </template>
              </q-input>
            </div>
            <div class="flex col">
              <q-btn outline icon="filter_list" label="Filters" @click="showFilterDialog = true" />
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

        <!--MATERIALS TABLE-->
        <q-table
          class="my-sticky-column-table"
          flat
          bordered
          title="Materials"
          :rows="store.filteredMaterialRows"
          :columns="materialColumns"
          row-key="id"
          :filter="store.materialSearch"
          @row-click="onMaterialRowClick"
          :pagination="{ rowsPerPage: 5 }"
          :loading="store.isLoading"
        >
          <template v-slot:loading>
            <q-inner-loading showing color="primary">
              <q-spinner-gears size="50px" color="primary" />
            </q-inner-loading>
          </template>
          <template v-slot:body-cell-actions="props">
            <q-td :props="props" auto-width>
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
                      @click.stop="editMaterial(props.row)"
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
                      @click.stop="deleteMaterial(props.row)"
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

        <!-- Existing Material Modal -->
        <ProductCardModal
          v-model="showMaterial"
          :image="selectedMaterial?.image || ''"
          :title="selectedMaterial?.name || ''"
          :price="0"
          :quantity="selectedMaterial?.quantity || 0"
          :details="`Supplier: ${selectedMaterial?.supplier}`"
          :unit_color="selectedMaterial?.category || ''"
          @addItem="addMaterialToCart"
        />

        <!-- Add Material Dialog -->
        <AddMaterialDialog v-model="showAddDialog" @add="addNewMaterial" :disable="!imageUrlValid || validatingImage">
          <q-input
            v-model="capitalizedName"
            label="Material Name"
            dense
            outlined
            :rules="[val => !!val || 'Field is required']"
          />
          <q-select
            v-model="newMaterial.category"
            label="Category"
            :options="categories"
            dense
            outlined
            :rules="[val => !!val || 'Field is required']"
          />
          <q-select
            v-model="newMaterial.supplier"
            label="Supplier"
            :options="suppliers"
            dense
            outlined
            :rules="[val => !!val || 'Field is required']"
          />
          <q-input
            v-model.number="newMaterial.quantity"
            label="Quantity"
            type="number"
            dense
            outlined
            :rules="[val => val >= 0 || 'Quantity cannot be negative']"
          />
          <q-input
            v-model="newMaterial.status"
            label="Status"
            dense
            outlined
            readonly
          />
          <!-- Enhanced Image Input Section -->
          <div class="column q-gutter-y-sm">
            <div class="row items-center q-gutter-x-sm">
              <q-input
                v-model="newMaterial.image"
                label="Image URL"
                dense
                outlined
                class="col"
                :rules="[val => imageUrlValid || 'Invalid image URL']"
                :loading="validatingImage"
                @blur="validateImageUrl(newMaterial.image)"
                clearable
                @clear="removeImage"
              >
                <template v-slot:prepend>
                  <q-icon name="link" />
                </template>
                <template v-slot:append>
                  <q-icon v-if="imageUrlValid && newMaterial.image && newMaterial.image !== defaultImageUrl" name="check_circle" color="positive" />
                  <q-icon v-else-if="!imageUrlValid" name="error" color="negative" />
                </template>
              </q-input>
              <q-btn
                icon="upload_file"
                flat
                round
                dense
                @click="triggerFileInput"
                title="Upload Image"
              />
            </div>
            <div
              class="upload-container q-pa-md"
              :class="{ 'dragging': isDragging }"
              @dragover.prevent
              @dragenter.prevent="isDragging = true"
              @dragleave.prevent="handleDragLeave"
              @drop.prevent="handleDrop"
            >
              <div class="text-center">
                <q-icon name="cloud_upload" size="lg" color="grey-7" />
                <p class="q-mt-sm text-grey-7">Drag & drop an image here, or click upload icon</p>
              </div>
              <input
                ref="fileInput"
                type="file"
                accept="image/png, image/jpeg, image/gif, image/webp"
                @change="handleFileSelect"
                style="display: none;"
              />
            </div>
            <!-- Image Preview and Remove Button -->
            <div v-if="newMaterial.image && imageUrlValid" class="row items-center justify-center q-mt-sm relative-position">
              <q-img
                :src="newMaterial.image"
                spinner-color="primary"
                style="max-height: 150px; max-width: 100%; border-radius: 4px;"
                alt="Image Preview"
                @error="imageUrlValid = false"
              >
                <template v-slot:error>
                  <div class="absolute-full flex flex-center bg-negative text-white">
                    Cannot load image
                  </div>
                </template>
              </q-img>
              <q-btn
                icon="close"
                flat
                round
                dense
                color="negative"
                size="sm"
                class="absolute-top-right q-ma-xs bg-white"
                @click="removeImage"
                style="z-index: 1; border-radius: 50%;"
              />
            </div>
          </div>
        </AddMaterialDialog>

        <!-- Filter Dialog -->
        <FilterMaterialDialog v-model="showFilterDialog">
          <!-- Default slot for filter options -->
          <q-select
            v-model="store.filterCategory"
            label="Filter by Category"
            :options="['All', ...categories]"
            dense
            outlined
            class="q-mb-sm"
          />
          <q-select
            v-model="store.filterSupplier"
            label="Filter by Supplier"
            :options="['All', ...suppliers]"
            dense
            outlined
            class="q-mb-sm"
          />
          <q-select
            v-model="store.filterStatus"
            label="Filter by Status"
            :options="['All', ...statuses]"
            dense
            outlined
          />
          <!-- Named slot for actions -->
          <template #actions>
            <q-btn flat label="Clear Filters" @click="store.resetFilters" />
            <q-btn flat label="Apply Filters" @click="applyFilters" />
          </template>
        </FilterMaterialDialog>

        <!-- Edit Material Dialog -->
        <EditMaterialDialog v-model="showEditDialog" @update="updateMaterial" :disable="!imageUrlValid || validatingImage">
          <q-input
            v-model="capitalizedName"
            label="Material Name"
            dense
            outlined
            :rules="[val => !!val || 'Field is required']"
          />
          <q-select
            v-model="selectedMaterial.category"
            label="Category"
            :options="categories"
            dense
            outlined
            :rules="[val => !!val || 'Field is required']"
          />
          <q-select
            v-model="selectedMaterial.supplier"
            label="Supplier"
            :options="suppliers"
            dense
            outlined
            :rules="[val => !!val || 'Field is required']"
          />
          <q-input
            v-model.number="selectedMaterial.quantity"
            label="Quantity"
            type="number"
            dense
            outlined
            :rules="[val => val >= 0 || 'Quantity cannot be negative']"
          />
          <q-select
            v-model="selectedMaterial.status"
            label="Status"
            :options="statuses"
            dense
            outlined
            :rules="[val => !!val || 'Field is required']"
          />
          <!-- Enhanced Image Input Section -->
          <div class="column q-gutter-y-sm">
            <div class="row items-center q-gutter-x-sm">
              <q-input
                v-model="selectedMaterial.image"
                label="Image URL"
                dense
                outlined
                class="col"
                :rules="[val => imageUrlValid || 'Invalid image URL']"
                :loading="validatingImage"
                @blur="validateImageUrl(selectedMaterial.image)"
                clearable
                @clear="removeImage"
              >
                <template v-slot:prepend>
                  <q-icon name="link" />
                </template>
                <template v-slot:append>
                  <q-icon v-if="imageUrlValid && selectedMaterial.image && selectedMaterial.image !== defaultImageUrl" name="check_circle" color="positive" />
                  <q-icon v-else-if="!imageUrlValid" name="error" color="negative" />
                </template>
              </q-input>
              <q-btn
                icon="upload_file"
                flat
                round
                dense
                @click="triggerFileInput"
                title="Upload Image"
              />
            </div>
            <div
              class="upload-container q-pa-md"
              :class="{ 'dragging': isDragging }"
              @dragover.prevent
              @dragenter.prevent="isDragging = true"
              @dragleave.prevent="handleDragLeave"
              @drop.prevent="handleDrop"
            >
              <div class="text-center">
                <q-icon name="cloud_upload" size="lg" color="grey-7" />
                <p class="q-mt-sm text-grey-7">Drag & drop an image here, or click upload icon</p>
              </div>
              <input
                ref="fileInput"
                type="file"
                accept="image/png, image/jpeg, image/gif, image/webp"
                @change="handleFileSelect"
                style="display: none;"
              />
            </div>
            <!-- Image Preview and Remove Button -->
            <div v-if="selectedMaterial.image && imageUrlValid" class="row items-center justify-center q-mt-sm relative-position">
              <q-img
                :src="selectedMaterial.image"
                spinner-color="primary"
                style="max-height: 150px; max-width: 100%; border-radius: 4px;"
                alt="Image Preview"
                @error="imageUrlValid = false"
              >
                <template v-slot:error>
                  <div class="absolute-full flex flex-center bg-negative text-white">
                    Cannot load image
                  </div>
                </template>
              </q-img>
              <q-btn
                icon="close"
                flat
                round
                dense
                color="negative"
                size="sm"
                class="absolute-top-right q-ma-xs bg-white"
                @click="removeImage"
                style="z-index: 1; border-radius: 50%;"
              />
            </div>
          </div>
        </EditMaterialDialog>

        <!-- Use the new Confirm Dialog Component -->
        <ConfirmDialog
          v-model="showDeleteDialog"
          title="Delete Material"
          :message="`Are you sure you want to delete ${materialToDelete?.name}? This action cannot be undone.`"
          confirmButtonLabel="Delete"
          confirmButtonColor="negative"
          @confirm="confirmDelete"
        />
      </div>
    </div>
  </q-page>
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

.z-top
  z-index: 1000

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

.action-button
  position: relative
  z-index: 1

.action-menu
  z-index: 1001 !important
</style>
