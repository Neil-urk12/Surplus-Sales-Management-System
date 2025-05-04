<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useCustomerStore } from '../stores/customerStore'; // Assuming history logic is added here
import { Sale, SaleItem } from '../types/models'; // Assuming Sale/SaleItem types exist

interface Props {
  modelValue: boolean; // Controls modal visibility (v-model)
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
    fetchPurchaseHistory(newCustomerId);
  }
}, { immediate: true }); // Check immediately on component mount if modal starts open

const closeModal = () => {
  emit('update:modelValue', false);
};

// Helper to format currency
const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(value);
  // Adjust currency code (e.g., 'PHP') as needed
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
          <q-list separator>
            <q-expansion-item
              v-for="sale in selectedCustomerHistory" 
              :key="sale.id"
              group="sales"
              icon="shopping_cart"
              :label="`Sale ID: ${sale.id} - Date: ${new Date(sale.saleDate).toLocaleDateString()}`"
              :caption="`Total: ${formatCurrency(sale.totalPrice)} - Sold By: ${sale.soldBy}`"
              header-class="text-primary"
            >
              <q-card>
                <q-card-section>
                  <div class="text-subtitle2 q-mb-sm">Items Sold:</div>
                  <q-list dense bordered padding class="rounded-borders">
                     <q-item v-if="!sale.items || sale.items.length === 0">
                      <q-item-section class="text-grey">No items recorded for this sale.</q-item-section>
                    </q-item>
                    <q-item v-for="item in sale.items" :key="item.id">
                      <q-item-section>
                        <q-item-label>Type: {{ item.itemType }}</q-item-label>
                         <q-item-label caption>Details: 
                           <span v-if="item.multiCabId">MultiCab ID: {{ item.multiCabId }}</span>
                           <span v-else-if="item.accessoryId">Accessory ID: {{ item.accessoryId }}</span>
                           <span v-else-if="item.materialId">Material ID: {{ item.materialId }}</span>
                           <span v-else>N/A</span>
                         </q-item-label>
                      </q-item-section>
                      <q-item-section side>
                         <q-item-label>Qty: {{ item.quantity }}</q-item-label>
                        <q-item-label caption>Price: {{ formatCurrency(item.unitPrice) }}</q-item-label>
                        <q-item-label caption>Subtotal: {{ formatCurrency(item.subtotal) }}</q-item-label>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-card-section>
              </q-card>
            </q-expansion-item>
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
/* Add any specific styles if needed */
</style>
