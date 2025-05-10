<script setup lang="ts">
import { ref, watch, onMounted, defineAsyncComponent } from 'vue';
import type { QTableColumn } from 'quasar';
import { useQuasar } from 'quasar';
import { useMaterialsStore } from '../stores/materials';
import type { MaterialRow, NewMaterialInput } from '../stores/materials';
import type { MaterialStatus, MaterialCategory, MaterialSupplier } from '../types/materials';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';

const ProductCardModal = defineAsyncComponent(() => import('../components/Global/ProductModal.vue'));
const DeleteDialog = defineAsyncComponent(() => import('../components/Global/DeleteDialog.vue'));
const AddMaterialDialog = defineAsyncComponent(() => import('../components/Materials/AddMaterialDialog.vue'));
const MaterialEditWrapper = defineAsyncComponent(() => import('../components/Materials/MaterialEditWrapper.vue'));
const FilterMaterialDialog = defineAsyncComponent(() => import('../components/Materials/FilterMaterialDialog.vue'));

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
const defaultImageUrl = 'https://loremflickr.com/600/400/material';

// Available options from store
const { categories, suppliers, statuses } = store;

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
watch([store.filterCategory, store.filterSupplier, store.filterStatus], async () => {
  try {
    await store.onRequest({ pagination: store.pagination });
  } catch (error) {
    console.error('Error updating filters:', error);
  }
});

// Watch for raw search input changes
watch(() => store.rawMaterialSearch, async () => {
  try {
    await store.onRequest({ pagination: store.pagination });
  } catch (error) {
    console.error('Error updating search:', error);
  }
});

function openAddDialog() {
  newMaterial.value = {
    name: '',
    category: '',
    supplier: '',
    quantity: 0,
    status: 'Out of Stock',
    image: defaultImageUrl
  }
  showAddDialog.value = true
}

async function handleAddMaterial(materialData: NewMaterialInput) {
  try {
    // Validate material data
    if (!materialData.name || materialData.name.trim() === '') {
      operationNotifications.validation.error('Material name is required');
      return;
    }

    // If image URL is empty, use default
    if (!materialData.image) {
      materialData.image = defaultImageUrl;
    }

    // Execute the store action and await its completion
    const result = await store.addMaterial(materialData);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showAddDialog.value = false;
      operationNotifications.add.success('material');
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

// Modify the watch for image URL changes to handle base64 validation
watch(() => newMaterial.value.image, (newUrl: string) => {
  if (!newUrl || newUrl === defaultImageUrl) {
    return;
  }
  try {
    if (newUrl.startsWith('data:image/')) {
      const validationResult = validateAndSanitizeBase64Image(newUrl);
      if (validationResult.isValid) {
        newMaterial.value.image = validationResult.sanitizedData!;
      } else {
        $q.notify({
          color: 'negative',
          message: validationResult.error || 'Invalid image data',
          position: 'top',
        });
      }
    } else if (newUrl.startsWith('http')) {
      // URL validation happens in the component now
    } else {
      $q.notify({
        color: 'negative',
        message: 'Invalid image URL format',
        position: 'top',
      });
    }
  } catch (error) {
    console.error('Error in image URL watcher:', error);
  }
});

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
    if (!materialToDelete.value || materialToDelete.value.id === 0) {
      console.warn('No material selected for deletion');
      return;
    }

    const result = await store.deleteMaterial(materialToDelete.value.id);
    if (result.success) {
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
    }
  } catch (error) {
    console.error('Error deleting material:', error);
    operationNotifications.delete.error('material');
  }
}

// Add ref for edit dialog
const showEditDialog = ref(false);

// Function to handle edit material
function editMaterial(material: MaterialRow) {
  // Deep copy selected material to the dedicated edit ref
  materialToEdit.value = JSON.parse(JSON.stringify(material));
  imageUrlValid.value = true; // Reset validation state for the dialog
  validatingImage.value = false; // Reset validation state
  showEditDialog.value = true;
}

// Function to handle updating material
async function handleUpdateMaterial(materialData: NewMaterialInput) {
  try {
    // Validate material data
    if (!materialData.name || materialData.name.trim() === '') {
      operationNotifications.validation.error('Material name is required');
      return;
    }

    const currentFile = fileInput.value?.files?.[0];
    if (currentFile) {
      if (currentFile.size > MAX_FILE_SIZE) {
        $q.notify({
          type: 'negative',
          message: `File size (${(currentFile.size / 1024 / 1024).toFixed(2)}MB) exceeds the 5MB limit`,
          position: 'top',
          timeout: 3000
        });
        return;
      }

      await handleFile(currentFile);
    }

    // Execute the store action
    const result = await store.updateMaterial(materialToEdit.value.id, materialData);
    if (result.success) {
      showEditDialog.value = false;
      operationNotifications.update.success('material');
    }
  } catch (error) {
    console.error('Error in handleUpdateMaterial:', error);
    operationNotifications.update.error('material');
  }
}

// Add missing image handling functions
async function handleFile(fileToHandle: File): Promise<void> {
  if (fileToHandle.size > MAX_FILE_SIZE) {
    $q.notify({
      type: 'negative',
      message: `File size (${(fileToHandle.size / 1024 / 1024).toFixed(2)}MB) exceeds the 5MB limit`,
      position: 'top',
      timeout: 3000
    });
    return;
  }

  try {
    const base64String = await new Promise<string>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const result = e.target?.result;
        if (result && typeof result === 'string') {
          resolve(result);
        } else {
          reject(new Error('Failed to read file: Invalid result'));
        }
      };
      reader.onerror = () => reject(new Error(reader.error?.message || 'Failed to read file'));
      reader.readAsDataURL(fileToHandle);
    });

    const validationResult = validateAndSanitizeBase64Image(base64String);
    if (validationResult.isValid) {
      newMaterial.value.image = validationResult.sanitizedData!;
      imageUrlValid.value = true;
    } else {
      imageUrlValid.value = false;
      $q.notify({
        type: 'negative',
        message: validationResult.error || 'Invalid image format',
        position: 'top',
        timeout: 3000
      });
    }
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

// Add missing refs and constants
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const fileInput = ref<HTMLInputElement | null>(null);
const imageUrlValid = ref(true);
const validatingImage = ref(false);

// Status colors are now directly applied in the template

// Update onMounted hook
onMounted(async () => {
  console.log('MaterialsPage mounted, initializing materials...');
  try {
    await store.initializeMaterials();
    console.log('Materials initialized:', {
      rowCount: store.materialRows.length,
      pagination: store.pagination
    });
  } catch (error) {
    console.error('Error initializing materials:', error);
    $q.notify({
      type: 'negative',
      message: 'Failed to initialize materials',
      position: 'top',
      timeout: 3000
    });
  }
});

function handleApplyFilters(filters: { category: string | null; supplier: string | null; status: MaterialStatus | null }) {
  store.filterCategory = (filters.category || 'All') as MaterialCategory | 'All';
  store.filterSupplier = (filters.supplier || 'All') as MaterialSupplier | 'All';
  store.filterStatus = filters.status || 'All';

  // Reset to first page when applying filters
  store.pagination.page = 1;

  // Trigger the request with updated filters
  void store.onRequest({
    pagination: {
      ...store.pagination,
      page: 1
    }
  });

  operationNotifications.filters.success();
  showFilterDialog.value = false;
}

function addMaterialToCart() {
  console.log('added material to cart', selectedMaterial.value.name)
  showMaterial.value = false
}
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
            <!-- Add + Download CSV Group -->
            <div
              class="flex items-center"
              :class="$q.screen.lt.md ? 'column full-width q-gutter-y-sm items-stretch' : 'row q-gutter-x-sm'"
            >
              <q-btn
                unelevated
                @click="openAddDialog"
                :disable="store.isLoading"
                :class="[
                  $q.dark.isActive ? 'text-black bg-white' : 'text-white bg-primary',
                  { 'full-width': $q.screen.lt.md }
                ]"
              >
                <q-icon name="add" :color="$q.dark.isActive ? 'black' : 'white'" />
                <span :class="$q.dark.isActive ? 'text-black' : 'text-white'">Add</span>
              </q-btn>
              <q-btn
                dense
                flat
                :disable="store.isLoading"
                :class="[
                  $q.dark.isActive ? 'bg-white text-black' : 'bg-primary text-white',
                  'q-pa-sm',
                  { 'full-width': $q.screen.lt.md }
                ]"
              >
                <q-icon name="download" :color="$q.dark.isActive ? 'black' : 'white'" />
                <span :class="$q.dark.isActive ? 'text-black' : 'text-white'">Download CSV</span>
              </q-btn>
            </div>
          </div>
        </div>
      </div>

      <!-- Materials Section -->
      <div class="q-mt-sm">
        <!--MATERIALS TABLE-->
        <q-table
          class="my-sticky-column-table custom-table-text"
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

          <template v-slot:body="props">
            <q-tr :props="props" @click="() => {
              selectedMaterial = props.row;
              showMaterial = true;
            }" class="cursor-pointer">
              <q-td key="id" :props="props">{{ props.row.id }}</q-td>
              <q-td key="materialName" :props="props">{{ props.row.name }}</q-td>
              <q-td key="category" :props="props">{{ props.row.category }}</q-td>
              <q-td key="supplier" :props="props">{{ props.row.supplier }}</q-td>
              <q-td key="quantity" :props="props">{{ props.row.quantity }}</q-td>
              <q-td key="status" :props="props">
                <q-badge :color="props.row.status === 'In Stock' ? 'green' : (props.row.status === 'Out of Stock' ? 'red' : 'orange')" :label="props.row.status" />
              </q-td>
              <q-td key="actions" :props="props" @click.stop>
                <q-btn flat round dense color="grey" icon="more_vert" class="action-button">
                  <q-menu class="action-menu">
                    <q-list style="min-width: 100px">
                      <q-item clickable v-close-popup @click="editMaterial(props.row)">
                        <q-item-section>
                          <q-item-label>
                            <q-icon name="edit" size="xs" class="q-mr-sm" />
                            Edit
                          </q-item-label>
                        </q-item-section>
                      </q-item>
                      <q-item clickable v-close-popup @click="deleteMaterial(props.row)" class="text-negative">
                        <q-item-section>
                          <q-item-label class="text-negative">
                            <q-icon name="delete" size="xs" class="q-mr-sm" />
                            Delete
                          </q-item-label>
                        </q-item-section>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-btn>
              </q-td>
            </q-tr>
          </template>
        </q-table>

        <!-- Existing Material Modal -->
        <ProductCardModal v-model="showMaterial" :image="selectedMaterial?.image || ''"
          :title="selectedMaterial?.name || ''" :quantity="selectedMaterial?.quantity || 0"
          :details="`Supplier: ${selectedMaterial?.supplier}`" :unit_color="selectedMaterial?.category || ''"
          :status="selectedMaterial?.status || ''" @addItem="addMaterialToCart" />

        <!-- Add Material Dialog -->
        <AddMaterialDialog
          v-model="showAddDialog"
          :categories="categories"
          :suppliers="suppliers"
          :default-image-url="defaultImageUrl"
          @add-material="handleAddMaterial"
        />

        <!-- Filter Dialog -->
        <FilterMaterialDialog
          v-model="showFilterDialog"
          :categories="categories"
          :suppliers="suppliers"
          :statuses="statuses"
          :initial-filter-category="store.filterCategory === 'All' ? null : store.filterCategory"
          :initial-filter-supplier="store.filterSupplier === 'All' ? null : store.filterSupplier"
          :initial-filter-status="store.filterStatus === 'All' ? null : store.filterStatus"
          @apply-filters="handleApplyFilters"
          @reset-filters="store.resetFilters"
        />

        <!-- Edit Material Dialog -->
        <MaterialEditWrapper
          :open="showEditDialog"
          @update:open="showEditDialog = $event"
          :material-data="materialToEdit"
          :categories="categories"
          :suppliers="suppliers"
          :default-image-url="defaultImageUrl"
          @update-material="handleUpdateMaterial"
        />

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
  color: rgba(255, 255, 255, 0.9) !important

  .body--light &
    color: rgba(0, 0, 0, 0.87) !important

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

:deep(.q-pagination)
  .q-btn
    color: #ffffff !important
    &.text-primary
      color: #ffffff !important
      background: var(--q-primary)
      font-weight: bold

  .body--light &
    .q-btn
      color: rgba(0, 0, 0, 0.87) !important
      &.text-primary
        color: var(--q-primary) !important
        background: transparent
        font-weight: bold
</style>
