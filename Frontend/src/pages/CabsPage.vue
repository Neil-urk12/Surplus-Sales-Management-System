<script setup lang="ts">
import { ref, onMounted, defineAsyncComponent } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
const ProductCardModal = defineAsyncComponent(() => import('src/components/Global/ProductModal.vue'));
const AddCabDialog = defineAsyncComponent(() => import('src/components/Cabs/AddCabDialog.vue'));
const EditCabDialog = defineAsyncComponent(() => import('src/components/Cabs/EditCabDialog.vue'));
const FilterDialog = defineAsyncComponent(() => import('src/components/Cabs/FilterDialog.vue'));
const DeleteDialog = defineAsyncComponent(() => import('src/components/Cabs/DeleteDialog.vue'));
const SellCabDialog = defineAsyncComponent(() => import('src/components/Cabs/SellCabDialog.vue'));
import { useQuasar } from 'quasar';
import { useCabsStore } from 'src/stores/cabs';
import { useAccessoriesStore } from 'src/stores/accessories';
import { useCustomerStore } from 'src/stores/customerStore';
import type { CabsRow, NewCabInput, CabStatus } from 'src/types/cabs';
import { getDefaultImage } from 'src/config/defaultImages';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';

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
  if (selected.value.image) {
    if (selected.value.image.startsWith('data:image/')) {
      const validationResult = validateAndSanitizeBase64Image(selected.value.image);
      if (validationResult.isValid && validationResult.sanitizedData) {
        selected.value.image = validationResult.sanitizedData;
      }
    }
  } else {
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

async function handleAddCab(cabData: NewCabInput) {
  try {
    const result = await store.addCab(cabData);
    if (result.success) {
      showAddDialog.value = false;
      operationNotifications.add.success(`cab: ${cabData.name}`);
    }
  } catch (error) {
    console.error('Error adding cab via handler:', error);
    operationNotifications.add.error('cab');
  }
}

function handleApplyFilters(filters: { make: string | null; color: string | null; status: CabStatus | null }) {
  store.filterMake = filters.make;
  store.filterColor = filters.color;
  store.filterStatus = filters.status;
  operationNotifications.filters.success();
  showFilterDialog.value = false;
}

function handleResetFilters() {
  store.resetFilters();
}

function editCab(cab: CabsRow) {
  cabToEdit.value = { ...cab };
  showEditDialog.value = true;
}

async function handleUpdateCab(updatedData: NewCabInput) {
  try {
    if (!cabToEdit.value || !cabToEdit.value.id) {
      throw new Error('No cab selected for update or missing ID');
    }
    const result = await store.updateCab(cabToEdit.value.id, updatedData);
    if (result.success) {
      showEditDialog.value = false;
      operationNotifications.update.success(`cab: ${updatedData.name}`);
      cabToEdit.value = null;
    }
  } catch (error) {
    console.error('Error updating cab via handler:', error);
    operationNotifications.update.error('cab');
  }
}

function deleteCab(cab: CabsRow) {
  cabToDelete.value = cab;
  showDeleteDialog.value = true;
}

async function handleConfirmDelete() {
  try {
    if (!cabToDelete.value) return;
    const result = await store.deleteCab(cabToDelete.value.id);
    if (result.success) {
      operationNotifications.delete.success('cab');
    } else {
      operationNotifications.delete.error('cab');
    }
    cabToDelete.value = null;
  } catch (error) {
    console.error('Error deleting cab:', error);
    operationNotifications.delete.error('cab');
    cabToDelete.value = null;
  }
}

function sellCab(cab: CabsRow) {
  cabToSell.value = cab;
  showSellDialog.value = true;
}

async function handleConfirmSell(payload: {
  customerId: string;
  quantity: number;
  accessories: Array<{ id: number; name: string; price: number; quantity: number; unitPrice: number }>
}) {
  try {
    if (!cabToSell.value) throw new Error('Cab to sell is not defined');

    const soldQuantity = payload.quantity;
    const cabName = cabToSell.value.name;

    if (soldQuantity <= 0 || soldQuantity > cabToSell.value.quantity) {
      operationNotifications.validation.error('Invalid quantity or not enough stock');
      return;
    }

    const newQuantity = cabToSell.value.quantity - soldQuantity;
    let newStatus: CabStatus = 'Available';
    if (newQuantity === 0) newStatus = 'Out of Stock';
    else if (newQuantity <= 2) newStatus = 'Low Stock';
    else if (newQuantity <= 5) newStatus = 'In Stock';

    const updatedCab: NewCabInput = {
      name: cabToSell.value.name,
      make: cabToSell.value.make,
      quantity: newQuantity,
      price: cabToSell.value.price,
      unit_color: cabToSell.value.unit_color,
      status: newStatus,
      image: cabToSell.value.image
    };

    const purchaseResult = await customerStore.recordCabPurchase(payload.customerId, {
      cabId: cabToSell.value.id,
      cabName: cabToSell.value.name,
      quantity: soldQuantity,
      unitPrice: cabToSell.value.price,
      accessories: payload.accessories
    });
    if (!purchaseResult.success) {
      operationNotifications.update.error('Failed purchase record');
      return;
    }

    const result = await store.updateCab(cabToSell.value.id, updatedCab);
    if (!result.success) {
      operationNotifications.update.error('Failed to update cab inventory');
      return;
    }

    try {
      await Promise.all(payload.accessories.map(async (acc) => {
        const accessory = accessoriesStore.accessoryRows.find(a => a.id === acc.id);
        if (accessory) {
          await accessoriesStore.updateAccessory(acc.id, {
            ...accessory,
            quantity: accessory.quantity - acc.quantity
          });
        }
      }));
    } catch (accError) {
      console.error('Failed to update accessory stock:', accError);
      $q.notify({ type: 'warning', message: `Failed to update stock for one or more accessories.` });
    }

    showSellDialog.value = false;
    operationNotifications.update.success(`Sold ${soldQuantity} ${cabName}`);
    cabToSell.value = null;

  } catch (error) {
    console.error('Error processing sell confirmation:', error);
    operationNotifications.update.error('cab sale processing');
  }
}

onMounted(async () => {
  try {
    await Promise.all([
      store.initializeCabs(),
      accessoriesStore.initializeAccessories(),
      customerStore.initializeCustomers()
    ]);
  } catch (error) {
    console.error('Error initializing data:', error);
  }
});

</script>

<template>
  <q-page class="flex q-pa-md">
    <div class="q-pa-sm full-width">
      <div class="flex row q-my-sm">
        <div class="flex full-width col">
          <div class="flex col q-mr-sm">
            <q-input v-model="store.rawCabSearch" outlined dense placeholder="Search" class="full-width">
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

      <q-table class="my-sticky-column-table" flat bordered title="Cabs" :rows="store.filteredCabRows"
        :columns="columns" row-key="id" :filter="store.cabSearch" @row-click="onRowClick"
        :pagination="{ rowsPerPage: 5 }" :loading="store.isLoading">
        <template v-slot:loading>
          <q-inner-loading showing color="primary">
            <q-spinner-gears size="50px" color="primary" />
          </q-inner-loading>
        </template>
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

      <ProductCardModal v-model="showProductCardModal" :image="selected?.image || ''" :title="selected?.name || ''"
        :unit_color="selected?.unit_color || ''" :price="selected?.price || 0" :quantity="selected?.quantity || 0"
        :details="`${selected?.make}`" :status="selected?.status || ''" @add="addToCart" />

      <AddCabDialog v-model="showAddDialog" :makes="makes" :colors="colors" :default-image-url="defaultImageUrl"
        @add-cab="handleAddCab" />

      <FilterDialog v-model="showFilterDialog" :makes="makes" :colors="colors" :statuses="statuses"
        :initial-filter-make="store.filterMake" :initial-filter-color="store.filterColor"
        :initial-filter-status="store.filterStatus" @apply-filters="handleApplyFilters"
        @reset-filters="handleResetFilters" />

      <EditCabDialog v-model="showEditDialog" :cab-data="cabToEdit" :makes="makes" :colors="colors"
        :default-image-url="defaultImageUrl" @update-cab="handleUpdateCab" />

      <DeleteDialog v-model="showDeleteDialog" item-type="cab" :item-name="cabToDelete?.name || 'this cab'"
        @confirm-delete="handleConfirmDelete" />

      <SellCabDialog v-model="showSellDialog" :cab-to-sell="cabToSell" :accessories="accessoriesStore.accessoryRows"
        :customer-store="customerStore" @confirm-sell="handleConfirmSell" />
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
