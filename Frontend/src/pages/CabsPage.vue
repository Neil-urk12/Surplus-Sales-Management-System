<script setup lang="ts">
import { ref, onMounted, defineAsyncComponent } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
const ProductCardModal = defineAsyncComponent(() => import('src/components/Global/ProductModal.vue'));
const AddCabDialog = defineAsyncComponent(() => import('src/components/Cabs/AddCabDialog.vue'));
const EditCabDialog = defineAsyncComponent(() => import('src/components/Cabs/EditCabDialog.vue'));
const FilterDialog = defineAsyncComponent(() => import('src/components/Cabs/FilterDialog.vue'));
const DeleteDialog = defineAsyncComponent(() => import('src/components/Cabs/DeleteDialog.vue'));
const SellCabDialog = defineAsyncComponent(() => import('src/components/Cabs/SellCabDialog.vue'));
const AdvancedSearch = defineAsyncComponent(() => import('src/components/Global/AdvancedSearch.vue'));
import { useQuasar } from 'quasar';
import { useCabsStore } from 'src/stores/cabs';
import { useAccessoriesStore } from 'src/stores/accessories';
import { useCustomerStore } from 'src/stores/customerStore';
import type { CabsRow, NewCabInput, CabStatus, CabMake, CabColor } from 'src/types/cabs';
import { getDefaultImage } from 'src/config/defaultImages';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';
import { exportToCsv } from '../utils/exportUtils';
import { AppError, errorHandler } from '../utils/errorHandling';

const $q = useQuasar();
const store = useCabsStore();
const accessoriesStore = useAccessoriesStore();
const customerStore = useCustomerStore();
const showFilterDialog = ref(false);
const showAddDialog = ref(false);
const showEditDialog = ref(false);
const showDeleteDialog = ref(false);
const showSellDialog = ref(false);
const cabToDelete = ref<CabsRow | null>(null);
const cabToSell = ref<CabsRow | null>(null);
const cabToEdit = ref<CabsRow | null>(null);
const showProductCardModal = ref(false);
const selected = ref<CabsRow | null>(null);

const defaultImageUrl = getDefaultImage('cab');

const { makes, colors, statuses } = store;

function getValidatedImage(image: string | null | undefined): string {
  if (image && image.startsWith('data:image/')) {
    const validationResult = validateAndSanitizeBase64Image(image);
    if (validationResult.isValid && validationResult.sanitizedData) {
      return validationResult.sanitizedData;
    }

    console.warn('Invalid Base64 image detected, using default.');
    return defaultImageUrl;
  }
  return image || defaultImageUrl;
}

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
  const target = evt.target as HTMLElement;
  if (target.closest('.action-button') || target.closest('.action-menu')) {
    return;
  }
  selected.value = { ...row as CabsRow };

  selected.value.image = getValidatedImage(selected.value.image);
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

async function handleAddCab(cabData: NewCabInput) {
  try {
    const result = await store.addCab(cabData);
    if (result.success) {
      showAddDialog.value = false;
      operationNotifications.add.success(`cab: ${cabData.name}`);
    } else {
      throw new AppError(
        `Failed to add cab: ${cabData.name}`,
        'database',
        'Failed to add cab to system'
      );
    }
  } catch (error) {
    errorHandler.handleOperation(error, 'add', 'cab');
  }
}

function handleApplyFilters(filters: { make: string | null; color: string | null; status: CabStatus | null }) {
  store.filterMake = filters.make === null ? '' : filters.make as CabMake;
  store.filterColor = filters.color === null ? '' : filters.color as CabColor;
  store.filterStatus = filters.status === null ? '' : filters.status;
  operationNotifications.filters.success();
  showFilterDialog.value = false;
}

async function handleResetFilters() {
  try {
    await store.resetFilters();
  } catch (error) {
    errorHandler.handle(error, 'resetting filters');
  }
}

function editCab(cab: CabsRow) {
  cabToEdit.value = { ...cab };
  showEditDialog.value = true;
}

async function handleUpdateCab(updatedData: NewCabInput) {
  try {
    if (!cabToEdit.value || !cabToEdit.value.id) {
      throw new AppError(
        'No cab selected for update or missing ID',
        'validation',
        'No cab selected for update'
      );
    }

    const result = await store.updateCab(cabToEdit.value.id, updatedData);
    if (result.success) {
      showEditDialog.value = false;
      operationNotifications.update.success(`cab: ${updatedData.name}`);
      cabToEdit.value = null;
    } else {
      throw new AppError(
        `Failed to update cab: ${updatedData.name}`,
        'database',
        'Failed to update cab in system'
      );
    }
  } catch (error) {
    errorHandler.handleOperation(error, 'update', 'cab');
  }
}

function deleteCab(cab: CabsRow) {
  cabToDelete.value = cab;
  showDeleteDialog.value = true;
}

async function handleConfirmDelete() {
  try {
    if (!cabToDelete.value) {
      throw new AppError(
        'No cab selected for deletion',
        'validation',
        'No cab selected for deletion'
      );
    }

    const result = await store.deleteCab(cabToDelete.value.id);
    if (result.success) {
      operationNotifications.delete.success('cab');
    } else {
      throw new AppError(
        `Failed to delete cab: ${cabToDelete.value.name}`,
        'database',
        'Failed to delete cab from system'
      );
    }
  } catch (error) {
    errorHandler.handleOperation(error, 'delete', 'cab');
  } finally {
    cabToDelete.value = null;
  }
}

function sellCab(cab: CabsRow) {
  cabToSell.value = cab;
  showSellDialog.value = true;
}

function calculateNewCabState(currentQuantity: number, soldQuantity: number): { newQuantity: number; newStatus: CabStatus } {
  const newQuantity = currentQuantity - soldQuantity;
  let newStatus: CabStatus = 'Available';

  if (newQuantity === 0) newStatus = 'Out of Stock';
  else if (newQuantity <= 2) newStatus = 'Low Stock';
  else if (newQuantity <= 5) newStatus = 'In Stock';

  return { newQuantity, newStatus };
}

async function updateAccessoryStock(accessories: Array<{ id: number; quantity: number }>) {
  const updatePromises = accessories.map(async (acc) => {
    const accessory = accessoriesStore.accessoryRows.find(a => a.id === acc.id);
    if (accessory && accessory.quantity >= acc.quantity) {
      await accessoriesStore.updateAccessory(acc.id, {
        ...accessory,
        quantity: accessory.quantity - acc.quantity
      });
    } else if (accessory) {
      const errorMessage = `Insufficient stock for accessory ID ${acc.id}`;
      throw new AppError(errorMessage, 'inventory', errorMessage, 'warning');
    } else {
      const errorMessage = `Accessory with ID ${acc.id} not found`;
      throw new AppError(errorMessage, 'inventory', errorMessage);
    }
  });

  await Promise.all(updatePromises);
}

async function processCabSale(
  cab: CabsRow,
  customerId: string,
  soldQuantity: number,
  accessories: Array<{ id: number; name: string; price: number; quantity: number; unitPrice: number }>
): Promise<void> {
  if (soldQuantity <= 0 || soldQuantity > cab.quantity) {
    throw new AppError(
      `Invalid quantity: ${soldQuantity} for cab with available quantity: ${cab.quantity}`,
      'validation',
      'Invalid quantity or not enough stock',
      'warning'
    );
  }

  const { newQuantity, newStatus } = calculateNewCabState(cab.quantity, soldQuantity);
  const updatedCabData: NewCabInput = {
    name: cab.name,
    make: cab.make,
    quantity: newQuantity,
    price: cab.price,
    unit_color: cab.unit_color,
    status: newStatus,
    image: cab.image
  };

  const purchaseResult = await customerStore.recordCabPurchase(customerId, {
    cabId: cab.id,
    cabName: cab.name,
    quantity: soldQuantity,
    unitPrice: cab.price,
    accessories: accessories
  });

  if (!purchaseResult.success) {
    throw new AppError(
      'Failed to record purchase in customer records',
      'purchase',
      'Failed to record purchase. Please try again.'
    );
  }

  await updateAccessoryStock(accessories);

  const result = await store.updateCab(cab.id, updatedCabData);
  if (!result.success) {
    throw new AppError(
      'Failed to update cab inventory after purchase recorded',
      'inventory',
      'Failed to update cab inventory. The purchase was recorded but inventory not updated.',
      'critical'
    );
  }
}

async function handleConfirmSell(payload: {
  customerId: string;
  quantity: number;
  accessories: Array<{ id: number; name: string; price: number; quantity: number; unitPrice: number }>
}) {
  if (!cabToSell.value) {
    const errorMessage = 'No cab selected for sale';
    $q.notify({ type: 'negative', message: errorMessage });
    operationNotifications.update.error('cab sale processing - internal error');
    return;
  }

  try {
    await processCabSale(
      cabToSell.value,
      payload.customerId,
      payload.quantity,
      payload.accessories
    );

    showSellDialog.value = false;
    operationNotifications.update.success(`Sold ${payload.quantity} ${cabToSell.value.name}`);
  } catch (error) {
    errorHandler.handleOperation(error, 'update', 'cab sale');
    
    if (error instanceof AppError && error.type === 'inventory') {
      await errorHandler.recoverFromInventoryError([
        async () => { await store.initializeCabs(); },
        async () => { await accessoriesStore.initializeAccessories(); }
        async () => { await store.initializeCabs(); },
        async () => { await accessoriesStore.initializeAccessories(); }
      ]);
    }
  } finally {
    cabToSell.value = null;
  }
}

function handleDownloadCsv() {
  const csvColumns = [
    { header: 'ID', field: 'id' as const },
    { header: 'Name', field: 'name' as const },
    { header: 'Make', field: 'make' as const },
    { header: 'Quantity', field: 'quantity' as const },
    { header: 'Price (PHP)', field: 'price' as const },
    { header: 'Status', field: 'status' as const },
    { header: 'Color', field: 'unit_color' as const }
  ];

  exportToCsv(store.filteredCabRows, 'cabs-inventory', csvColumns);
}

/* eslint-disable @typescript-eslint/no-unused-vars */
function filterCabs(
  rows: readonly Record<string, unknown>[], 
  terms: string, 
  cols: readonly QTableColumn[], 
  getCellValue: (col: QTableColumn, row: Record<string, unknown>) => unknown
): readonly Record<string, unknown>[] {
/* eslint-enable @typescript-eslint/no-unused-vars */
  return rows.filter(row => {
    const searchTerms = terms.toLowerCase();
    
    // Helper function to safely convert values to searchable strings
    const safeString = (value: unknown): string => {
      if (value === null || value === undefined) return '';
      if (typeof value === 'string') return value;
      if (typeof value === 'number') return value.toString();
      if (typeof value === 'boolean') return value.toString();
      return ''; // For objects, arrays, or other complex types, return empty string
    };
    
    return (
      safeString(row.name).toLowerCase().includes(searchTerms) ||
      safeString(row.make).toLowerCase().includes(searchTerms) ||
      safeString(row.status).toLowerCase().includes(searchTerms) ||
      safeString(row.unit_color).toLowerCase().includes(searchTerms)
    );
  });
}

onMounted(async () => {
  try {
    await Promise.all([
      store.initializeCabs(),
      accessoriesStore.initializeAccessories(),
      customerStore.fetchCustomers()
    ]);
  } catch (error) {
    errorHandler.handle(error, 'initializing application data');
  }
});

</script>

<template>
  <q-page class="flex q-pa-md">
    <div class="q-pa-sm full-width">
      <div class="flex row q-my-sm">
        <div class="flex full-width col">
          <div class="flex col q-mr-sm">
            <AdvancedSearch v-model="store.search.searchInput" placeholder="Search cabs" @clear="store.resetFilters"
              color="primary" />
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
            <q-btn dense flat class="bg-primary text-white q-pa-sm" @click="handleDownloadCsv">
              <q-icon name="download" color="white" />
              Download CSV
            </q-btn>
          </div>
        </div>
      </div>


      <template v-if="store.isLoading">
        <q-inner-loading showing color="primary">
          <q-spinner-gears size="50px" color="primary" />
        </q-inner-loading>
      </template>


      <template v-else>
        <q-table class="my-sticky-column-table" flat bordered title="Cabs" :rows="store.filteredCabRows || []"
          :columns="columns" row-key="id" :filter="store.search.searchValue" @row-click="onRowClick"
          :filter-method="filterCabs" :pagination="{ rowsPerPage: 5 }">
          <template v-slot:body-cell-actions="props">
            <q-td :props="props" auto-width :key="props.row.id">
              <q-btn flat round dense color="grey" icon="more_vert" class="action-button"
                :aria-label="'Actions for ' + props.row.name">
                <q-menu class="action-menu" :aria-label="'Available actions for ' + props.row.name">
                  <q-list style="min-width: 100px">
                    <q-item clickable v-close-popup @click.stop="sellCab(props.row)" role="button"
                      :aria-label="'Sell ' + props.row.name" v-if="props.row.quantity > 0">
                      <q-item-section>
                        <q-item-label>
                          <q-icon name="sell" size="xs" class="q-mr-sm" aria-hidden="true" />
                          Sell
                        </q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-item clickable v-close-popup @click.stop="editCab(props.row)" role="button"
                      :aria-label="'Edit ' + props.row.name">
                      <q-item-section>
                        <q-item-label>
                          <q-icon name="edit" size="xs" class="q-mr-sm" aria-hidden="true" />
                          Edit
                        </q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-item clickable v-close-popup @click.stop="deleteCab(props.row)" role="button"
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
      </template>

      <ProductCardModal v-model="showProductCardModal" :image="selected?.image || ''" :title="selected?.name || ''"
        :unit_color="selected?.unit_color || ''" :price="selected?.price || 0" :quantity="selected?.quantity || 0"
        :details="`${selected?.make}`" :status="selected?.status || ''" @add="addToCart" />

      <AddCabDialog v-model="showAddDialog" :makes="makes" :colors="colors" :default-image-url="defaultImageUrl"
        @add-cab="handleAddCab" />

      <FilterDialog v-model="showFilterDialog" :makes="makes" :colors="colors" :statuses="statuses"
        :initial-filter-make="store.filterMake === '' ? null : store.filterMake"
        :initial-filter-color="store.filterColor === '' ? null : store.filterColor"
        :initial-filter-status="store.filterStatus === '' ? null : store.filterStatus"
        @apply-filters="handleApplyFilters" @reset-filters="handleResetFilters" />


      <EditCabDialog v-if="cabToEdit" v-model="showEditDialog" :cab-data="cabToEdit" :makes="makes" :colors="colors"
        :default-image-url="defaultImageUrl" @update-cab="handleUpdateCab" />

      <DeleteDialog v-model="showDeleteDialog" item-type="cab" :item-name="cabToDelete?.name || 'this cab'"
        @confirm-delete="handleConfirmDelete" />


      <SellCabDialog v-if="cabToSell" v-model="showSellDialog" :cab-to-sell="cabToSell"
        :accessories="accessoriesStore.accessoryRows" :customer-store="customerStore"
        @confirm-sell="handleConfirmSell" />
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

.cab-page-z-top
  z-index: 1000

.action-button
  position: relative
  z-index: 1

.action-menu
  z-index: 1001 !important

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

</style>
