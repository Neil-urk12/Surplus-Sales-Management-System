<script setup lang="ts">
import { ref, watch, computed, onMounted, defineAsyncComponent } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import { useQuasar } from 'quasar';
import { useMaterialsStore } from 'src/stores/materials';
import type { MaterialRow, NewMaterialInput } from 'src/stores/materials';
import type { MaterialCategory, MaterialSupplier, MaterialStatus, UpdateMaterialInput } from 'src/types/materials';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';
const ProductCardModal = defineAsyncComponent(() => import('src/components/Global/ProductModal.vue'));
const DeleteDialog = defineAsyncComponent(() => import('src/components/Global/DeleteDialog.vue'));
const AddMaterialDialog = defineAsyncComponent(() => import('src/components/Materials/AddMaterialDialog.vue'));
const MaterialEditWrapper = defineAsyncComponent(() => import('src/components/Materials/MaterialEditWrapper.vue'));
const FilterMaterialDialog = defineAsyncComponent(() => import('src/components/Materials/FilterMaterialDialog.vue'));
const AdvancedSearch = defineAsyncComponent(() => import('src/components/Global/AdvancedSearch.vue'));
const ImageUploader = defineAsyncComponent(() => import('src/components/Global/ImageUploader.vue'));

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
  showAddDialog.value = true
}

async function handleAddMaterial(materialData: NewMaterialInput) {
  try {
    // Validate material data
    if (!materialData.name || materialData.name.trim() === '') {
      operationNotifications.validation.error('Material name is required');
      return;
    }
    
    // Validate category and supplier
    if (!materialData.category) {
      operationNotifications.validation.error('Material category is required');
      return;
    }
    
    if (!materialData.supplier) {
      operationNotifications.validation.error('Material supplier is required');
      return;
    }
    
    // Ensure image is provided
    if (!materialData.image) {
      materialData.image = defaultImageUrl;
    }
    
    const result = await store.addMaterial(materialData);
    if (result.success) {
      showAddDialog.value = false;
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

async function handleUpdateMaterial(materialData: UpdateMaterialInput) {
  try {
    console.log('handleUpdateMaterial called in MaterialsPage');
    
    // Validate material selection
    if (!materialToEdit.value || !materialToEdit.value.id) {
      throw new Error('No material selected for update or missing ID');
    }
    
    // Validate material data
    if (!materialData.name || materialData.name.trim() === '') {
      operationNotifications.validation.error('Material name is required');
      return;
    }
    
    // Validate category and supplier
    if (!materialData.category) {
      operationNotifications.validation.error('Material category is required');
      return;
    }
    
    if (!materialData.supplier) {
      operationNotifications.validation.error('Material supplier is required');
      return;
    }
    
    // Ensure image is provided
    if (!materialData.image) {
      materialData.image = defaultImageUrl;
    }

    const result = await store.updateMaterial(materialToEdit.value.id, materialData);
    console.log('Update result:', result);
    if (result.success) {
      closeEditDialog(); // Use the dedicated function
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
          <ImageUploader 
            :model-value="showEditDialog ? materialToEdit.image : newMaterial.image"
            @update:model-value="val => showEditDialog ? materialToEdit.image = val : newMaterial.image = val"
            :default-image-url="defaultImageUrl"
          />
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
          :itemName="materialToDelete?.name || 'Unknown material'" 
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
