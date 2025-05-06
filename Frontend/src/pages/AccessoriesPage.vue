<script setup lang="ts">
import { ref, onMounted, defineAsyncComponent } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
const ProductCardModal = defineAsyncComponent(() => import('src/components/Global/ProductModal.vue'));
const DeleteDialog = defineAsyncComponent(() => import('src/components/Global/DeleteDialog.vue'));
const AddAccessoryDialog = defineAsyncComponent(() => import('src/components/accessories/AddAccessoryDialog.vue'));
const EditAccessoryDialog = defineAsyncComponent(() => import('src/components/accessories/EditAccessoryDialog.vue'));
import { useAccessoriesStore } from 'src/stores/accessories';
import type { AccessoryRow, NewAccessoryInput } from 'src/types/accessories';
import { getDefaultImage } from 'src/config/defaultImages';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';

const store = useAccessoriesStore();
const showFilterDialog = ref(false);
const showAddDialog = ref(false);
const showEditDialog = ref(false);
const showDeleteDialog = ref(false);
const accessoryToDelete = ref<AccessoryRow | null>(null);
const showProductCardModal = ref(false);

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
  { name: 'price', label: 'Price', field: 'price', sortable: true, format: (val: number) =>
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

async function handleAddAccessory(newAccessoryData: NewAccessoryInput) {
  try {
    // Execute the store action and await its completion
    const result = await store.addAccessory(newAccessoryData);

    // Show notification after operation successfully completes
    if (result.success) {
      operationNotifications.add.success(`accessory: ${newAccessoryData.name}`);
    }
  } catch (error) {
    console.error('Error adding accessory:', error);
    operationNotifications.add.error('accessory');
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
    // Execute the store action and await its completion
    const result = await store.updateAccessory(id, updatedAccessory);

    // Show notification after operation successfully completes
    if (result.success) {
      operationNotifications.update.success(`accessory: ${updatedAccessory.name}`);
    }
  } catch (error) {
    console.error('Error updating accessory:', error);
    operationNotifications.update.error('accessory');
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

    await store.deleteAccessory(accessoryToDelete.value.id);
    showDeleteDialog.value = false;
    accessoryToDelete.value = null;
    operationNotifications.delete.success('accessory');
  } catch (error) {
    console.error('Error deleting accessory:', error);
    operationNotifications.delete.error('accessory');
  }
}

function applyFilters() {
  showFilterDialog.value = false;
}

onMounted(() => {
  void store.initializeAccessories();
});
</script>

<template>
  <div class="q-pa-md">
    <div class="q-pa-sm full-width">
      <div class="flex row q-my-sm">
        <div class="flex full-width col">
          <div class="flex col q-mr-sm">
            <q-input
              v-model="store.rawAccessorySearch"
              outlined
              dense
              placeholder="Search"
              class="full-width"
            >
              <template v-slot:prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>
          <div class="flex col">
            <q-btn
              outline
              icon="filter_list"
              label="Filters"
              @click="showFilterDialog = true"
            />
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

      <!--ACCESSORIES TABLE-->
      <q-table
        class="my-sticky-column-table"
        flat
        bordered
        title="Accessories"
        :rows="store.filteredAccessoryRows"
        :columns="columns"
        row-key="id"
        :filter="store.accessorySearch"
        @row-click="onRowClick"
        :pagination="{ rowsPerPage: 5 }"
        :loading="store.isLoading"
      >
        <template v-slot:loading>
          <q-inner-loading showing color="primary">
            <q-spinner-gears size="50px" color="primary" />
          </q-inner-loading>
        </template>
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" auto-width :key="props.row.id">
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
                    @click.stop="editAccessory(props.row)"
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
                    @click.stop="deleteAccessory(props.row)"
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

      <!-- Existing Accessory Modal -->
      <ProductCardModal
        v-model="showProductCardModal"
        :image="selected?.image || ''"
        :title="selected?.name || ''"
        :unit_color="selected?.unit_color || ''"
        :price="selected?.price || 0"
        :quantity="selected?.quantity || 0"
        :details="`${selected?.make}`"
        :status="selected?.status || ''"
        @add="addToCart"
      />

      <!-- Add Accessory Dialog -->
      <AddAccessoryDialog
        v-model="showAddDialog"
        :makes="makes"
        :colors="colors"
        @submit="handleAddAccessory"
      />

      <!-- Edit Accessory Dialog -->
      <EditAccessoryDialog
        v-model="showEditDialog"
        :makes="makes"
        :colors="colors"
        :accessory="selected"
        @submit="handleUpdateAccessory"
      />

      <!-- Filter Dialog -->
      <q-dialog v-model="showFilterDialog">
        <q-card style="min-width: 350px">
          <q-card-section class="row items-center">
            <div class="text-h6">Filter Accessories</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>

          <q-card-section class="q-pt-none">
            <q-select
              v-model="store.filterMake"
              :options="makes"
              label="Make"
              clearable
              outlined
              class="q-mb-md"
            />

            <q-select
              v-model="store.filterColor"
              :options="colors"
              label="Color"
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

      <!-- Delete Confirmation Dialog -->
      <DeleteDialog
        v-model="showDeleteDialog"
        itemType="accessory"
        :itemName="accessoryToDelete?.name || ''"
        @confirm-delete="confirmDelete"
      />
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
</style> 
