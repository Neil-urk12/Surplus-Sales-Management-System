<script setup lang="ts">
import { ref, watch, computed, onMounted, defineAsyncComponent } from 'vue';
import type { QTableColumn } from 'quasar';
import ProductCardModal from 'src/components/Global/ProductModal.vue'
import { useQuasar } from 'quasar';
import { useMaterialsStore } from 'src/stores/materials';
import type { MaterialRow, NewMaterialInput } from 'src/stores/materials';
import type { UpdateMaterialInput } from 'src/types/materials';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';
const DeleteDialog = defineAsyncComponent(() => import('src/components/Global/DeleteDialog.vue'));
const AddMaterialDialog = defineAsyncComponent(() => import('../components/AddMaterialDialog.vue'))
const EditMaterialDialog = defineAsyncComponent(() => import('../components/EditMaterialDialog.vue'))
const FilterMaterialDialog = defineAsyncComponent(() => import('../components/FilterMaterialDialog.vue'))

const $q = useQuasar();
const store = useMaterialsStore();
console.log('Store initialized:', {
  rawMaterialSearch: store.rawMaterialSearch,
  filterCategory: store.filterCategory,
  filterSupplier: store.filterSupplier,
  filterStatus: store.filterStatus,
  pagination: store.pagination
});
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
  {
    name: 'status',
    label: 'Status',
    field: 'status',
    align: 'center'
  },
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

// Watch for filter changes and trigger data refresh
watch([store.filterCategory, store.filterSupplier, store.filterStatus], () => {
  void store.onRequest({ pagination: store.pagination });
});

// Watch for raw search input changes
watch(() => store.rawMaterialSearch, () => {
  void store.onRequest({ pagination: store.pagination });
});

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

    // If image URL is empty, use default (only relevant if image is part of update)
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
  void store.onRequest({ pagination: store.pagination });
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

// Function to handle edit material
function editMaterial(material: MaterialRow) {
  // Deep copy selected material to the dedicated edit ref
  materialToEdit.value = JSON.parse(JSON.stringify(material));
  imageUrlValid.value = true; // Reset validation state for the dialog
  validatingImage.value = false; // Reset validation state
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

    // If image URL is empty, use default (only relevant if image is part of update)
    if (!materialToEdit.value.image) {
      materialToEdit.value.image = defaultImageUrl;
    }

    // Prepare the update payload from materialToEdit
    const updatePayload: UpdateMaterialInput = {
      name: materialToEdit.value.name,
      category: materialToEdit.value.category,
      supplier: materialToEdit.value.supplier,
      quantity: materialToEdit.value.quantity,
      status: materialToEdit.value.status,
      image: materialToEdit.value.image
    };

    // Execute the store action and await its completion
    const result = await store.updateMaterial(materialToEdit.value.id, updatePayload);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showEditDialog.value = false;
      operationNotifications.update.success(`material: ${materialToEdit.value.name}`);
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

// Add a computed property for status color mapping
const getStatusColor = (status: string) => {
  switch (status) {
    case 'In Stock':
      return 'positive';
    case 'Low Stock':
      return 'warning';
    case 'Out of Stock':
      return 'negative';
    default:
      return 'grey';
  }
};

// Update onMounted hook
onMounted(async () => {
  console.log('MaterialsPage mounted, initializing materials...');
  await store.initializeMaterials();
  console.log('Materials initialized:', {
    rowCount: store.materialRows.length,
    pagination: store.pagination
  });
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
              <q-input
                v-model="store.rawMaterialSearch"
                outlined
                dense
                placeholder="Search by name, ID, category, or supplier..."
                class="full-width"
                clearable
                debounce="300"
                @clear="() => {
                  store.rawMaterialSearch = '';
                  store.resetFilters();
                  void store.onRequest({ pagination: store.pagination });
                }"
                @update:model-value="() => {
                  console.log('Search input updated:', store.rawMaterialSearch);
                  void store.onRequest({ pagination: store.pagination });
                }"
              >
                <template v-slot:prepend>
                  <q-icon name="search" />
                </template>
                <template v-slot:hint>
                  Type to search all materials
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
          :rows="store.materialRows"
          :columns="materialColumns"
          row-key="id"
          :pagination="store.pagination"
          @request="store.onRequest"
          :loading="store.isLoading"
          binary-state-sort
          :rows-per-page-options="[10, 20, 50, 0]"
        >
          <template v-slot:loading>
            <q-inner-loading showing color="primary">
              <q-spinner-gears size="50px" color="primary" />
            </q-inner-loading>
          </template>

          <template v-slot:no-data>
            <div class="full-width row flex-center q-pa-md">
              <q-icon name="search_off" size="2em" color="grey-7" class="q-mr-sm" />
              <span class="text-grey-7">No materials found matching your search</span>
            </div>
          </template>

          <template v-slot:bottom>
            <div class="row items-center justify-end q-mt-md">
              <q-pagination
                :model-value="store.pagination?.page || 1"
                :max="store.pagination ? Math.ceil((store.pagination.rowsNumber || 0) / (store.pagination.rowsPerPage || 10)) : 1"
                :max-pages="6"
                boundary-numbers
                direction-links
                flat
                color="primary"
                @update:model-value="(val) => store.onRequest({ pagination: { ...store.pagination, page: val } })"
              />
            </div>
          </template>

          <template v-slot:body-cell-status="props">
            <q-td :props="props">
              <q-chip
                dense
                :color="getStatusColor(props.value)"
                text-color="white"
                class="q-px-sm"
              >
                {{ props.value }}
              </q-chip>
            </q-td>
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
            v-model="materialToEdit.name"
            label="Material Name"
            dense
            outlined
            :rules="[val => !!val || 'Field is required']"
          />
          <q-select
            v-model="materialToEdit.category"
            label="Category"
            :options="categories"
            dense
            outlined
            :rules="[val => !!val || 'Field is required']"
          />
          <q-select
            v-model="materialToEdit.supplier"
            label="Supplier"
            :options="suppliers"
            dense
            outlined
            :rules="[val => !!val || 'Field is required']"
          />
          <q-input
            v-model.number="materialToEdit.quantity"
            label="Quantity"
            type="number"
            dense
            outlined
            :rules="[val => val >= 0 || 'Quantity cannot be negative']"
          />
          <q-select
            v-model="materialToEdit.status"
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
                v-model="materialToEdit.image"
                label="Image URL (Optional)"
                dense
                outlined
                class="col"
                :rules="[val => imageUrlValid || 'Invalid image URL']"
                :loading="validatingImage"
                @blur="validateImageUrl(materialToEdit.image)"
                clearable
                @clear="removeImage"
              >
                <template v-slot:prepend>
                  <q-icon name="link" />
                </template>
                <template v-slot:append>
                  <q-icon v-if="imageUrlValid && materialToEdit.image && materialToEdit.image !== defaultImageUrl" name="check_circle" color="positive" />
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
            <div v-if="materialToEdit.image && imageUrlValid" class="row items-center justify-center q-mt-sm relative-position">
              <q-img
                :src="materialToEdit.image"
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

        <DeleteDialog
          v-model="showDeleteDialog"
          itemType="material"
          :itemName="materialToDelete?.name || ''"
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
</style>
