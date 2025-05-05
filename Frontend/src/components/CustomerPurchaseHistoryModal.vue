<script setup lang="ts">
import { watch, computed, ref } from 'vue';
import { storeToRefs } from 'pinia';
import { useCustomerStore } from '../stores/customerStore';
import type { PurchaseHistoryItem } from 'src/types/customers'; // Adjust path if needed

interface Props {
  modelValue: boolean;
  customerId: string | null;
}

const props = defineProps<Props>();
const emit = defineEmits(['update:modelValue']);

const customerStore = useCustomerStore();
const { selectedCustomerHistory, isLoadingHistory, historyError } = storeToRefs(customerStore);
const { fetchPurchaseHistory } = customerStore;

// Date filter refs
const dateRange = ref({
  from: '',
  to: ''
});

// Computed property for filtered and virtual scroll data
const filteredHistory = computed(() => {
  if (!selectedCustomerHistory.value) return [];

  let filtered = [...selectedCustomerHistory.value];

  if (dateRange.value.from || dateRange.value.to) {
    filtered = filtered.filter(sale => {
      const saleDate = new Date(sale.saleDate);
      const fromDate = dateRange.value.from ? new Date(dateRange.value.from) : null;
      const toDate = dateRange.value.to ? new Date(dateRange.value.to) : null;

      if (fromDate && toDate) {
        return saleDate >= fromDate && saleDate <= toDate;
      } else if (fromDate) {
        return saleDate >= fromDate;
      } else if (toDate) {
        return saleDate <= toDate;
      }
      return true;
    });
  }

  // Sort by date, most recent first
  return filtered.sort((a, b) => new Date(b.saleDate).getTime() - new Date(a.saleDate).getTime());
});

// Virtual scroll settings
const virtualScrollProps = {
  virtualScrollSliceSize: 15,
  virtualScrollSliceRatioBefore: 1,
  virtualScrollSliceRatioAfter: 1,
  virtualScrollItemSize: 60,
};

// Reset filters
const resetFilters = () => {
  dateRange.value = {
    from: '',
    to: ''
  };
};

// Fetch history when the modal becomes visible and customerId is valid
watch(() => [props.modelValue, props.customerId], ([newVisible, newCustomerId]) => {
  if (newVisible && newCustomerId && typeof newCustomerId === 'string') {
    void fetchPurchaseHistory(newCustomerId);
    resetFilters(); // Reset filters when modal opens
  }
}, { immediate: true });

const closeModal = () => {
  emit('update:modelValue', false);
};

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(value);
};

</script>

<template>
  <q-dialog :model-value="modelValue" @update:model-value="closeModal" persistent>
    <q-card class="purchase-history-modal">
      <q-card-section class="row items-center q-pb-none q-mb-sm">
        <div class="text-h6">Purchase History</div>
        <q-space />
        <div class="row items-center q-gutter-sm">
          <q-input v-model="dateRange.from" type="date" label="From" dense outlined class="col-auto"
            style="width: 150px" />
          <q-input v-model="dateRange.to" type="date" label="To" dense outlined class="col-auto" style="width: 150px" />
          <q-btn icon="restart_alt" flat round dense @click="resetFilters" :disabled="!dateRange.from && !dateRange.to">
            <q-tooltip>Reset Filters</q-tooltip>
          </q-btn>
        </div>
        <q-btn icon="close" flat round dense @click="closeModal" class="q-ml-sm">
          <q-tooltip>Close</q-tooltip>
        </q-btn>
      </q-card-section>

      <q-separator />

      <q-card-section class="scroll-section">
        <div v-if="filteredHistory.length" class="text-caption q-mb-sm">
          Showing {{ filteredHistory.length }} {{ filteredHistory.length === 1 ? 'record' : 'records' }}
        </div>
        <div v-if="isLoadingHistory" class="text-center q-pa-md">
          <q-spinner color="primary" size="3em" />
          <div class="q-mt-sm">Loading history...</div>
        </div>
        <q-banner v-else-if="historyError" inline-actions class="text-white bg-red">
          {{ historyError }}
          <template v-slot:action>
            <q-btn v-if="customerId" flat color="white" label="Retry" @click="fetchPurchaseHistory(customerId)"
              :loading="isLoadingHistory" />
          </template>
        </q-banner>
        <div v-else-if="!filteredHistory.length" class="text-center text-grey q-pa-md">
          {{ selectedCustomerHistory?.length ? 'No records match the selected filters.' :
            'No purchase history found for this customer.' }}
        </div>
        <div v-else>
          <q-virtual-scroll :items="filteredHistory" v-slot="{ item: sale }"
            :virtual-scroll-item-size="virtualScrollProps.virtualScrollItemSize"
            :virtual-scroll-slice-size="virtualScrollProps.virtualScrollSliceSize"
            :virtual-scroll-slice-ratio-before="virtualScrollProps.virtualScrollSliceRatioBefore"
            :virtual-scroll-slice-ratio-after="virtualScrollProps.virtualScrollSliceRatioAfter"
            class="virtual-scroll-list">
            <q-expansion-item :key="sale.id" group="sales" icon="shopping_cart"
              :label="`Sale ID: ${sale.id} - Date: ${new Date(sale.saleDate).toLocaleDateString(undefined, { dateStyle: 'medium' })}`"
              :caption="`Total: ${formatCurrency(sale.totalPrice)} - Sold By: ${sale.soldBy}`"
              header-class="text-primary">
              <q-card>
                <q-card-section>
                  <div class="text-subtitle2 q-mb-sm">Items Sold:</div>
                  <q-list dense bordered padding class="rounded-borders">
                    <q-item v-if="!sale.items?.length">
                      <q-item-section class="text-grey">No items recorded for this sale.</q-item-section>
                    </q-item>
                    <template v-else>
                      <!-- Display Cab first -->
                      <q-item v-for="item in sale.items.filter((i: PurchaseHistoryItem) => i.itemType === 'Cab')"
                        :key="item.id">
                        <q-item-section>
                          <q-item-label class="text-weight-bold">{{ item.name }}</q-item-label>
                          <q-item-label caption>
                            Quantity: {{ item.quantity }} | Price: {{ formatCurrency(item.unitPrice) }}
                          </q-item-label>
                        </q-item-section>
                        <q-item-section side>
                          <q-item-label>{{ formatCurrency(item.subtotal) }}</q-item-label>
                        </q-item-section>
                      </q-item>

                      <!-- Display Accessories if any -->
                      <q-item v-if="sale.items.some((i: PurchaseHistoryItem) => i.itemType === 'Accessory')">
                        <q-item-section>
                          <q-item-label class="text-weight-medium">Additional Accessories:</q-item-label>
                          <q-list dense class="q-ml-md">
                            <q-item
                              v-for="acc in sale.items.filter((i: PurchaseHistoryItem) => i.itemType === 'Accessory')"
                              :key="acc.id">
                              <q-item-section>
                                <div class="row items-center">
                                  <q-icon name="fiber_manual_record" size="xs" class="q-mr-sm" />
                                  <span>{{ acc.name }}</span>
                                  <span class="q-ml-sm text-caption">
                                    (Qty: {{ acc.quantity }} × {{ formatCurrency(acc.unitPrice) }})
                                  </span>
                                </div>
                              </q-item-section>
                              <q-item-section side>
                                {{ formatCurrency(acc.subtotal) }}
                              </q-item-section>
                            </q-item>
                          </q-list>
                        </q-item-section>
                      </q-item>

                      <!-- Total Section -->
                      <q-separator class="q-my-sm" />
                      <q-item>
                        <q-item-section>
                          <q-item-label class="text-weight-bold">Total</q-item-label>
                        </q-item-section>
                        <q-item-section side>
                          <q-item-label class="text-weight-bold text-primary">
                            {{ formatCurrency(sale.totalPrice) }}
                          </q-item-label>
                        </q-item-section>
                      </q-item>
                    </template>
                  </q-list>
                </q-card-section>
              </q-card>
            </q-expansion-item>
          </q-virtual-scroll>
        </div>
      </q-card-section>

      <q-separator />

      <q-card-actions align="right">
        <q-btn label="Close" color="primary" flat @click="closeModal" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style scoped>
.purchase-history-modal {
  min-width: 60vw;
  max-width: 80vw;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.scroll-section {
  flex: 1;
  overflow-y: auto;
  max-height: calc(90vh - 120px);

  /* Hide scrollbar for Chrome, Safari and Opera */
  &::-webkit-scrollbar {
    display: none;
  }

  /* Hide scrollbar for IE, Edge and Firefox */
  -ms-overflow-style: none;
  /* IE and Edge */
  scrollbar-width: none;
  /* Firefox */
}

.q-item__section--side {
  text-align: right;
  min-width: 100px;
}

.q-list.dense .q-item {
  min-height: 32px;
}

.virtual-scroll-list {
  height: calc(90vh - 180px);

  /* Hide scrollbar for Chrome, Safari and Opera */
  &::-webkit-scrollbar {
    display: none;
  }

  /* Hide scrollbar for IE, Edge and Firefox */
  -ms-overflow-style: none;
  /* IE and Edge */
  scrollbar-width: none;
  /* Firefox */
}
</style>
