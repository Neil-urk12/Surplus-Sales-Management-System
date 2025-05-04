<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import ProductCardModal from 'src/components/Global/ProductModal.vue'
import { useQuasar } from 'quasar';
import { useMaterialsStore } from 'src/stores/materials';
import type { MaterialRow } from 'src/stores/materials';
import { validateImageFile, validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';

const $q = useQuasar();
const store = useMaterialsStore();
const showFilterDialog = ref(false);
const selectedMaterial = ref<MaterialRow>({
  name: '',
  id: 0,
  category: '',
  supplier: '',
  quantity: 0,
  status: '',
  image: ''
})

const newMaterial = ref<Omit<MaterialRow, 'id'>>({
  name: '',
  category: '',
  supplier: '',
  quantity: 0,
  status: '',
  image: ''
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
const isDragging = ref(false);
const previewUrl = ref('');

// Function to handle file selection
function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files[0]) {
    const file = input.files[0];
    handleFile(file);
  }
}

// Function to handle drag and drop
function handleDrop(event: DragEvent) {
  event.preventDefault();
  isDragging.value = false;

  if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
    const file = event.dataTransfer.files[0];
    handleFile(file);
  }
}

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
    }
  };
  reader.readAsDataURL(file);
}

// Function to remove image
function removeImage(event: MouseEvent) {
  event.stopPropagation(); // Prevent triggering file input click
  previewUrl.value = '';
  newMaterial.value.image = defaultImageUrl;
  if (fileInput.value) {
    fileInput.value.value = ''; // Clear the file input
  }
}

// Function to trigger file input
function triggerFileInput() {
  fileInput.value?.click();
}

// Function to clear image input and preview
function clearImageInput() {
  previewUrl.value = '';
  newMaterial.value.image = defaultImageUrl;
  if (fileInput.value) {
    fileInput.value.value = '';
  }
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
        >
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
                      @dragleave.prevent="isDragging = false"
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
                      @dragleave.prevent="isDragging = false"
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
