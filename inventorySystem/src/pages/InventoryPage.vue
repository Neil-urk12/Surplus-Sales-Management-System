<template>
  <q-page class="flex q-pa-md">
    <div class="q-pa-sm full-width">
      <div class="flex justify-between items-center">
        <div class="flex gap-4">
          <q-input
            v-model="search"
            outlined
            dense
            placeholder="Search"
            class="w-64"
          >
            <template v-slot:prepend>
              <q-icon name="search" />
            </template>
          </q-input>
          <q-btn
            outline
            color="primary"
            icon="filter_list"
            label="Filters"
          />
        </div>
      </div>
      <q-table
        class="my-sticky-column-table"
        flat
        bordered
        title="Treats"
        :rows="rows"
        :columns="columns"
        row-key="name"
        :filter="search"
      />
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { QTableColumn } from 'quasar';

const search = ref('');

interface DessertRow {
  name: string;
  calories: number;
  fat: number;
  carbs: number;
  protein: number;
  sodium: number;
  calcium: string;
  iron: string;
}

const columns: QTableColumn[] = [
  {
    name: 'name',
    required: true,
    label: 'Dessert (100g serving)',
    align: 'left',
    field: 'name',
    sortable: true
  },
  { name: 'calories', align: 'center', label: 'Calories', field: 'calories', sortable: true },
  { name: 'fat', label: 'Fat (g)', field: 'fat', sortable: true },
  { name: 'carbs', label: 'Carbs (g)', field: 'carbs' },
  { name: 'protein', label: 'Protein (g)', field: 'protein' },
  { name: 'sodium', label: 'Sodium (mg)', field: 'sodium' },
  {
    name: 'calcium',
    label: 'Calcium (%)',
    field: 'calcium',
    sortable: true,
    sort: (a: string, b: string) => parseInt(a) - parseInt(b)
  },
  {
    name: 'iron',
    label: 'Iron (%)',
    field: 'iron',
    sortable: true,
    sort: (a: string, b: string) => parseInt(a) - parseInt(b)
  }
];

const rows: DessertRow[] = [ ];
</script>
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

