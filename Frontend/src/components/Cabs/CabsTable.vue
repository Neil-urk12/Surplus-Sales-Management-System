<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
import type { QTableColumn, QTableProps } from 'quasar';
import type { CabsRow } from 'src/types/cabs';

defineProps({
  rows: {
    type: Array as () => CabsRow[],
    required: true
  },
  columns: {
    type: Array as () => QTableColumn[],
    required: true
  },
  isLoading: {
    type: Boolean,
    default: false
  },
  searchValue: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['row-click', 'edit', 'delete', 'sell']);

const onRowClick: QTableProps['onRowClick'] = (evt, row) => {
  const target = evt.target as HTMLElement;
  if (target.closest('.action-button') || target.closest('.action-menu')) {
    return;
  }
  emit('row-click', evt, row);
};

function sellCab(row: CabsRow) {
  emit('sell', row);
}

function editCab(row: CabsRow) {
  emit('edit', row);
}

function deleteCab(row: CabsRow) {
  emit('delete', row);
}

function filterCabs(rows: readonly CabsRow[], terms: string): CabsRow[] {
  if (!terms) {
    return rows as CabsRow[];
  }

  const lowerTerms = terms.toLowerCase();
  return rows.filter(row => {
    return (
      (row.name && row.name.toLowerCase().includes(lowerTerms)) ||
      (row.make && row.make.toLowerCase().includes(lowerTerms)) ||
      (row.status && row.status.toLowerCase().includes(lowerTerms)) ||
      (row.unit_color && row.unit_color.toLowerCase().includes(lowerTerms))
    );
  });
}
</script>

<template>
  <template v-if="isLoading">
    <q-inner-loading showing color="primary">
      <q-spinner-gears size="50px" color="primary" />
    </q-inner-loading>
  </template>

  <template v-else>
    <q-table 
      class="my-sticky-column-table custom-table-text" 
      flat 
      bordered 
      :rows="rows" 
      :columns="columns" 
      row-key="id" 
      :filter="searchValue" 
      @row-click="onRowClick"
      :filter-method="filterCabs" 
      :pagination="{ rowsPerPage: 10 }" 
      :rows-per-page-options="[10]"
    >
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
                  @click.stop="sellCab(props.row)" 
                  role="button"
                  :aria-label="'Sell ' + props.row.name" 
                  v-if="props.row.quantity > 0"
                >
                  <q-item-section>
                    <q-item-label>
                      <q-icon name="sell" size="xs" class="q-mr-sm" aria-hidden="true" />
                      Sell
                    </q-item-label>
                  </q-item-section>
                </q-item>
                <q-item 
                  clickable 
                  v-close-popup 
                  @click.stop="editCab(props.row)" 
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
                  @click.stop="deleteCab(props.row)" 
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
      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-badge 
            :color="props.row.status === 'In Stock' ? 'green' : (props.row.status === 'Out of Stock' || props.row.status === 'Low Stock' ? 'red' : 'grey')" 
            :label="props.row.status" 
          />
        </q-td>
      </template>
    </q-table>
  </template>
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

.action-button
  position: relative
  z-index: 1

.action-menu
  z-index: 1001 !important

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
