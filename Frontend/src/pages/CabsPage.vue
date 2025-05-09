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
import { useDashboardStore } from 'src/stores/dashboardStore';
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
const dashboardStore = useDashboardStore();
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

// Utility function for formatting currency
function formatCurrency(value: number): string {
  return `₱ ${value.toLocaleString('en-PH', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })}`;
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
    name: 'price', label: 'Price', field: 'price', sortable: true, 
    format: (val: number) => formatCurrency(val)
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
  let newStatus: CabStatus;

  if (newQuantity === 0) {
    newStatus = 'Out of Stock';
  } else if (newQuantity <= 7) {
    newStatus = 'Low Stock';
  } else {
    newStatus = 'In Stock';
  }

  return { newQuantity, newStatus };
}

async function updateAccessoryStock(accessories: Array<{ id: number; quantity: number }>) {
  const updatePromises = accessories.map(async (acc) => {
    // Validate that quantity is non-negative
    if (acc.quantity < 0) {
      const errorMessage = `Invalid negative quantity for accessory ID ${acc.id}`;
      throw new AppError(errorMessage, 'validation', errorMessage, 'warning');
    }
    
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
  // Validate that quantity is positive and not more than available stock
  if (soldQuantity <= 0) {
    throw new AppError(
      `Invalid quantity: ${soldQuantity}. Quantity must be positive`,
      'validation',
      'Invalid quantity. Must be positive',
      'warning'
    );
  }
  
  if (soldQuantity > cab.quantity) {
    throw new AppError(
      `Not enough stock: requested ${soldQuantity}, available ${cab.quantity}`,
      'validation',
      'Not enough stock available',
      'warning'
    );
  }

  // Validate that all accessory quantities are non-negative
  for (const acc of accessories) {
    if (acc.quantity < 0) {
      throw new AppError(
        `Invalid quantity for accessory ${acc.name}: ${acc.quantity}`,
        'validation',
        'Accessory quantities must be non-negative',
        'warning'
      );
    }
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
  
  // Calculate the total sale amount including the cab and accessories
  const cabTotal = cab.price * soldQuantity;
  const accessoriesTotal = accessories.reduce((total, acc) => total + (acc.price * acc.quantity), 0);
  
  // Get customer info for the activity record
  const customer = customerStore.customers.find(c => c.id === customerId);
  const customerName = customer ? customer.fullName : 'a customer';
  
  // Record the cab sale in the dashboard
  dashboardStore.recordSale({
    itemName: cab.name,
    amount: cabTotal,
    quantity: soldQuantity,
    type: 'cab',
    customerName
  });
  
  // If accessories were sold, record them separately
  if (accessoriesTotal > 0) {
    // Create a concatenated name for multiple accessories
    const accessoryNames = accessories.map(a => a.name).join(', ');
    
    // Calculate total quantity of accessories
    const totalAccessoryQuantity = accessories.reduce((total, acc) => total + acc.quantity, 0);
    
    // Record the accessories sale
    dashboardStore.recordSale({
      itemName: accessoryNames,
      amount: accessoriesTotal,
      quantity: totalAccessoryQuantity,
      type: 'accessory',
      customerName
    });
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
    // Validate quantity is positive
    if (payload.quantity <= 0) {
      throw new AppError(
        'Quantity must be greater than zero',
        'validation',
        'Invalid quantity. Please enter a positive number.',
        'warning'
      );
    }
    
    await processCabSale(
      cabToSell.value,
      payload.customerId,
      payload.quantity,
      payload.accessories
    );

    // Calculate total value for the notification
    const cabTotal = cabToSell.value.price * payload.quantity;
    const accessoriesTotal = payload.accessories.reduce((total, acc) => total + (acc.price * acc.quantity), 0);
    const totalAmount = cabTotal + accessoriesTotal;
    
    // Format the price for display using the reusable function
    const formattedTotal = formatCurrency(totalAmount);

    showSellDialog.value = false;
    operationNotifications.update.success(`Sold ${payload.quantity} ${cabToSell.value.name} for ${formattedTotal}`);
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

function filterCabs(rows: readonly CabsRow[], terms: string): CabsRow[] {
  return rows.filter(row => {
    const searchTerms = terms.toLowerCase();
    return (
      row.name.toLowerCase().includes(searchTerms) ||
      row.make.toLowerCase().includes(searchTerms) ||
      row.status.toLowerCase().includes(searchTerms) ||
      row.unit_color.toLowerCase().includes(searchTerms)
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

      <template v-if="store.isLoading">
        <q-inner-loading showing color="primary">
          <q-spinner-gears size="50px" color="primary" />
        </q-inner-loading>
      </template>

      <template v-else>
        <q-table class="my-sticky-column-table custom-table-text" flat bordered :rows="store.filteredCabRows || []"
          :columns="columns" row-key="id" :filter="store.search.searchValue" @row-click="onRowClick"
          :filter-method="filterCabs" :pagination="{ rowsPerPage: 10 }" :rows-per-page-options="[10]">
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
          <template v-slot:body-cell-status="props">
            <q-td :props="props">
              <q-badge :color="props.row.status === 'In Stock' ? 'green' : (props.row.status === 'Out of Stock' || props.row.status === 'Low Stock' ? 'red' : 'grey')" :label="props.row.status" />
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
