<script setup lang="ts">
import { watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useCustomerStore } from '../stores/customerStore'; 

interface Props {
  modelValue: boolean;
  customerId: string | null;
}

const props = defineProps<Props>();
const emit = defineEmits(['update:modelValue']);

const customerStore = useCustomerStore();
const { selectedCustomerHistory, isLoadingHistory, historyError } = storeToRefs(customerStore);
const { fetchPurchaseHistory } = customerStore;

// Fetch history when the modal becomes visible and customerId is valid
watch(() => [props.modelValue, props.customerId], ([newVisible, newCustomerId]) => {
  // Ensure customerId is a valid string before fetching
  if (newVisible && newCustomerId && typeof newCustomerId === 'string') {
    void fetchPurchaseHistory(newCustomerId);
  }
}, { immediate: true }); // Check immediately on component mount if modal starts open

const closeModal = () => {
  emit('update:modelValue', false);
};

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('en-PH', { 
    style: 'currency', 
    currency: 'PHP',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(value);
};

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-PH', { 
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

</script>

<template>
  <q-dialog :model-value="modelValue" @update:model-value="closeModal" persistent>
    <q-card style="min-width: 60vw; max-width: 80vw;">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">Purchase History</div>
        <q-space />
        <q-btn icon="close" flat round dense @click="closeModal" />
      </q-card-section>

      <q-separator />

      <q-card-section style="max-height: 70vh" class="scroll">
        <div v-if="isLoadingHistory" class="text-center q-pa-md">
          <q-spinner color="primary" size="3em" />
          <div class="q-mt-sm">Loading history...</div>
        </div>
        <q-banner v-else-if="historyError" inline-actions class="text-white bg-red">
          {{ historyError }}
           <template v-slot:action>
            <q-btn v-if="customerId" flat color="white" label="Retry" @click="fetchPurchaseHistory(customerId)" :loading="isLoadingHistory" />
          </template>
        </q-banner>
        <div v-else-if="!selectedCustomerHistory || selectedCustomerHistory.length === 0">
          <q-item>
            <q-item-section class="text-center text-grey">
              No purchase history found for this customer.
            </q-item-section>
          </q-item>
        </div>
        <div v-else>
          <q-list separator class="purchase-history-list">
            <q-item v-for="sale in selectedCustomerHistory" :key="sale.id" class="q-py-md">
              <q-item-section>
                <q-item-label class="text-subtitle1 text-weight-medium">
                  {{ formatDate(sale.saleDate) }}
                </q-item-label>
                
                <q-item-label class="q-mt-sm">
                  <div class="row items-center">
                    <!-- Main cab item -->
                    <div class="col-12 q-mb-sm">
                      <div v-for="item in sale.items.filter(i => i.itemType === 'Cab')" :key="item.id">
                        <span class="text-weight-medium">{{ item.quantity }}x {{ item.name || item.itemType }}</span>
                        <span class="q-mx-xs text-grey">•</span>
                        <span>{{ formatCurrency(item.unitPrice) }} each</span>
                      </div>
                    </div>
                    
                    <!-- Accessories with center dots -->
                    <div class="col-12 text-grey-8" v-if="sale.items.some(i => i.itemType === 'Accessory')">
                      <div class="text-caption q-mb-xs">Additional Accessories:</div>
                      <div 
                        v-for="(item) in sale.items.filter(i => i.itemType === 'Accessory')" 
                        :key="item.id"
                        class="q-ml-md text-caption accessory-item"
                      >
                        <span class="text-grey q-mr-sm">•</span>
                        <span>{{ item.quantity }}x {{ item.name || 'Unknown Accessory' }}</span>
                      </div>
                    </div>
                  </div>
                </q-item-label>
              </q-item-section>
              
              <q-item-section side>
                <div class="text-subtitle2 text-weight-medium text-primary">
                  {{ formatCurrency(sale.totalPrice) }}
                </div>
              </q-item-section>
            </q-item>
          </q-list>
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
.purchase-history-list {
  max-height: calc(70vh - 120px);
  overflow-y: auto;
}

.accessory-item {
  line-height: 1.5;
  display: flex;
  align-items: center;
}

/* Custom scrollbar styles */
.purchase-history-list::-webkit-scrollbar {
  width: 8px;
}

.purchase-history-list::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.purchase-history-list::-webkit-scrollbar-thumb {
  background: #ccc;
  border-radius: 4px;
}

.purchase-history-list::-webkit-scrollbar-thumb:hover {
  background: #999;
}
</style>
