<script setup lang="ts">
import { ref } from 'vue';
import type { QTableColumn } from 'quasar';
import ProductCardModal from 'src/components/Global/ProductModal.vue'
const search = ref('');
const selected = ref<DessertRow>({
  name: '',
  id: 0,
  make: '',
  quantity: 0,
  price: 0,
  status: '',
  image: ''
})

interface DessertRow {
  name: string
  id: number
  make: string
  quantity: number
  price: number
  status: string
  image: string

}

const columns: QTableColumn[] = [
  {
    name: 'unitName',
    required: true,
    label: 'Unit Name',
    align: 'left',
    field: 'name',
    sortable: true
  },
  { name: 'id', align: 'center', label: 'ID', field: 'id', sortable: true },
  { name: 'make', label: 'Make', field: 'make' },
  { name: 'quantity', label: 'Quantity', field: 'quantity', sortable: true },
  { name: 'price', label: 'Price', field: 'price', sortable: true },
  { name: 'status', label: 'Status', field: 'status' },
  // { name: 'action', label: 'Action', field: 'action' },
  // {
  //   name: 'calcium',
  //   label: 'Calcium (%)',
  //   field: 'calcium',
  //   sortable: true,
  //   sort: (a: string, b: string) => parseInt(a) - parseInt(b)
  // },
  // {
  //   name: 'iron',
  //   label: 'Iron (%)',
  //   field: 'iron',
  //   sortable: true,
  //   sort: (a: string, b: string) => parseInt(a) - parseInt(b)
  // }
];

const rows: DessertRow[] = [
  {
    name: 'RX‑7',
    id: 1,
    make: 'Mazda',
    quantity: 4,
    price: 7_000_000,
    status: 'Available',
    image: 'https://loremflickr.com/600/400/mazda',      // ✅ working link
  },
  {
    name: '911 GT3',
    id: 2,
    make: 'Porsche',
    quantity: 2,
    price: 10_000_000,
    status: 'Available',
    image: 'https://loremflickr.com/600/400/porsche',    // ✅ working link
  },
]




import type { QTableProps } from 'quasar'
const show = ref(false)
const title  = ref('')

const onRowClick: QTableProps['onRowClick'] = (_e, row) => {
  selected.value = row as DessertRow
  show.value = true
}

function addToCart () {
  console.log('added to cart for', title.value)
  show.value = false           // close after adding (optional)
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
              color="primary"
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
        :price="selected?.price || 0"
        :quantity="selected?.quantity || 0"
        :details="`Make: ${selected?.make}`"
        @add="addToCart"
      />
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
</style>

