<script setup lang="ts">
import { ref, watch, computed, onMounted, defineAsyncComponent } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import { useQuasar } from 'quasar';
import { useMaterialsStore } from 'src/stores/materials';
import type { MaterialRow, NewMaterialInput } from 'src/stores/materials';
import type { MaterialCategory, MaterialSupplier, MaterialStatus } from 'src/types/materials';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';
const ProductCardModal = defineAsyncComponent(() => import('src/components/Global/ProductModal.vue'));
const DeleteDialog = defineAsyncComponent(() => import('src/components/Global/DeleteDialog.vue'));
const AddMaterialDialog = defineAsyncComponent(() => import('src/components/Materials/AddMaterialDialog.vue'));
const MaterialEditWrapper = defineAsyncComponent(() => import('src/components/Materials/MaterialEditWrapper.vue'));
const FilterMaterialDialog = defineAsyncComponent(() => import('src/components/Materials/FilterMaterialDialog.vue'));
const AdvancedSearch = defineAsyncComponent(() => import('src/components/Global/AdvancedSearch.vue'));

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

const materialToEdit = ref<MaterialRow>({
  id: 0,
  name: '',
  category: 'Building',
  supplier: 'Steel Co.',
  quantity: 0,
  status: 'Out of Stock',
  image: ''
});

// Image validation
const imageUrlValid = ref(true);
const validatingImage = ref(false);
const defaultImageUrl = 'https://loremflickr.com/600/400/material';

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

async function handleAddMaterial(materialData: NewMaterialInput) {
  try {
    const result = await store.addMaterial(materialData);
    if (result.success) {
      showAddDialog.value = false;
      operationNotifications.add.success(`material: ${materialData.name}`);
    }
  } catch (error) {
    console.error('Error adding material:', error);
    operationNotifications.add.error('material');
  }
}

// Add a dedicated function to close the edit dialog
function closeEditDialog() {
  console.log('closeEditDialog called in MaterialsPage');
  showEditDialog.value = false;
  materialToEdit.value = {
    id: 0,
    name: '',
    category: 'Building',
    supplier: 'Steel Co.',
    quantity: 0,
    status: 'Out of Stock',
    image: ''
  };
}

async function handleUpdateMaterial(materialData: NewMaterialInput) {
  try {
    console.log('handleUpdateMaterial called in MaterialsPage');
    if (!materialToEdit.value || !materialToEdit.value.id) {
      throw new Error('No material selected for update or missing ID');
    }

    const result = await store.updateMaterial(materialToEdit.value.id, materialData);
    console.log('Update result:', result);
    if (result.success) {
      closeEditDialog(); // Use the dedicated function
      operationNotifications.update.success(`material: ${materialData.name}`);
    }
  } catch (error) {
    console.error('Error updating material:', error);
    operationNotifications.update.error('material');
  }
}

function editMaterial(material: MaterialRow) {
  materialToEdit.value = { ...material };
  showEditDialog.value = true;
}

// Add watch for quantity changes
watch(() => newMaterial.value.quantity, (newQuantity) => {
  if (newQuantity === 0) {
    newMaterial.value.status = 'Out of Stock';
  } else if (newQuantity <= 10) {
    newMaterial.value.status = 'Low Stock';
  } else {
    newMaterial.value.status = 'In Stock';
  }
});

// Watcher for materialToEdit quantity
watch(() => materialToEdit.value.quantity, (newQuantity) => {
  if (newQuantity === 0) {
    materialToEdit.value.status = 'Out of Stock';
  } else if (newQuantity <= 10) {
    materialToEdit.value.status = 'Low Stock';
  } else {
    materialToEdit.value.status = 'In Stock';
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
  // Check if the mouse left the container (not just moved between child elements)
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
  const x = event.clientX;
  const y = event.clientY;

  // Check if the mouse is outside the container's bounds
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

// Update handleFile function
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
          previewUrl.value = defaultImageUrl;
          return;
        }

        console.log('Base64 validation passed, updating image');
        if (showEditDialog.value) {
          materialToEdit.value.image = base64ValidationResult.sanitizedData!;
        } else {
          newMaterial.value.image = base64ValidationResult.sanitizedData!;
        }
        imageUrlValid.value = true;

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
      previewUrl.value = defaultImageUrl;
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
    previewUrl.value = defaultImageUrl;
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

// Update validateFile function
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

// Add validateImageDimensions function
function validateImageDimensions(file: File): Promise<{ isValid: boolean; error?: string }> {
  return new Promise((resolve) => {
    const img = new Image();
    const objectUrl = URL.createObjectURL(file);

    const cleanup = () => {
      URL.revokeObjectURL(objectUrl);
    };

    img.onload = () => {
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
      cleanup();
      console.error('Error loading image for dimension validation');
      resolve({
        isValid: false,
        error: 'Error loading image. Please ensure it is a valid image file.'
      });
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

    img.src = objectUrl;

    // Clear timeout when image loads or errors
    img.onload = () => {
      clearTimeout(timeout);
      cleanup();
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
  });
}

// Update removeImage function
function removeImage(event?: Event) {
  if (event) {
    event.stopPropagation();
  }
  clearImageInput();
}

// Update clearImageInput function
function clearImageInput() {
  if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(previewUrl.value);
  }
  previewUrl.value = defaultImageUrl;
  if (showEditDialog.value) {
    materialToEdit.value.image = defaultImageUrl;
  } else {
    newMaterial.value.image = defaultImageUrl;
  }
  imageUrlValid.value = true;
  if (fileInput.value) {
    fileInput.value.value = '';
  }
  isUploadingImage.value = false;
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

// Add new ref for delete dialog
const showDeleteDialog = ref(false);
const materialToDelete = ref<MaterialRow>({
  id: 0,
  name: '',
  category: 'Building',
  supplier: 'Steel Co.',
  quantity: 0,
  status: 'Out of Stock',
  image: ''
});

// Function to handle delete material
function deleteMaterial(material: MaterialRow) {
  materialToDelete.value = { ...material };
  showDeleteDialog.value = true;
}

// Function to confirm and execute delete
async function confirmDelete() {
  try {
    if (!materialToDelete.value || materialToDelete.value.id === 0) return;

    await store.deleteMaterial(materialToDelete.value.id);
    showDeleteDialog.value = false;
    materialToDelete.value = {
      id: 0,
      name: '',
      category: 'Building',
      supplier: 'Steel Co.',
      quantity: 0,
      status: 'Out of Stock',
      image: ''
    };
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

function handleApplyFilters(filters: { category: string | null; supplier: string | null; status: string | null }) {
  store.filterCategory = filters.category === null ? '' : filters.category as MaterialCategory;
  store.filterSupplier = filters.supplier === null ? '' : filters.supplier as MaterialSupplier;
  store.filterStatus = filters.status === null ? '' : filters.status as MaterialStatus;
  operationNotifications.filters.success();
  showFilterDialog.value = false;
}
</script>

<template>
  <q-page class="flex inventory-page-padding">
    <div class="q-pa-sm full-width">
      <!-- Modified Header Structure -->
      <div class="q-mb-md">
        <div class="flex row items-center justify-between">
          <div class="col">
            <div class="text-h5">Materials</div>
          </div>
        </div>
        <div>
          <div class="text-caption text-grey q-mt-sm">Manage your inventory items, track stock levels, and monitor product details.</div>
          <!-- Main Controls Container -->
          <div
            class="flex items-center q-mt-sm"
            :class="$q.screen.lt.md ? 'column q-gutter-y-sm items-stretch' : 'row justify-between'"
          >
            <!-- Search + Filters Group -->
            <div
              class="flex items-center"
              :class="$q.screen.lt.md ? 'column full-width q-gutter-y-sm items-stretch' : 'row q-gutter-x-sm'"
            >
              <AdvancedSearch
                v-model="store.search.searchInput"
                placeholder="Search materials"
                @clear="store.resetFilters"
                color="primary"
                :disable="store.isLoading"
                :style="$q.screen.lt.md ? { width: '100%' } : { width: '400px' }"
              />
              <q-btn
                outline
                icon="filter_list"
                label="Filters"
                @click="showFilterDialog = true"
                :disable="store.isLoading"
                :class="{ 'full-width': $q.screen.lt.md }"
              />
            </div>
            <!-- Add + Download CSV Group -->
            <div
              class="flex items-center"
              :class="$q.screen.lt.md ? 'column full-width q-gutter-y-sm items-stretch' : 'row q-gutter-x-sm'"
            >
              <q-btn
                unelevated
                @click="openAddDialog"
                :disable="store.isLoading"
                :class="['text-white bg-primary', { 'full-width': $q.screen.lt.md }]"
              >
                <q-icon name="add" color="white" />
                Add
              </q-btn>
              <q-btn
                dense
                flat
                :disable="store.isLoading"
                :class="['bg-primary text-white q-pa-sm', { 'full-width': $q.screen.lt.md }]"
              >
                <q-icon name="download" color="white" />
                Download CSV
              </q-btn>
            </div>
          </div>
        </div>
      </div>

      <!-- Materials Section -->
      <div class="q-mt-sm">
        <!--MATERIALS TABLE-->
        <q-table class="my-sticky-column-table custom-table-text" flat bordered :rows="store.filteredMaterialRows" 
          :columns="materialColumns" row-key="id" :filter="store.search.searchValue" @row-click="onMaterialRowClick"
          :pagination="{ rowsPerPage: 10 }" :rows-per-page-options="[10]" :loading="store.isLoading">
          <template v-slot:loading>
            <q-inner-loading showing color="primary">
              <q-spinner-gears size="50px" color="primary" />
            </q-inner-loading>
          </template>
          <template v-slot:body-cell-status="props">
            <q-td :props="props">
              <q-badge :color="props.row.status === 'In Stock' ? 'green' : (props.row.status === 'Out of Stock' || props.row.status === 'Low Stock' ? 'red' : 'grey')" :label="props.row.status" />
            </q-td>
          </template>
          <template v-slot:body-cell-actions="props">
            <q-td :props="props" auto-width>
              <q-btn flat round dense color="grey" icon="more_vert" class="action-button"
                :aria-label="'Actions for ' + props.row.name">
                <q-menu class="action-menu" :aria-label="'Available actions for ' + props.row.name">
                  <q-list style="min-width: 100px">
                    <q-item clickable v-close-popup @click.stop="editMaterial(props.row)" role="button"
                      :aria-label="'Edit ' + props.row.name">
                      <q-item-section>
                        <q-item-label>
                          <q-icon name="edit" size="xs" class="q-mr-sm" aria-hidden="true" />
                          Edit
                        </q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-item clickable v-close-popup @click.stop="deleteMaterial(props.row)" role="button"
                      :aria-label="'Delete ' + props.row.name" class="text-negative">
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
        <ProductCardModal v-model="showMaterial" :image="selectedMaterial?.image || ''"
          :title="selectedMaterial?.name || ''" :price="0" :quantity="selectedMaterial?.quantity || 0"
          :details="`Supplier: ${selectedMaterial?.supplier}`" :unit_color="selectedMaterial?.category || ''"
          @addItem="addMaterialToCart" />

        <!-- Add Material Dialog -->
        <AddMaterialDialog
          v-model="showAddDialog"
          :categories="categories"
          :suppliers="suppliers"
          :default-image-url="defaultImageUrl"
          :material-name="capitalizedName"
          @add-material="handleAddMaterial"
        />

        <!-- Image Upload Zone (hidden but functional) -->
        <div class="hidden">
          <input
            ref="fileInput"
            type="file"
            accept="image/jpeg,image/png,image/gif"
            @change="handleFileSelect"
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

        <!-- Filter Dialog -->
        <FilterMaterialDialog v-model="showFilterDialog" :categories="categories" :suppliers="suppliers"
          :statuses="statuses" :initial-filter-category="store.filterCategory === '' ? null : store.filterCategory"
          :initial-filter-supplier="store.filterSupplier === '' ? null : store.filterSupplier"
          :initial-filter-status="store.filterStatus === '' ? null : store.filterStatus"
          @apply-filters="handleApplyFilters" @reset-filters="store.resetFilters" />

        <!-- Edit Material Dialog -->
        <MaterialEditWrapper
          v-model:open="showEditDialog"
          :material-data="materialToEdit"
          :categories="categories"
          :suppliers="suppliers"
          :default-image-url="defaultImageUrl"
          @update-material="handleUpdateMaterial"
        />

        <DeleteDialog 
          v-model="showDeleteDialog" 
          itemType="material" 
          :itemName="materialToDelete.name" 
          @confirm-delete="confirmDelete" 
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

.custom-table-text
  td,
  th
    font-size: 1.15em
    font-weight: 400

    .q-badge
      font-size: 0.9em
      font-weight: 600
</style>
