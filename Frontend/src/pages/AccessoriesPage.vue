<script setup lang="ts">
import { ref, onMounted, defineAsyncComponent, computed } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import { useQuasar } from 'quasar';
const ProductCardModal = defineAsyncComponent(() => import('src/components/Global/ProductModal.vue'));
const DeleteDialog = defineAsyncComponent(() => import('src/components/Global/DeleteDialog.vue'));
const AddAccessoryDialog = defineAsyncComponent(() => import('src/components/accessories/AddAccessoryDialog.vue'));
const EditAccessoryDialog = defineAsyncComponent(() => import('src/components/accessories/EditAccessoryDialog.vue'));
const AdvancedSearch = defineAsyncComponent(() => import('src/components/Global/AdvancedSearch.vue'));
const FilterDialog = defineAsyncComponent(() => import('src/components/Global/FilterDialog.vue'));
import { useAccessoriesStore } from 'src/stores/accessories';
import { useCustomerStore } from 'src/stores/customerStore';
import type { AccessoryRow, NewAccessoryInput, AccessoryMakeInput, AccessoryColorInput, AccessoryStatus } from 'src/types/accessories';
import { getDefaultImage } from 'src/config/defaultImages';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';

const store = useAccessoriesStore();
const customerStore = useCustomerStore();
const $q = useQuasar();
const showFilterDialog = ref(false);
const showAddDialog = ref(false);
const showEditDialog = ref(false);
const showDeleteDialog = ref(false);
const accessoryToDelete = ref<AccessoryRow | null>(null);
const showProductCardModal = ref(false);
const isAddLoading = ref(false);
const isEditLoading = ref(false);
const isDeleteLoading = ref(false);
const pageLoading = ref(true);

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

// Image validation
const defaultImageUrl = getDefaultImage('accessory');

// Available options from store
const { makes, colors, statuses } = store;

// Current filter values
const filterValues = computed(() => ({
  make: store.filterMake === '' ? null : store.filterMake,
  color: store.filterColor === '' ? null : store.filterColor,
  status: store.filterStatus === '' ? null : store.filterStatus
}));

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
  {
    name: 'price', label: 'Price', field: 'price', sortable: true, format: (val: number) =>
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
  showAddDialog.value = true;
}

// Add error handling for API connections
const apiErrorMessage = computed(() => store.apiError);
const hasApiError = computed(() => !!store.apiError);

// Function to retry loading data after an error
async function retryLoading() {
  await new Promise(resolve => setTimeout(resolve, 1500));
  await store.initializeAccessories();
}

// Function to show error notification
function showErrorNotification(message: string) {
  // Use the custom notification utility
  operationNotifications.validation.error(message || 'An error occurred');
}

// Enhanced error handling in CRUD operations
async function handleAddAccessory(newAccessoryData: NewAccessoryInput) {
  try {
    isAddLoading.value = true;
    const result = await store.addAccessory(newAccessoryData);

    if (result.success) {
      operationNotifications.add.success(`accessory: ${newAccessoryData.name}`);
      showAddDialog.value = false;
    } else if (result.error) {
      showErrorNotification(`Failed to add accessory: ${result.error}`);
    }
  } catch (error) {
    console.error('Error adding accessory:', error);
    operationNotifications.add.error('accessory');
  } finally {
    isAddLoading.value = false;
  }
}

// Function to handle edit accessory
function editAccessory(accessory: AccessoryRow) {
  selected.value = { ...accessory };
  showEditDialog.value = true;
}

// Handler for the edit dialog submit event
async function handleUpdateAccessory(id: number, updatedAccessory: NewAccessoryInput) {
  try {
    isEditLoading.value = true;
    const result = await store.updateAccessory(id, updatedAccessory);

    if (result.success) {
      operationNotifications.update.success(`accessory: ${updatedAccessory.name}`);
      showEditDialog.value = false;
    } else if (result.error) {
      showErrorNotification(`Failed to update accessory: ${result.error}`);
    }
  } catch (error) {
    console.error('Error updating accessory:', error);
    operationNotifications.update.error('accessory');
  } finally {
    isEditLoading.value = false;
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

    isDeleteLoading.value = true;
    const result = await store.deleteAccessory(accessoryToDelete.value.id);

    if (result.success) {
      operationNotifications.delete.success('accessory');
      accessoryToDelete.value = null;
      showDeleteDialog.value = false;
    } else if (result.error) {
      showErrorNotification(`Failed to delete accessory: ${result.error}`);
    }
  } catch (error) {
    console.error('Error deleting accessory:', error);
    operationNotifications.delete.error('accessory');
  } finally {
    if (showDeleteDialog.value) showDeleteDialog.value = false;
    isDeleteLoading.value = false;
  }
}

// Function to handle applying filters
function handleApplyFilters(filters: Record<string, string | null>) {
  try {
    // Convert null values to empty strings for the store
    store.filterMake = filters.make ? filters.make as AccessoryMakeInput : '';
    store.filterColor = filters.color ? filters.color as AccessoryColorInput : '';
    store.filterStatus = filters.status ? filters.status as AccessoryStatus : '';

    // Get count of active filters for notification message
    const activeFilterCount = Object.values(filters).filter(value => value !== null).length;

    // Show appropriate notification based on if filters were applied or cleared
    if (activeFilterCount > 0) {
      operationNotifications.filters.success();
    } else {
      // If all filters were cleared, consider resetting search as well
      operationNotifications.filters.success();
    }

    // Close the filter dialog
    showFilterDialog.value = false;
  } catch (error) {
    console.error('Error applying filters:', error);
    showErrorNotification('Failed to apply filters');
  }
}

onMounted(async () => {
  pageLoading.value = true;
  try {
    await Promise.all([
      store.initializeAccessories(),
      customerStore.fetchCustomers()
    ]);
  } finally {
    pageLoading.value = false;
  }
});
</script>

<template>
  <div class="q-pa-md page-height">
    <div class="q-pa-sm full-width">

      <div class="q-mb-md">
        <div class="flex row items-center justify-between">
          <div class="col">
            <div class="text-h5">Accessories</div>
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
                placeholder="Search accessories"
                @clear="store.resetFilters"
                color="primary"
                :disable="pageLoading"
                :style="$q.screen.lt.md ? { width: '100%' } : { width: '400px' }"
              />
              <q-btn
                outline
                icon="filter_list"
                label="Filters"
                @click="showFilterDialog = true"
                :disable="pageLoading"
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
                :disable="pageLoading"
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
                :disable="pageLoading"
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

      <!-- Skeleton Loader -->
      <div v-if="pageLoading" class="q-my-md">
        <div class="row items-center q-mb-md">
          <q-skeleton type="text" class="text-h6" width="180px" />
          <q-space />
          <q-skeleton type="QBtn" width="50px" />
        </div>

        <!-- Skeleton Table Header -->
        <div class="row q-py-sm bg-grey-2">
          <div class="col-1"><q-skeleton type="text" /></div>
          <div class="col-3"><q-skeleton type="text" /></div>
          <div class="col-2"><q-skeleton type="text" /></div>
          <div class="col-1"><q-skeleton type="text" /></div>
          <div class="col-2"><q-skeleton type="text" /></div>
          <div class="col-1"><q-skeleton type="text" /></div>
          <div class="col-1"><q-skeleton type="text" /></div>
          <div class="col-1"><q-skeleton type="text" /></div>
        </div>

        <!-- Skeleton Table Rows -->
        <div v-for="n in 5" :key="n" class="row q-py-md q-gutter-y-sm items-center"
          :class="n % 2 === 0 ? 'bg-grey-1' : ''">
          <div class="col-1"><q-skeleton type="text" width="30px" /></div>
          <div class="col-3"><q-skeleton type="text" width="90%" /></div>
          <div class="col-2"><q-skeleton type="text" width="70%" /></div>
          <div class="col-1"><q-skeleton type="text" width="40px" /></div>
          <div class="col-2"><q-skeleton type="text" width="80px" /></div>
          <div class="col-1"><q-skeleton type="text" width="80%" /></div>
          <div class="col-1"><q-skeleton type="text" width="60%" /></div>
          <div class="col-1"><q-skeleton type="QBtn" width="40px" /></div>
        </div>

        <!-- Skeleton Pagination -->
        <div class="row justify-end q-mt-md">
          <q-skeleton type="QBtn" width="40px" class="q-mr-sm" />
          <q-skeleton type="text" width="80px" class="q-mr-sm" />
          <q-skeleton type="QBtn" width="40px" />
        </div>
      </div>

      <!-- API Error Message -->
      <div v-if="hasApiError && !pageLoading" class="q-py-md text-center">
        <q-banner rounded class="bg-negative text-white">
          <template v-slot:avatar>
            <q-icon name="error" />
          </template>
          <div class="text-body1 q-mb-sm">Failed to load accessories</div>
          <div class="text-caption q-mb-md">{{ apiErrorMessage }}</div>
          <q-btn color="white" text-color="negative" label="Retry" @click="retryLoading" />
        </q-banner>
      </div>

      <!--ACCESSORIES TABLE - Only show when not loading and no errors -->
      <q-table v-if="!pageLoading && !hasApiError" class="my-sticky-column-table custom-table-text" flat bordered
        :rows="store.filteredAccessoryRows" :columns="columns" row-key="id" :filter="store.search.searchValue"
        @row-click="onRowClick" :pagination="{ rowsPerPage: 10 }" :rows-per-page-options="[10]" :loading="store.isLoading">
        <template v-slot:loading>
          <q-inner-loading showing color="primary">
            <q-spinner-gears size="50px" color="primary" />
            <div class="q-mt-sm text-primary">Loading data...</div>
          </q-inner-loading>
        </template>
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge :color="props.row.status === 'In Stock' ? 'green' : (props.row.status === 'Out of Stock' || props.row.status === 'Low Stock' ? 'red' : 'grey')" :label="props.row.status" />
          </q-td>
        </template>
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" auto-width :key="props.row.id">
            <q-btn flat round dense color="grey" icon="more_vert" class="action-button"
              :aria-label="'Actions for ' + props.row.name" :disable="isDeleteLoading || isEditLoading">
              <q-menu class="action-menu" :aria-label="'Available actions for ' + props.row.name">
                <q-list style="min-width: 100px">
                  <q-item clickable v-close-popup @click.stop="editAccessory(props.row)" role="button"
                    :aria-label="'Edit ' + props.row.name" :disable="isEditLoading">
                    <q-item-section>
                      <q-item-label>
                        <q-icon name="edit" size="xs" class="q-mr-sm" aria-hidden="true" />
                        Edit
                      </q-item-label>
                    </q-item-section>
                  </q-item>
                  <q-item clickable v-close-popup @click.stop="deleteAccessory(props.row)" role="button"
                    :aria-label="'Delete ' + props.row.name" class="text-negative" :disable="isDeleteLoading">
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
      <ProductCardModal v-model="showProductCardModal" :image="selected?.image || ''" :title="selected?.name || ''"
        :unit_color="selected?.unit_color || ''" :price="selected?.price || 0" :quantity="selected?.quantity || 0"
        :details="`${selected?.make}`" :status="selected?.status || ''" @add="addToCart" />

      <!-- Add Accessory Dialog -->
      <AddAccessoryDialog v-model="showAddDialog" :makes="makes" :colors="colors" @submit="handleAddAccessory"
        :loading="isAddLoading" />

      <!-- Edit Accessory Dialog -->
      <EditAccessoryDialog v-model="showEditDialog" :makes="makes" :colors="colors" :accessory="selected"
        @submit="handleUpdateAccessory" :loading="isEditLoading" />

      <!-- Filter Dialog with direct passing of filter data -->
      <FilterDialog v-model="showFilterDialog" title="Filter Accessories" :filter-data="{
        make: { label: 'Make', options: makes, value: filterValues.make },
        color: { label: 'Color', options: colors, value: filterValues.color },
        status: { label: 'Status', options: statuses, value: filterValues.status }
      }" @apply-filters="handleApplyFilters" @reset-filters="store.resetFilters" />

      <!-- Delete Confirmation Dialog -->
      <DeleteDialog v-model="showDeleteDialog" itemType="accessory" :itemName="accessoryToDelete?.name || ''"
        @confirm-delete="confirmDelete" :loading="isDeleteLoading" />
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
.page-height
  height: 100vh

.custom-table-text
  td,
  th
    font-size: 1.15em
    font-weight: 400

    .q-badge
      font-size: 0.9em
      font-weight: 600
</style>
