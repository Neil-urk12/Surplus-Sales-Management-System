<script setup lang="ts">
import { ref } from 'vue';
import type { QTableColumn, QTableProps} from 'quasar';
import ProductCardModal from 'src/components/Global/ProductModal.vue'
const search = ref('');
const show = ref(false)
const title  = ref('')
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

interface CabsRow {
  name: string
  id: number
  make: string
  quantity: number
  price: number
  status: string
  unit_color: string
  image: string
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
  { name: 'price', label: 'Price', field: 'price', sortable: true, format: (val: number) =>
        `₱ ${val.toLocaleString('en-PH', {
          minimumFractionDigits: 2,
          maximumFractionDigits: 2
        })}`
  },

  { name: 'status', label: 'Status', field: 'status' },
  {name: 'color', label: 'Color', field: 'unit_color'},
];

const rows: CabsRow[] = [
  {
    name: 'RX‑7',
    id: 1,
    make: 'Mazda',
    quantity: 4,
    price: 7_000_000,
    status: 'In Stock',
    unit_color: 'Black',
    image: 'https://loremflickr.com/600/400/mazda',
  },
  {
    name: '911 GT3',
    id: 2,
    make: 'Porsche',
    quantity: 2,
    price: 10_000_000,
    status: 'In Stock',
    unit_color: 'Black',
    image: 'https://loremflickr.com/600/400/porsche',
  },
  {
    name: '911 GT3',
    id: 3,
    make: 'Porsche',
    quantity: 2,
    price: 10_000_000,
    status: 'Available',
    unit_color: 'Black',
    image: 'https://loremflickr.com/600/400/porsche',
  },
  {
    name: 'Corolla',
    id: 4,
    make: 'Toyota',
    quantity: 2,
    price: 10_000_000,
    status: 'In Stock',
    unit_color: 'Black',
    image: 'https://loremflickr.com/600/400/porsche',
  },
  {
    name: 'Navara',
    id: 5,
    make: 'Nissan',
    quantity: 2,
    price: 10_000_000,
    status: 'In Stock',
    unit_color: 'Black',
    image: 'src/assets/images/Cars/navara.avif',
  },
  {
    name: 'Vios',
    id: 6,
    make: 'Toyota',
    quantity: 2,
    price: 10_000_000,
    status: 'In Stock',
    unit_color: 'Black',
    image: 'https://loremflickr.com/600/400/porsche',
  },
  {
    name: 'Ranger',
    id: 7,
    make: 'Ford',
    quantity: 2,
    price: 10_000_000,
    status: 'In Stock',
    unit_color: 'Black',
    image: 'https://loremflickr.com/600/400/porsche',
  },
]

const onRowClick: QTableProps['onRowClick'] = (_e, row) => {
  selected.value = row as CabsRow
  show.value = true
}
//for future function
function addToCart () {
  console.log('added to cart for', title.value)
  show.value = false
}

</script>

<template>
  <q-page class="flex q-pa-md">
    <div class="q-pa-sm full-width">
      <div class="flex row q-my-sm ">

        <div class="flex full-width col">
          <div class="flex col q-mr-sm">
            <q-input
              v-model="search"
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
            />
          </div>
        </div>

        <div class="flex row q-gutter-x-sm">
          <q-btn class="text-white bg-primary">
            <q-icon color="white" name="add" />
            Add
          </q-btn>
          <div class="flex row">
            <q-btn dense flat class="bg-primary text-white q-pa-sm"  >
              <q-icon color="white" name="download" />
              Download CVS
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
        :rows="rows"
        :columns="columns"
        row-key="name"
        :filter="search"
        @row-click="onRowClick"
      />
      <ProductCardModal
        v-model="show"
        :image="selected?.image || ''"
        :title="selected?.name  || ''"
        :unit_color="selected?.unit_color || ''"
        :price="selected?.price || 0"
        :quantity="selected?.quantity || 0"
        :details="`${selected?.make}`"
        :status="selected?.status || ''"
        @add="addToCart"
      />
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

.body--dark .my-sticky-column-table th:nth-child(2),
.body--dark .my-sticky-column-table td:nth-child(2)
  color: black
</style>
