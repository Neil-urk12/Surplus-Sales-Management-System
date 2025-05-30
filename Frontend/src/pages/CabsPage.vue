<script setup lang="ts">
import { ref, onMounted, defineAsyncComponent } from 'vue';
import type { QTableColumn } from 'quasar';
const ProductCardModal = defineAsyncComponent(() => import('src/components/Global/ProductModal.vue'));
const AddCabDialog = defineAsyncComponent(() => import('src/components/Cabs/AddCabDialog.vue'));
const EditCabDialog = defineAsyncComponent(() => import('src/components/Cabs/EditCabDialog.vue'));
const FilterDialog = defineAsyncComponent(() => import('src/components/Cabs/FilterDialog.vue'));
const DeleteDialog = defineAsyncComponent(() => import('src/components/Cabs/DeleteDialog.vue'));
const SellCabDialog = defineAsyncComponent(() => import('src/components/Cabs/SellCabDialog.vue'));
const AdvancedSearch = defineAsyncComponent(() => import('src/components/Global/AdvancedSearch.vue'));
const CabsTable = defineAsyncComponent(() => import('src/components/Cabs/CabsTable.vue'));
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

function handleRowClick(evt: Event, row: CabsRow) {
  selected.value = { ...row };
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

// Function removed as it's now handled by the store's sellCab function

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

  // Record the purchase in customer records
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

  // Update accessory stock
  await updateAccessoryStock(accessories);

  // Process the cab sale using the new sellCab function
  const saleResult = await store.sellCab(cab.id, {
    customerId,
    quantity: soldQuantity,
    accessories: accessories.map(acc => ({
      id: acc.id,
      name: acc.name,
      price: acc.price,
      quantity: acc.quantity,
      unitPrice: acc.unitPrice
    }))
  });

  if (!saleResult.success) {
    throw new AppError(
      'Failed to process cab sale',
      'inventory',
      saleResult.error || 'Failed to update cab inventory. The purchase was recorded but inventory not updated.',
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
    const { type } = errorHandler.handleOperation(error, 'update', 'cab sale');

    if (type === 'inventory') {
      await errorHandler.recoverFromInventoryError([
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
      <div class="q-mb-md">
        <div class="flex row items-center justify-between">
          <div class="col">
            <div class="text-h5">Cabs</div>
          </div>
        </div>
        <div>
          <div class="text-caption text-grey q-mt-sm">Manage your inventory items, track stock levels, and monitor product details.</div>
          <div
            class="flex items-center q-mt-sm"
            :class="$q.screen.lt.md ? 'column q-gutter-y-sm items-stretch' : 'row justify-between'"
          >
            <div
              class="flex items-center"
              :class="$q.screen.lt.md ? 'column full-width q-gutter-y-sm items-stretch' : 'row q-gutter-x-sm'"
            >
              <AdvancedSearch
                v-model="store.search.searchInput"
                placeholder="Search cabs"
                @clear="handleResetFilters"
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
                @click="handleDownloadCsv"
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

      <CabsTable
        :rows="store.filteredCabRows || []"
        :columns="columns"
        :is-loading="store.isLoading"
        :search-value="store.search.searchValue"
        @row-click="handleRowClick"
        @edit="editCab"
        @delete="deleteCab"
        @sell="sellCab"
      />

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
.cab-page-z-top
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

// Custom styles for table text
.custom-table-text
  td,
  th
    font-size: 1.05em
    font-weight: 500

    .q-badge
      font-size: 1em
      font-weight: 600

</style>
