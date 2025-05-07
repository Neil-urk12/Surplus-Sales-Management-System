<script setup lang="ts">
import { ref, watch, computed, onMounted, defineAsyncComponent } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import { useQuasar } from 'quasar';
import { useMaterialsStore } from 'src/stores/materials';
import type { MaterialRow, NewMaterialInput } from 'src/stores/materials';
import type { UpdateMaterialInput, MaterialCategoryInput, MaterialSupplierInput, MaterialStatus } from 'src/types/materials';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';
const ProductCardModal = defineAsyncComponent(() => import('src/components/Global/ProductModal.vue'));
const DeleteDialog = defineAsyncComponent(() => import('src/components/Global/DeleteDialog.vue'));
const AddMaterialDialog = defineAsyncComponent(() => import('../components/AddMaterialDialog.vue'))
const EditMaterialDialog = defineAsyncComponent(() => import('../components/EditMaterialDialog.vue'))
const FilterMaterialDialog = defineAsyncComponent(() => import('../components/FilterMaterialDialog.vue'))
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

/* eslint-disable @typescript-eslint/no-unused-vars */
/**
 * Computed property for capitalizing material names.
 * This is currently not actively used in the UI but is kept
 * for potential future use in the application.
 */
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
/* eslint-enable @typescript-eslint/no-unused-vars */

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
  };
  imageUrlValid.value = true;
  showAddDialog.value = true;
}

// Add function to refresh materials list
async function refreshMaterials() {
  console.log('Refreshing materials list');
  await store.initializeMaterials();
}

async function addNewMaterial(materialInput: NewMaterialInput) {
  try {
    // Validate image URL before proceeding
    if (!imageUrlValid.value) {
      operationNotifications.validation.error('Please provide a valid image URL');
      return;
    }

    // If image URL is empty, use default (only relevant if image is part of update)
    if (!materialInput.image) {
      materialInput.image = defaultImageUrl;
    }

    // Execute the store action and await its completion
    const result = await store.addMaterial(materialInput);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showAddDialog.value = false;
      operationNotifications.add.success(`material: ${materialInput.name}`);
      // Refresh the materials list to ensure we have the latest data
      await refreshMaterials();
    }
  } catch (error) {
    console.error('Error adding material:', error);
    operationNotifications.add.error('material');
  }
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

// Function to handle edit material
function editMaterial(material: MaterialRow) {
  // Deep copy selected material to the dedicated edit ref
  materialToEdit.value = JSON.parse(JSON.stringify(material));
  imageUrlValid.value = true; // Reset validation state for the dialog
  validatingImage.value = false; // Reset validation state
  showEditDialog.value = true;
}

// Function to handle update material
async function updateMaterial(updatedMaterial: MaterialRow) {
  try {
    // Validate image URL before proceeding
    if (!imageUrlValid.value) {
      operationNotifications.validation.error('Please provide a valid image URL');
      return;
    }

    // If image URL is empty, use default (only relevant if image is part of update)
    if (!updatedMaterial.image) {
      updatedMaterial.image = defaultImageUrl;
    }

    // Prepare the update payload from the received updatedMaterial
    const updatePayload: UpdateMaterialInput = {
      name: updatedMaterial.name,
      category: updatedMaterial.category,
      supplier: updatedMaterial.supplier,
      quantity: updatedMaterial.quantity,
      status: updatedMaterial.status,
      image: updatedMaterial.image
    };

    console.log('Sending update for material:', updatedMaterial.id, updatePayload);

    // Execute the store action and await its completion
    const result = await store.updateMaterial(updatedMaterial.id, updatePayload);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showEditDialog.value = false;
      operationNotifications.update.success(`material: ${updatedMaterial.name}`);
      // Refresh the materials list to ensure we have the latest data
      await refreshMaterials();
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

    const result = await store.deleteMaterial(materialToDelete.value.id);
    
    if (result.success) {
      showDeleteDialog.value = false;
      materialToDelete.value = null;
      operationNotifications.delete.success('material');
      // Refresh the materials list to ensure we have the latest data
      await refreshMaterials();
    }
  } catch (error) {
    console.error('Error deleting material:', error);
    operationNotifications.delete.error('material');
  }
}

// Add ref for edit dialog
const showEditDialog = ref(false);

// Function to handle applying filters from the FilterMaterialDialog
function handleApplyFilters(filterData: { category: string | null; supplier: string | null; status: string | null }) {
  // Type assertion is safe here because we're ensuring the values match the expected types from the dialog
  store.filterCategory = (filterData.category || '') as MaterialCategoryInput;
  store.filterSupplier = (filterData.supplier || '') as MaterialSupplierInput;
  store.filterStatus = (filterData.status || '') as MaterialStatus | '';
  showFilterDialog.value = false;
}

// Update onMounted hook
onMounted(async () => {
  console.log('MaterialsPage mounted, initializing materials');
  await refreshMaterials();
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
              <AdvancedSearch v-model="store.search.searchInput" placeholder="Search materials"
                @clear="store.resetFilters" color="primary" />
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
        <q-table class="my-sticky-column-table" flat bordered title="Materials" :rows="store.filteredMaterialRows"
          :columns="materialColumns" row-key="id" :filter="store.search.searchValue" @row-click="onMaterialRowClick"
          :pagination="{ rowsPerPage: 5 }" :loading="store.isLoading">
          <template v-slot:loading>
            <q-inner-loading showing color="primary">
              <q-spinner-gears size="50px" color="primary" />
            </q-inner-loading>
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
          @add="addNewMaterial" 
          :disable="!imageUrlValid || validatingImage"
          :materialData="newMaterial"
          :categories="categories"
          :suppliers="suppliers"
          :defaultImageUrl="defaultImageUrl"
        />

        <!-- Filter Dialog -->
        <FilterMaterialDialog v-model="showFilterDialog" :categories="categories" :suppliers="suppliers"
          :statuses="statuses" :initial-filter-category="store.filterCategory === '' ? null : store.filterCategory"
          :initial-filter-supplier="store.filterSupplier === '' ? null : store.filterSupplier"
          :initial-filter-status="store.filterStatus === '' ? null : store.filterStatus"
          @apply-filters="handleApplyFilters" @reset-filters="store.resetFilters" />

        <!-- Edit Material Dialog -->
        <EditMaterialDialog 
          v-model="showEditDialog" 
          @update="updateMaterial" 
          :disable="!imageUrlValid || validatingImage"
          :materialData="materialToEdit"
          :categories="categories"
          :suppliers="suppliers"
          :statuses="statuses"
          :defaultImageUrl="defaultImageUrl"
        />

        <DeleteDialog v-model="showDeleteDialog" itemType="material" :itemName="materialToDelete?.name || ''"
          @confirm-delete="confirmDelete" />
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
