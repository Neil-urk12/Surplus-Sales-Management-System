<script setup lang="ts">
import { ref, watch } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import ProductCardModal from 'src/components/Global/ProductModal.vue'
import { useQuasar } from 'quasar';
import { useMaterialsStore } from 'src/stores/materials';
import type { MaterialRow } from 'src/stores/materials';

const $q = useQuasar();
const store = useMaterialsStore();
const showFilterDialog = ref(false);
const selectedMaterial = ref<MaterialRow>({
  name: '',
  id: 0,
  category: '',
  supplier: '',
  quantity: 0,
  status: '',
  image: ''
})

const newMaterial = ref<Omit<MaterialRow, 'id'>>({
  name: '',
  category: '',
  supplier: '',
  quantity: 0,
  status: '',
  image: ''
})

// Available options from store
const { categories, suppliers, statuses } = store;

const materialColumns: QTableColumn[] = [
  {
    name: 'materialName',
    required: true,
    label: 'Material Name',
    align: 'left',
    field: 'name',
    sortable: true
  },
  { name: 'id', align: 'center', label: 'ID', field: 'id', sortable: true },
  { name: 'category', label: 'Category', field: 'category' },
  { name: 'supplier', label: 'Supplier', field: 'supplier' },
  { name: 'quantity', label: 'Quantity', field: 'quantity', sortable: true },
  { name: 'status', label: 'Status', field: 'status' },
];

const showMaterial = ref(false)
const showAddDialog = ref(false)

const onMaterialRowClick: QTableProps['onRowClick'] = (_e, row) => {
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
    image: 'https://loremflickr.com/600/400/material'
  }
  showAddDialog.value = true
}

function addNewMaterial() {
  try {
    store.addMaterial(newMaterial.value);
    
    // Close the dialog first
    showAddDialog.value = false;
    
    // Show success notification after dialog is closed
    setTimeout(() => {
      $q.notify({
        color: 'positive',
        message: `Added new material: ${newMaterial.value.name}`,
        position: 'top',
        timeout: 2000
      });
    }, 300);
  } catch (error) {
    console.error('Error adding material:', error);
    $q.notify({
      color: 'negative',
      message: 'Failed to add material',
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
watch(() => newMaterial.value.quantity, (newQuantity) => {
  if (newQuantity === 0) {
    newMaterial.value.status = 'Out of Stock';
  } else if (newQuantity < 10) {
    newMaterial.value.status = 'Low Stock';
  } else {
    newMaterial.value.status = 'In Stock';
  }
});

// Add the preview method in the script section
function previewImage(url: string) {
  if (url) {
    window.open(url, '_blank');
  }
}
</script>

<template>
  <q-page class="flex q-pa-md">
    <div class="q-pa-sm full-width">
      <!-- Materials Section -->
      <div class="q-mt-xl">
        <div class="text-h6 q-mb-md">Materials</div>
        <div class="flex row q-my-sm">
          <div class="flex full-width col">
            <div class="flex col q-mr-sm">
              <q-input v-model="store.materialSearch" outlined dense placeholder="Search" class="full-width">
                <template v-slot:prepend>
                  <q-icon name="search" />
                </template>
              </q-input>
            </div>
            <div class="flex col">
              <q-btn outline color="primary" icon="filter_list" label="Filters" @click="showFilterDialog = true" />
            </div>
          </div>

          <div class="flex row q-gutter-x-sm">
            <q-btn class="text-white bg-primary" @click="openAddDialog">
              <q-icon color="white" name="add" />
              Add
            </q-btn>
            <div class="flex row">
              <q-btn dense flat class="bg-primary text-white q-pa-sm">
                <q-icon color="white" name="download" />
                Download CSV
              </q-btn>
            </div>
          </div>
        </div>

        <!--MATERIALS TABLE-->
        <q-table 
          class="my-sticky-column-table" 
          flat 
          bordered 
          title="Materials" 
          :rows="store.filteredMaterialRows"
          :columns="materialColumns" 
          row-key="id" 
          :filter="store.materialSearch"
          @row-click="onMaterialRowClick" 
        />
        
        <!-- Existing Material Modal -->
        <ProductCardModal v-model="showMaterial" :image="selectedMaterial?.image || ''"
          :title="selectedMaterial?.name || ''" :price="0" :quantity="selectedMaterial?.quantity || 0"
          :details="`Supplier: ${selectedMaterial?.supplier}`" @add="addMaterialToCart" />
        
        <!-- Add Material Dialog - Minimalistic Design -->
        <q-dialog v-model="showAddDialog" persistent>
          <q-card style="min-width: 400px; max-width: 95vw">
            <q-card-section class="row items-center q-pb-none">
              <div class="text-h6">New Material</div>
              <q-space />
              <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>

            <q-card-section>
              <q-form @submit.prevent="addNewMaterial" class="q-gutter-sm">
                <q-input 
                  v-model="newMaterial.name" 
                  label="Material Name" 
                  dense 
                  outlined 
                  required 
                  :rules="[val => !!val || 'Name is required']"
                >
                  <template v-slot:prepend>
                    <q-icon name="inventory_2" />
                  </template>
                </q-input>

                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-6">
                    <q-select 
                      v-model="newMaterial.category" 
                      :options="categories" 
                      label="Category" 
                      dense 
                      outlined 
                      required 
                      :rules="[val => !!val || 'Category is required']"
                    >
                      <template v-slot:prepend>
                        <q-icon name="category" />
                      </template>
                    </q-select>
                  </div>
                  
                  <div class="col-12 col-sm-6">
                    <q-select 
                      v-model="newMaterial.supplier" 
                      :options="suppliers" 
                      label="Supplier" 
                      dense 
                      outlined 
                      required 
                      :rules="[val => !!val || 'Supplier is required']"
                    >
                      <template v-slot:prepend>
                        <q-icon name="local_shipping" />
                      </template>
                    </q-select>
                  </div>
                </div>

                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-6">
                    <q-input 
                      v-model.number="newMaterial.quantity" 
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
                      v-model="newMaterial.status" 
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

                <q-input 
                  v-model="newMaterial.image" 
                  label="Image URL" 
                  dense 
                  outlined
                  hint="Enter the URL for the material image"
                >
                  <template v-slot:prepend>
                    <q-icon name="image" />
                  </template>
                  <template v-slot:append>
                    <q-icon 
                      name="preview" 
                      class="cursor-pointer"
                      @click="previewImage(newMaterial.image)"
                    >
                      <q-tooltip>Preview Image</q-tooltip>
                    </q-icon>
                  </template>
                </q-input>
              </q-form>
            </q-card-section>

            <q-card-actions align="right" class="bg-dark text-primary q-pa-md">
              <q-btn flat label="Cancel" color="negative" v-close-popup />
              <q-btn unelevated color="primary" label="Add Material" @click="addNewMaterial" />
            </q-card-actions>
          </q-card>
        </q-dialog>

        <!-- Filter Dialog -->
        <q-dialog v-model="showFilterDialog">
          <q-card style="min-width: 350px">
            <q-card-section class="row items-center">
              <div class="text-h6">Filter Materials</div>
              <q-space />
              <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>

            <q-card-section class="q-pt-none">
              <q-select
                v-model="store.filterCategory"
                :options="categories"
                label="Category"
                clearable
                outlined
                class="q-mb-md"
              />
              
              <q-select
                v-model="store.filterSupplier"
                :options="suppliers"
                label="Supplier"
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
      </div>
    </div>
  </q-page>
</template>

<style lang="sass">
.my-sticky-column-table
  max-width: 100%

  thead tr:first-child th:first-child
    background-color: #00b4ff

  td:first-child
    background-color: #00b4ff

  th:first-child,
  td:first-child
    position: sticky
    left: 0
    z-index: 1

.z-top
  z-index: 1000
</style>
