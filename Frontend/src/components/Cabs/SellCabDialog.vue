<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import type { PropType } from 'vue';
import { useQuasar } from 'quasar';
import type { CabsRow } from 'src/types/cabs';
import type { AccessoryRow } from 'src/types/accessories';
import type { useCustomerStore } from 'src/stores/customerStore';
import { operationNotifications } from '../../utils/notifications';

const props = defineProps({
    modelValue: {
        type: Boolean,
        required: true,
    },
    cabToSell: {
        type: Object as PropType<CabsRow | null>,
        required: true,
    },
    accessories: {
        type: Array as PropType<AccessoryRow[]>,
        required: true,
    },
    customerStore: {
        type: Object as PropType<ReturnType<typeof useCustomerStore>>,
        required: true,
    },
});

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'confirm-sell', payload: {
        customerId: string;
        quantity: number;
        accessories: Array<{ id: number; name: string; price: number; quantity: number; unitPrice: number }>
    }): void
}>();

const $q = useQuasar();

// Local state for the dialog
const customerId = ref<string>('');
const sellQuantity = ref(1);
const selectedAccessoryId = ref<number | null>(null);
const selectedAccessories = ref<Array<{
    id: number;
    name: string;
    price: number;
    quantity: number;
    availableQuantity: number;
}>>([]);
const accessoryQuantity = ref(1);

const totalAccessoriesPrice = computed(() => {
    return selectedAccessories.value.reduce((total, acc) => total + (acc.price * acc.quantity), 0);
});

const totalPrice = computed(() => {
    const cabPrice = (props.cabToSell?.price || 0) * sellQuantity.value;
    return cabPrice + totalAccessoriesPrice.value;
});

const isCustomerIdValid = computed(() => {
    return props.customerStore.validateCustomerId(customerId.value).isValid;
});

// Reset local state when dialog opens/cab changes
watch(() => props.modelValue, (newValue) => {
    if (newValue) {
        customerId.value = '';
        sellQuantity.value = 1;
        selectedAccessories.value = [];
        selectedAccessoryId.value = null;
        accessoryQuantity.value = 1;
    }
});

// Functions
function addAccessory() {
    if (!selectedAccessoryId.value || accessoryQuantity.value <= 0) return;

    const accessory = props.accessories.find(a => a.id === selectedAccessoryId.value);
    if (!accessory) return;

    const existing = selectedAccessories.value.find(a => a.id === accessory.id);
    if (existing) {
        operationNotifications.validation.error(`${accessory.name} is already added`);
        return;
    }

    if (accessoryQuantity.value > accessory.quantity) {
        operationNotifications.validation.error(`Not enough ${accessory.name} in stock`);
        return;
    }

    selectedAccessories.value.push({
        id: accessory.id,
        name: accessory.name,
        price: accessory.price,
        quantity: accessoryQuantity.value,
        availableQuantity: accessory.quantity
    });

    selectedAccessoryId.value = null;
    accessoryQuantity.value = 1;
}

function removeAccessory(id: number) {
    selectedAccessories.value = selectedAccessories.value.filter(acc => acc.id !== id);
}

function validateCustomerInput(id: string | number | null): void {
    const validation = props.customerStore.validateCustomerId(String(id || ''));
    if (validation.isValid) {
        const customer = validation.customer;
        $q.notify({
            type: 'positive',
            message: `Customer found: ${customer?.fullName}`,
            position: 'top',
            timeout: 2000
        });
    }
}

function handleSellClick() {
    // Local validation before emitting
    if (!customerId.value) return operationNotifications.validation.error('Customer ID required');
    if (!isCustomerIdValid.value) return operationNotifications.validation.error('Invalid Customer ID');
    if (sellQuantity.value <= 0) return operationNotifications.validation.error('Quantity must be > 0');
    if (sellQuantity.value > (props.cabToSell?.quantity || 0)) return operationNotifications.validation.error('Not enough stock');

    for (const acc of selectedAccessories.value) {
        const availableQuantity = props.accessories.find(a => a.id === acc.id)?.quantity || 0;
        if (acc.quantity > availableQuantity) {
            return operationNotifications.validation.error(`Not enough ${acc.name} in stock`);
        }
    }

    emit('confirm-sell', {
        customerId: customerId.value,
        quantity: sellQuantity.value,
        accessories: selectedAccessories.value.map(acc => ({
            id: acc.id,
            name: acc.name,
            price: acc.price,
            quantity: acc.quantity,
            unitPrice: acc.price // Add unitPrice for the payload expected by parent
        }))
    });
}

function closeDialog() {
    emit('update:modelValue', false);
}

</script>

<template>
    <q-dialog :model-value="modelValue" persistent @update:model-value="closeDialog">
        <q-card style="min-width: 400px">
            <q-card-section class="row items-center">
                <q-avatar icon="sell" color="primary" text-color="white" />
                <span class="q-ml-sm text-h6">Sell Cab</span>
            </q-card-section>

            <q-card-section>
                <div class="text-body1 q-mb-md">Selling {{ cabToSell?.name }}</div>

                <q-input v-model="customerId" label="Customer ID *" dense outlined class="q-mb-md"
                    :rules="[val => !!val || 'ID required', val => customerStore.validateCustomerId(String(val || '')).isValid || 'Invalid ID']"
                    @update:model-value="validateCustomerInput" lazy-rules>
                    <template v-slot:prepend><q-icon name="person" /></template>
                    <template v-slot:append>
                        <q-icon :name="isCustomerIdValid ? 'check_circle' : 'error'"
                            :color="isCustomerIdValid ? 'positive' : 'negative'" v-if="customerId" />
                    </template>
                </q-input>

                <div class="text-body2 q-mb-sm">Available quantity: {{ cabToSell?.quantity }}</div>
                <q-input v-model.number="sellQuantity" type="number" label="Quantity to sell" dense outlined
                    class="q-mb-md"
                    :rules="[val => val > 0 || 'Qty > 0', val => val <= (cabToSell?.quantity || 0) || 'Not enough stock']"
                    lazy-rules>
                    <template v-slot:prepend><q-icon name="numbers" /></template>
                </q-input>

                <q-separator class="q-my-md" />

                <div class="text-subtitle2 q-mb-sm">Additional Accessories</div>
                <div class="row q-col-gutter-sm">
                    <div class="col-8">
                        <q-select v-model="selectedAccessoryId" :options="accessories" option-value="id"
                            option-label="name" label="Select Accessory" dense outlined emit-value map-options>
                            <template v-slot:option="scope">
                                <q-item v-bind="scope.itemProps">
                                    <q-item-section>
                                        <q-item-label>{{ scope.opt.name }}</q-item-label>
                                        <q-item-label caption>₱ {{ scope.opt.price.toLocaleString('en-PH') }} | Avail:
                                            {{ scope.opt.quantity }}</q-item-label>
                                    </q-item-section>
                                </q-item>
                            </template>
                            <template v-slot:no-option>
                                <q-item><q-item-section class="text-grey">No accessories
                                        available</q-item-section></q-item>
                            </template>
                        </q-select>
                    </div>
                    <div class="col-4">
                        <q-input v-model.number="accessoryQuantity" type="number" min="1" label="Quantity" dense
                            outlined :disable="selectedAccessoryId === null" />
                    </div>
                    <div class="col-12">
                        <q-btn color="primary" icon="add" label="Add Accessory" class="full-width"
                            :disable="selectedAccessoryId === null || accessoryQuantity <= 0" @click="addAccessory" />
                    </div>
                </div>

                <div v-if="selectedAccessories.length > 0" class="q-mt-md">
                    <q-list bordered separator>
                        <q-item v-for="acc in selectedAccessories" :key="acc.id">
                            <q-item-section>
                                <q-item-label>{{ acc.name }}</q-item-label>
                                <q-item-label caption>Qty: {{ acc.quantity }} | Price: ₱ {{
                                    acc.price.toLocaleString('en-PH')
                                }}</q-item-label>
                            </q-item-section>
                            <q-item-section side>
                                <div class="text-subtitle2">₱ {{ (acc.price * acc.quantity).toLocaleString('en-PH') }}
                                </div>
                            </q-item-section>
                            <q-item-section side><q-btn flat round dense color="negative" icon="close"
                                    @click="removeAccessory(acc.id)" /></q-item-section>
                        </q-item>
                    </q-list>
                </div>

                <q-separator class="q-my-md" />

                <div class="row justify-between q-mt-md">
                    <div class="text-subtitle2">Cab Total:</div>
                    <div>₱ {{ ((cabToSell?.price || 0) * sellQuantity).toLocaleString('en-PH') }}</div>
                </div>
                <div class="row justify-between q-mt-sm">
                    <div class="text-subtitle2">Accessories Total:</div>
                    <div>₱ {{ totalAccessoriesPrice.toLocaleString('en-PH') }}</div>
                </div>
                <div class="row justify-between q-mt-sm text-bold">
                    <div class="text-subtitle2">Grand Total:</div>
                    <div>₱ {{ totalPrice.toLocaleString('en-PH') }}</div>
                </div>
            </q-card-section>

            <q-card-actions align="right">
                <q-btn flat label="Cancel" @click="closeDialog" />
                <q-btn flat label="Sell" color="primary" @click="handleSellClick"
                    :disable="!isCustomerIdValid || sellQuantity <= 0 || sellQuantity > (cabToSell?.quantity || 0)" />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<style scoped lang="sass">
/* Add any specific styles if needed */
</style>
