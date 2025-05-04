<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import ProductCardModal from 'src/components/Global/ProductModal.vue'
import { useQuasar } from 'quasar';
import { useMaterialsStore } from 'src/stores/materials';
<<<<<<< HEAD
import type { MaterialRow, NewMaterialInput } from 'src/stores/materials';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
=======
import type { MaterialRow } from 'src/stores/materials';
import { validateImageFile, validateAndSanitizeBase64Image } from '../utils/imageValidation';
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
import { operationNotifications } from '../utils/notifications';

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

<<<<<<< HEAD
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
=======
// Function to handle the file
function handleFile(file: File) {
  const validationResult = validateImageFile(file);
  if (!validationResult.isValid) {
    $q.notify({
      color: 'negative',
      message: validationResult.error || 'Invalid file',
      position: 'top',
    });
    return;
  }

  const reader = new FileReader();
  reader.onload = (e) => {
    if (e.target?.result) {
      const base64String = e.target.result as string;
      const base64ValidationResult = validateAndSanitizeBase64Image(base64String);
      
      if (!base64ValidationResult.isValid) {
        $q.notify({
          color: 'negative',
          message: base64ValidationResult.error || 'Invalid image data',
          position: 'top',
        });
        return;
      }

      previewUrl.value = base64ValidationResult.sanitizedData!;
      newMaterial.value.image = base64ValidationResult.sanitizedData!;
      imageUrlValid.value = true;
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
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

<<<<<<< HEAD
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
=======
// Function to remove image
function removeImage(event: MouseEvent) {
  event.stopPropagation(); // Prevent triggering file input click
  previewUrl.value = '';
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  newMaterial.value.image = defaultImageUrl;
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
<<<<<<< HEAD
          :loading="store.isLoading"
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
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

        <!-- Add Material Dialog - Minimalistic Design -->
        <q-dialog
          v-model="showAddDialog"
          persistent
          @hide="clearImageInput"
        >
          <q-card style="min-width: 400px; max-width: 95vw">
            <q-card-section class="row items-center q-pb-none">
              <div class="text-h6">New Material</div>
              <q-space />
              <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>

            <q-card-section>
              <q-form @submit.prevent="addNewMaterial" class="q-gutter-sm">
                <q-input
                  v-model="capitalizedName"
                  label="Material Name"
                  dense
                  outlined
                  required
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
                      :rules="[val => !!val || 'Category is required']"
                    >
                      <template v-slot:prepend>
                        <q-icon name="category" />
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
                      :rules="[val => !!val || 'Supplier is required']"
                    >
                      <template v-slot:prepend>
                        <q-icon name="local_shipping" />
                      </template>
                    </q-select>
                  </div>
                </div>

                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-6">
                    <q-input
                      v-model.number="newMaterial.quantity"
                      type="number"
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
                      v-model="newMaterial.status"
                      label="Status"
                      dense
                      outlined
                      readonly
                      disable
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
                      class="upload-container q-pa-md"
                      :class="{ 'dragging': isDragging }"
                      @dragenter.prevent="isDragging = true"
                      @dragover.prevent="isDragging = true"
                      @dragleave.prevent="handleDragLeave"
                      @drop.prevent="handleDrop"
                      @click="triggerFileInput"
                    >
                      <input
                        type="file"
                        ref="fileInput"
                        accept="image/*"
                        class="hidden"
                        @change="handleFileSelect"
                      >
                      <div class="text-center" v-if="!previewUrl">
                        <q-icon name="cloud_upload" size="48px" color="primary" />
                        <div class="text-body1 q-mt-sm">
                          Drag and drop an image here or click to select
                        </div>
                        <div class="text-caption text-grey">
                          Supported formats: JPG, PNG, GIF
                        </div>
                      </div>
                      <div v-else class="row items-center">
                        <div class="col-8">
                          <img :src="previewUrl" class="preview-image" />
                        </div>
                        <div class="col-4 text-center">
                          <q-btn
                            flat
                            round
                            color="negative"
                            icon="close"
                            @mousedown.stop="removeImage($event)"
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
                label="Add Material"
                @click="addNewMaterial"
                :disable="!newMaterial.name || !newMaterial.category || !newMaterial.supplier || newMaterial.quantity < 0"
              />
            </q-card-actions>
          </q-card>
        </q-dialog>

        <!-- Filter Dialog -->
        <q-dialog v-model="showFilterDialog">
          <q-card style="min-width: 350px">
            <q-card-section class="row items-center">
              <div class="text-h6">Filter Materials</div>
              <q-space />
              <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>

            <q-card-section class="q-pt-none">
              <q-select
                v-model="store.filterCategory"
                :options="categories"
                label="Category"
                clearable
                outlined
                class="q-mb-md"
              />

              <q-select
                v-model="store.filterSupplier"
                :options="suppliers"
                label="Supplier"
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

        <!-- Edit Material Dialog -->
        <q-dialog
          v-model="showEditDialog"
          persistent
          @hide="clearImageInput"
        >
          <q-card style="min-width: 400px; max-width: 95vw">
            <q-card-section class="row items-center q-pb-none">
              <div class="text-h6">Edit Material</div>
              <q-space />
              <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>

            <q-card-section>
              <q-form @submit.prevent="updateMaterial" class="q-gutter-sm">
                <q-input
                  v-model="capitalizedName"
                  label="Material Name"
                  dense
                  outlined
                  required
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
                      :rules="[val => !!val || 'Category is required']"
                    >
                      <template v-slot:prepend>
                        <q-icon name="category" />
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
                      :rules="[val => !!val || 'Supplier is required']"
                    >
                      <template v-slot:prepend>
                        <q-icon name="local_shipping" />
                      </template>
                    </q-select>
                  </div>
                </div>

                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-6">
                    <q-input
                      v-model.number="newMaterial.quantity"
                      type="number"
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
                      v-model="newMaterial.status"
                      label="Status"
                      dense
                      outlined
                      readonly
                      disable
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
                      class="upload-container q-pa-md"
                      :class="{ 'dragging': isDragging }"
                      @dragenter.prevent="isDragging = true"
                      @dragover.prevent="isDragging = true"
                      @dragleave.prevent="handleDragLeave"
                      @drop.prevent="handleDrop"
                      @click="triggerFileInput"
                    >
                      <input
                        type="file"
                        ref="fileInput"
                        accept="image/*"
                        class="hidden"
                        @change="handleFileSelect"
                      >
                      <div v-if="!previewUrl && !newMaterial.image" class="text-center">
                        <q-icon name="cloud_upload" size="48px" color="primary" />
                        <div class="text-body1 q-mt-sm">
                          Drag and drop an image here or click to select
                        </div>
                        <div class="text-caption text-grey">
                          Supported formats: JPG, PNG, GIF
                        </div>
                      </div>
                      <div v-else class="row items-center">
                        <div class="col-8">
                          <img :src="previewUrl || newMaterial.image" class="preview-image" />
                        </div>
                        <div class="col-4 text-center">
                          <q-btn
                            flat
                            round
                            color="negative"
                            icon="close"
                            @mousedown.stop="removeImage($event)"
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
                label="Update Material"
                @click="updateMaterial"
                :disable="!newMaterial.name || !newMaterial.category || !newMaterial.supplier || newMaterial.quantity < 0"
              />
            </q-card-actions>
          </q-card>
        </q-dialog>

        <!-- Delete Confirmation Dialog -->
        <q-dialog v-model="showDeleteDialog" persistent>
          <q-card>
            <q-card-section class="row items-center">
              <q-avatar icon="warning" color="negative" text-color="white" />
              <span class="q-ml-sm text-h6">Delete Material</span>
            </q-card-section>

            <q-card-section>
              Are you sure you want to delete {{ materialToDelete?.name }}? This action cannot be undone.
            </q-card-section>

            <q-card-actions align="right">
              <q-btn flat label="Cancel" v-close-popup />
              <q-btn flat label="Delete" color="negative" @click="confirmDelete" />
            </q-card-actions>
          </q-card>
        </q-dialog>
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
