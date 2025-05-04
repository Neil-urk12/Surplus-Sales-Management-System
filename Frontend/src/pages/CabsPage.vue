<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import ProductCardModal from 'src/components/Global/ProductModal.vue'
import { useQuasar } from 'quasar';
import { useCabsStore } from 'src/stores/cabs';
import type { CabsRow } from 'src/stores/cabs';

const $q = useQuasar();
const store = useCabsStore();
const showFilterDialog = ref(false);
const showAddDialog = ref(false);
const showEditDialog = ref(false);
const showDeleteDialog = ref(false);
const cabToDelete = ref<CabsRow | null>(null);
const show = ref(false);

const selected = ref<CabsRow>({
  name: '',
  id: 0,
  make: '',
  quantity: 0,
  price: 0,
  unit_color: '',
  status: '',
  image: '',
})

const newCab = ref<Omit<CabsRow, 'id'>>({
  name: '',
  make: '',
  quantity: 0,
  price: 0,
  unit_color: '',
  status: 'Out of Stock',
  image: 'https://loremflickr.com/600/400/car',
})

// Image validation
const imageUrlValid = ref(true);
const validatingImage = ref(false);
const defaultImageUrl = 'https://loremflickr.com/600/400/car';

// Available options from store
const { makes, colors, statuses } = store;

const capitalizedName = computed({
  get: () => newCab.value.name,
  set: (value: string) => {
    if (value) {
      newCab.value.name = value.charAt(0).toUpperCase() + value.slice(1);
    } else {
      newCab.value.name = value;
    }
  }
});

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
    align: 'center'
  }
];

const onRowClick: QTableProps['onRowClick'] = (evt, row) => {
  // Check if the click originated from the action button or its menu
  const target = evt.target as HTMLElement;
  if (target.closest('.action-button') || target.closest('.action-menu')) {
    return; // Do nothing if clicked on action button or its menu
  }
  selected.value = row as CabsRow
  show.value = true
}

function addToCart() {
  console.log('added to cart for', selected.value.name)
  show.value = false
}

function openAddDialog() {
  newCab.value = {
    name: '',
    make: '',
    quantity: 0,
    price: 0,
    unit_color: '',
    status: 'Out of Stock',
    image: defaultImageUrl
  }
  imageUrlValid.value = true;
  showAddDialog.value = true
}

async function addNewCab() {
  try {
    // Validate image URL before proceeding
    if (!imageUrlValid.value) {
      $q.notify({
        color: 'negative',
        message: 'Please provide a valid image URL',
        position: 'top',
        timeout: 2000
      });
      return;
    }

    // If image URL is empty, use default
    if (!newCab.value.image) {
      newCab.value.image = defaultImageUrl;
    }

    // Execute the store action and await its completion
    const result = await store.addCab(newCab.value);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showAddDialog.value = false;

      $q.notify({
        color: 'positive',
        message: `Added new cab: ${newCab.value.name}`,
        position: 'top',
        timeout: 2000
      });
    }
  } catch (error) {
    console.error('Error adding cab:', error);
    $q.notify({
      color: 'negative',
      message: 'Failed to add cab',
      position: 'top',
      timeout: 2000
    });
  }
}

function applyFilters() {
  showFilterDialog.value = false;

  $q.notify({
    color: 'positive',
    message: 'Filters applied successfully',
    position: 'top',
    timeout: 2000
  });
}

// Add watch for quantity changes
watch(() => newCab.value.quantity, (newQuantity) => {
  if (newQuantity === 0) {
    newCab.value.status = 'Out of Stock';
  } else if (newQuantity <= 2) {
    newCab.value.status = 'Low Stock';
  } else if (newQuantity <= 5) {
    newCab.value.status = 'In Stock';
  } else {
    newCab.value.status = 'Available';
  }
});

// Function to validate if URL is a valid image
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

  try {
    const result = await new Promise<boolean>((resolve) => {
      const img = new Image();

      const cleanup = () => {
        img.onload = null;
        img.onerror = null;
      };

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
      setTimeout(() => {
        cleanup();
        imageUrlValid.value = false;
        validatingImage.value = false;
        resolve(false);
      }, 5000);

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

// Modify the watch for image URL changes to handle the default image case
watch(() => newCab.value.image, async (newUrl: string) => {
  if (!newUrl || newUrl === defaultImageUrl) {
    imageUrlValid.value = true; // Default image or empty should be valid
    return;
  }
  try {
    if (newUrl.startsWith('data:image/')) {
      imageUrlValid.value = true; // Base64 image data is valid
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
  if (!file.type.startsWith('image/')) {
    $q.notify({
      color: 'negative',
      message: 'Please upload an image file',
      position: 'top',
    });
    return;
  }

  const reader = new FileReader();
  reader.onload = (e) => {
    if (e.target?.result) {
      previewUrl.value = e.target.result as string;
      newCab.value.image = e.target.result as string;
      imageUrlValid.value = true;
    }
  };
  reader.readAsDataURL(file);
}

// Function to remove image
function removeImage(event: Event) {
  event.stopPropagation(); // Prevent triggering file input click
  previewUrl.value = '';
  newCab.value.image = defaultImageUrl;
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
  newCab.value.image = defaultImageUrl;
  if (fileInput.value) {
    fileInput.value.value = '';
  }
}

// Function to handle edit cab
function editCab(cab: CabsRow) {
  selected.value = { ...cab };
  newCab.value = {
    name: cab.name,
    make: cab.make,
    quantity: cab.quantity,
    price: cab.price,
    unit_color: cab.unit_color,
    status: cab.status,
    image: cab.image
  };
  previewUrl.value = cab.image; // Set the preview URL for the existing image
  showEditDialog.value = true;
}

// Function to handle update cab
async function updateCab() {
  try {
    // Validate image URL before proceeding
    if (!imageUrlValid.value) {
      $q.notify({
        color: 'negative',
        message: 'Please provide a valid image URL',
        position: 'top',
        timeout: 2000
      });
      return;
    }

    // If image URL is empty, use default
    if (!newCab.value.image) {
      newCab.value.image = defaultImageUrl;
    }

    // Execute the store action and await its completion
    const result = await store.updateCab(selected.value.id, newCab.value);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showEditDialog.value = false;
      clearImageInput();

      $q.notify({
        color: 'positive',
        message: `Updated cab: ${newCab.value.name}`,
        position: 'top',
        timeout: 2000
      });
    }
  } catch (error) {
    console.error('Error updating cab:', error);
    $q.notify({
      color: 'negative',
      message: 'Failed to update cab',
      position: 'top',
      timeout: 2000
    });
  }
}

// Function to handle delete cab
function deleteCab(cab: CabsRow) {
  cabToDelete.value = cab;
  showDeleteDialog.value = true;
}

// Function to confirm and execute delete
async function confirmDelete() {
  try {
    if (!cabToDelete.value) return;

    await store.deleteCab(cabToDelete.value.id);
    showDeleteDialog.value = false;
    cabToDelete.value = null;

    $q.notify({
      color: 'positive',
      message: `Successfully deleted cab`,
      position: 'top',
      timeout: 2000
    });
  } catch (error) {
    console.error('Error deleting cab:', error);
    $q.notify({
      color: 'negative',
      message: 'Failed to delete cab',
      position: 'top',
      timeout: 2000
    });
  }
}

</script>

<template>
  <q-page class="flex q-pa-md">
    <div class="q-pa-sm full-width">
      <div class="flex row q-my-sm">
        <div class="flex full-width col">
          <div class="flex col q-mr-sm">
            <q-input
              v-model="store.rawCabSearch"
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

      <!--INVENTORY TABLE-->
      <q-table
        class="my-sticky-column-table"
        flat
        bordered
        title="Cabs"
        :rows="store.filteredCabRows"
        :columns="columns"
        row-key="id"
        :filter="store.cabSearch"
        @row-click="onRowClick"
        :pagination="{ rowsPerPage: 5 }"
      >
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" auto-width>
            <q-btn flat round dense color="grey" icon="more_vert" class="action-button">
              <q-menu class="action-menu">
                <q-list style="min-width: 100px">
                  <q-item clickable v-close-popup @click.stop="editCab(props.row)">
                    <q-item-section>
                      <q-item-label>
                        <q-icon name="edit" size="xs" class="q-mr-sm" />
                        Edit
                      </q-item-label>
                    </q-item-section>
                  </q-item>
                  <q-item clickable v-close-popup @click.stop="deleteCab(props.row)">
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
        </template>
      </q-table>

      <!-- Existing Cab Modal -->
      <ProductCardModal
        v-model="show"
        :image="selected?.image || ''"
        :title="selected?.name || ''"
        :unit_color="selected?.unit_color || ''"
        :price="selected?.price || 0"
        :quantity="selected?.quantity || 0"
        :details="`${selected?.make}`"
        :status="selected?.status || ''"
        @add="addToCart"
      />

      <!-- Add Cab Dialog -->
      <q-dialog
        v-model="showAddDialog"
        persistent
        @hide="clearImageInput"
      >
        <q-card style="min-width: 400px; max-width: 95vw">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6">New Cab</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>

          <q-card-section>
            <q-form @submit.prevent="addNewCab" class="q-gutter-sm">
              <q-input
                v-model="capitalizedName"
                label="Cab Name"
                dense
                outlined
                required
                :rules="[val => !!val || 'Name is required']"
              >
                <template v-slot:prepend>
                  <q-icon name="directions_car" />
                </template>
              </q-input>

              <div class="row q-col-gutter-sm">
                <div class="col-12 col-sm-6">
                  <q-select
                    v-model="newCab.make"
                    :options="makes"
                    label="Make"
                    dense
                    outlined
                    required
                    :rules="[val => !!val || 'Make is required']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="business" />
                    </template>
                  </q-select>
                </div>

                <div class="col-12 col-sm-6">
                  <q-select
                    v-model="newCab.unit_color"
                    :options="colors"
                    label="Color"
                    dense
                    outlined
                    required
                    :rules="[val => !!val || 'Color is required']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="palette" />
                    </template>
                  </q-select>
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-12 col-sm-6">
                  <q-input
                    v-model.number="newCab.quantity"
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
                    v-model.number="newCab.price"
                    type="number"
                    label="Price"
                    dense
                    outlined
                    required
                    :rules="[val => val > 0 || 'Price must be greater than 0']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="attach_money" />
                    </template>
                  </q-input>
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-12">
                  <q-input
                    v-model="newCab.status"
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
                          @click.stop="removeImage($event)"
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
              label="Add Cab"
              @click="addNewCab"
              :disable="!newCab.name || !newCab.make || !newCab.unit_color || newCab.quantity < 0 || newCab.price <= 0"
            />
          </q-card-actions>
        </q-card>
      </q-dialog>

      <!-- Filter Dialog -->
      <q-dialog v-model="showFilterDialog">
        <q-card style="min-width: 350px">
          <q-card-section class="row items-center">
            <div class="text-h6">Filter Cabs</div>
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

      <!-- Edit Cab Dialog -->
      <q-dialog
        v-model="showEditDialog"
        persistent
        @hide="clearImageInput"
      >
        <q-card style="min-width: 400px; max-width: 95vw">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6">Edit Cab</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>

          <q-card-section>
            <q-form @submit.prevent="updateCab" class="q-gutter-sm">
              <q-input
                v-model="capitalizedName"
                label="Cab Name"
                dense
                outlined
                required
                :rules="[val => !!val || 'Name is required']"
              >
                <template v-slot:prepend>
                  <q-icon name="directions_car" />
                </template>
              </q-input>

              <div class="row q-col-gutter-sm">
                <div class="col-12 col-sm-6">
                  <q-select
                    v-model="newCab.make"
                    :options="makes"
                    label="Make"
                    dense
                    outlined
                    required
                    :rules="[val => !!val || 'Make is required']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="business" />
                    </template>
                  </q-select>
                </div>

                <div class="col-12 col-sm-6">
                  <q-select
                    v-model="newCab.unit_color"
                    :options="colors"
                    label="Color"
                    dense
                    outlined
                    required
                    :rules="[val => !!val || 'Color is required']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="palette" />
                    </template>
                  </q-select>
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-12 col-sm-6">
                  <q-input
                    v-model.number="newCab.quantity"
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
                    v-model.number="newCab.price"
                    type="number"
                    label="Price"
                    dense
                    outlined
                    required
                    :rules="[val => val > 0 || 'Price must be greater than 0']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="attach_money" />
                    </template>
                  </q-input>
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-12">
                  <q-input
                    v-model="newCab.status"
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
                    <div v-if="!previewUrl && !newCab.image" class="text-center">
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
                        <img :src="previewUrl || newCab.image" class="preview-image" />
                      </div>
                      <div class="col-4 text-center">
                        <q-btn
                          flat
                          round
                          color="negative"
                          icon="close"
                          @click.stop="removeImage($event)"
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
              label="Update Cab"
              @click="updateCab"
              :disable="!newCab.name || !newCab.make || !newCab.unit_color || newCab.quantity < 0 || newCab.price <= 0"
            />
          </q-card-actions>
        </q-card>
      </q-dialog>

      <!-- Delete Confirmation Dialog -->
      <q-dialog v-model="showDeleteDialog" persistent>
        <q-card>
          <q-card-section class="row items-center">
            <q-avatar icon="warning" color="negative" text-color="white" />
            <span class="q-ml-sm text-h6">Delete Cab</span>
          </q-card-section>

          <q-card-section>
            Are you sure you want to delete {{ cabToDelete?.name }}? This action cannot be undone.
          </q-card-section>

          <q-card-actions align="right">
            <q-btn flat label="Cancel" v-close-popup />
            <q-btn flat label="Delete" color="negative" @click="confirmDelete" />
          </q-card-actions>
        </q-card>
      </q-dialog>
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
