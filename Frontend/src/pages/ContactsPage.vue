<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useQuasar } from 'quasar';
import type { Customer } from '../types/models';
import { useCustomerStore } from '../stores/customerStore';
import CustomerPurchaseHistoryModal from '../components/CustomerPurchaseHistoryModal.vue';
import ConfirmationDialog from '../components/ConfirmationDialog.vue';

const customerStore = useCustomerStore();
const { customers, isLoading, error } = storeToRefs(customerStore);
const { fetchCustomers, addCustomer, updateCustomer, deleteCustomer } = customerStore;

const $q = useQuasar();

const isCustomerModalOpen = ref(false);
const editingCustomer = ref<Customer | null>(null);
const newCustomer = ref<Omit<Customer, 'id' | 'createdAt' | 'updatedAt' | 'dateRegistered'>>({
  fullName: '',
  email: '',
  phone: '',
  address: '',
});

const isHistoryModalOpen = ref(false);
const selectedCustomerIdForHistory = ref<string | null>(null);

const isConfirmDialogOpen = ref(false);
const selectedCustomerId = ref<string | null>(null);

onMounted(async () => {
  await fetchCustomers();
});

const handleAddCustomer = async () => {
  if (newCustomer.value.fullName.trim() && newCustomer.value.email.trim() && newCustomer.value.phone.trim()) {
    await addCustomer(newCustomer.value);
    if (!error.value) {
      resetCustomerForm();
      isCustomerModalOpen.value = false;
      $q.notify({ type: 'positive', message: 'Customer added successfully!' });
    } else {
      $q.notify({ type: 'negative', message: `Error adding customer: ${error.value}` });
    }
  }
};

const handleEditCustomer = async () => {
  if (editingCustomer.value && newCustomer.value.fullName.trim() && newCustomer.value.email.trim() && newCustomer.value.phone.trim()) {
    await updateCustomer(editingCustomer.value.id, newCustomer.value);
    if (!error.value) {
      resetCustomerForm();
      isCustomerModalOpen.value = false;
      $q.notify({ type: 'positive', message: 'Customer updated successfully!' });
    } else {
      $q.notify({ type: 'negative', message: `Error updating customer: ${error.value}` });
    }
  }
};

const executeDelete = () => {
  if (selectedCustomerId.value) {
    deleteCustomer(selectedCustomerId.value);
    if (!error.value) {
      $q.notify({ type: 'positive', message: 'Customer deleted successfully!' });
    } else {
      $q.notify({ type: 'negative', message: `Error deleting customer: ${error.value}` });
    }
    selectedCustomerId.value = null;
  }
};

const handleDeleteCustomer = (customerId: string) => {
  selectedCustomerId.value = customerId;
  isConfirmDialogOpen.value = true;
};

const openAddCustomerModal = () => {
  editingCustomer.value = null;
  resetCustomerForm();
  isCustomerModalOpen.value = true;
};

const openEditCustomerModal = (customer: Customer) => {
  editingCustomer.value = customer;
  const { ...editableFields } = customer;
  newCustomer.value = { ...editableFields };
  isCustomerModalOpen.value = true;
};

const closeCustomerModal = () => {
  isCustomerModalOpen.value = false;
  editingCustomer.value = null;
  resetCustomerForm();
};

const resetCustomerForm = () => {
  newCustomer.value = {
    fullName: '',
    email: '',
    phone: '',
    address: '',
  };
};

const viewPurchaseHistory = (customerId: string) => {
  console.log('View history for customer:', customerId);
  selectedCustomerIdForHistory.value = customerId;
  isHistoryModalOpen.value = true;
};

</script>

<template>
  <q-page class="q-pa-md q-mt-md">
    <div class="q-pa-sm full-width">
      <div class="row justify-between items-center q-mb-md">
        <div class="text-h5">Customer Management</div>
        <div class="row q-gutter-x-sm items-center">
          <q-btn class="text-white bg-primary" @click="openAddCustomerModal" :disable="isLoading">
            <q-icon color="white" name="add" />
            Add Customer
          </q-btn>
        </div>
      </div>

      <div v-if="isLoading && !customers.length" class="text-center q-pa-md">
        <q-spinner-dots color="primary" size="40px" />
        <div>Loading customers...</div>
      </div>

      <q-banner v-if="error" inline-actions class="text-white bg-red q-mb-md">
        {{ error }}
        <template v-slot:action>
          <q-btn flat color="white" label="Retry" @click="fetchCustomers" :loading="isLoading" />
        </template>
      </q-banner>

      <!-- Empty state message -->
      <div v-if="!isLoading && customers.length === 0" class="text-center q-pa-xl text-grey">
        <q-icon name="person_search" size="4rem" />
        <div class="text-h6 q-mt-md">No customers found</div>
        <div class="q-mt-sm">Click 'Add Customer' to get started.</div>
      </div>

      <!-- Grid of customer cards -->
      <div v-else-if="!isLoading && customers.length > 0" class="row q-col-gutter-md">
        <div v-for="customer in customers" :key="customer.id" class="col-12 col-sm-6 col-md-4 col-lg-3">
          <q-card class="customer-card" flat bordered>
            <q-card-section>
              <div class="row items-center no-wrap">
                <div class="col">
                  <div class="text-h6">{{ customer.fullName }}</div>
                </div>
                <div class="col-auto">
                  <q-btn flat round dense icon="more_vert">
                    <q-menu>
                      <q-list style="min-width: 100px">
                        <q-item clickable v-close-popup @click="openEditCustomerModal(customer)">
                          <q-item-section avatar>
                            <q-icon name="edit" color="primary" />
                          </q-item-section>
                          <q-item-section>Edit</q-item-section>
                        </q-item>
                        <q-item clickable v-close-popup @click="handleDeleteCustomer(customer.id)">
                          <q-item-section avatar>
                            <q-icon name="delete" color="negative" />
                          </q-item-section>
                          <q-item-section>Delete</q-item-section>
                        </q-item>
                        <q-item clickable v-close-popup @click="viewPurchaseHistory(customer.id)">
                          <q-item-section avatar>
                            <q-icon name="receipt_long" color="info" />
                          </q-item-section>
                          <q-item-section>History</q-item-section>
                        </q-item>
                      </q-list>
                    </q-menu>
                  </q-btn>
                </div>
              </div>
            </q-card-section>

            <q-separator />

            <q-card-section>
              <div class="q-gutter-y-sm">
                <div class="row items-center">
                  <q-icon name="email" size="xs" class="q-mr-sm text-grey" />
                  <a :href="`mailto:${customer.email}`" class="text-primary">{{ customer.email }}</a>
                </div>
                <div class="row items-center">
                  <q-icon name="phone" size="xs" class="q-mr-sm text-grey" />
                  <a :href="`tel:${customer.phone}`" class="text-primary">{{ customer.phone }}</a>
                </div>
                <div class="row items-center">
                  <q-icon name="home" size="xs" class="q-mr-sm text-grey" />
                  <span>{{ customer.address || 'No address provided' }}</span>
                </div>
                <div class="row items-center text-caption text-grey">
                  <q-icon name="event" size="xs" class="q-mr-sm" />
                  <span>Registered: {{ new Date(customer.dateRegistered).toLocaleDateString() }}</span>
                </div>
              </div>
            </q-card-section>

            <q-card-actions align="right">
              <q-btn flat dense color="info" icon="receipt_long" @click="viewPurchaseHistory(customer.id)">
                <q-tooltip>View Purchase History</q-tooltip>
              </q-btn>
              <q-btn flat dense color="primary" icon="edit" @click="openEditCustomerModal(customer)">
                <q-tooltip>Edit Customer</q-tooltip>
              </q-btn>
              <q-btn flat dense color="negative" icon="delete" @click="handleDeleteCustomer(customer.id)">
                <q-tooltip>Delete Customer</q-tooltip>
              </q-btn>
            </q-card-actions>
          </q-card>
        </div>
      </div>

      <q-inner-loading :showing="isLoading && customers.length > 0">
        <q-spinner-gears size="50px" color="primary" />
      </q-inner-loading>

      <!-- Modals -->
      <q-dialog v-model="isCustomerModalOpen" @hide="closeCustomerModal" persistent>
        <q-card style="min-width: 350px;">
          <q-card-section class="row items-center">
            <div class="text-h6">{{ editingCustomer ? 'Edit Customer' : 'Add New Customer' }}</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>

          <q-form @submit.prevent="editingCustomer ? handleEditCustomer() : handleAddCustomer()">
            <q-card-section class="q-pt-none">
              <q-input
                filled
                v-model="newCustomer.fullName"
                label="Name *"
                autofocus
                :rules="[val => !!val || 'Name is required']"
                :disable="isLoading"
              />
              <q-input
                filled
                v-model="newCustomer.email"
                label="Email *"
                type="email"
                :rules="[val => !!val || 'Email is required', val => /.+@.+\..+/.test(val) || 'Invalid email format']"
                :disable="isLoading"
              />
              <q-input
                filled
                v-model="newCustomer.phone"
                label="Phone *"
                lazy-rules
                :rules="[
                  val => val && val.length > 0 || 'Please type the phone number',
                  val => /^\+639\d{9}$/.test(val) || 'Phone number must be in the format +639xxxxxxxxx'
                ]"
                :disable="isLoading"
              />
              <q-input
                filled
                v-model="newCustomer.address"
                label="Address"
                autogrow
                :disable="isLoading"
              />
            </q-card-section>

            <q-card-actions align="right">
              <q-btn flat label="Cancel" color="primary" v-close-popup :disable="isLoading" />
              <q-btn
                flat
                :label="editingCustomer ? 'Save Changes' : 'Add Customer'"
                color="primary"
                type="submit"
                :loading="isLoading"
                :disable="!newCustomer.fullName || !newCustomer.email || !newCustomer.phone || isLoading" />
            </q-card-actions>
          </q-form>
        </q-card>
      </q-dialog>

      <customer-purchase-history-modal
        v-model="isHistoryModalOpen"
        :customer-id="selectedCustomerIdForHistory"
      />

      <confirmation-dialog
        v-model="isConfirmDialogOpen"
        title="Delete Customer"
        message="Are you sure you want to delete this customer? This action cannot be undone."
        @confirm="executeDelete"
      />

    </div> <!-- End wrapper div -->
  </q-page>
</template>

<style scoped>
.customer-card {
  transition: all 0.3s;
  min-height: 220px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.customer-card:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.12);
}

.customer-card .text-h6 {
  font-size: 1.1rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 90%;
}

.customer-card .q-card-section {
  padding: 12px 16px;
}

.customer-card .q-card-actions {
  padding: 8px 16px 12px 16px;
}

.customer-card .row.items-center {
  flex-wrap: nowrap;
}

.customer-card .row.items-center > * {
  min-width: 0;
}

.customer-card .q-gutter-y-sm > .row {
  margin-bottom: 0.3em;
}

.customer-card a {
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
  display: inline-block;
  max-width: 100%;
}

@media (max-width: 599px) {
  .customer-card {
    min-height: 200px;
    font-size: 0.97rem;
  }
  .q-col-xs-12 {
    width: 100% !important;
    max-width: 100% !important;
  }
}

@media (max-width: 959px) and (min-width: 600px) {
  .customer-card {
    min-height: 210px;
    font-size: 1rem;
  }
}

@media (max-width: 1279px) and (min-width: 960px) {
  .customer-card {
    min-height: 220px;
    font-size: 1.05rem;
  }
}
</style>
