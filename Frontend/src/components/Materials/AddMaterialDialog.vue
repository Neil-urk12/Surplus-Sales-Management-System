<script setup lang="ts">
import { ref, watch } from 'vue';
import type { NewMaterialInput } from 'src/types/materials';
import type { PropType } from 'vue';
import { operationNotifications } from '../../utils/notifications';
import MaterialFormCommon from './MaterialFormCommon.vue';

const props = defineProps({
    modelValue: {
        type: Boolean,
        required: true,
    },
    categories: {
        type: Array as PropType<string[]>,
        required: true,
    },
    suppliers: {
        type: Array as PropType<string[]>,
        required: true,
    },
    defaultImageUrl: {
        type: String,
        required: true,
    },
});

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'add-material', materialData: NewMaterialInput): void
}>();

// --- State --- 
const newMaterial = ref<NewMaterialInput>({
    name: '',
    category: '',
    supplier: '',
    quantity: 0,
    status: 'Out of Stock',
    image: props.defaultImageUrl,
});
const isProcessing = ref(false);

// --- Watchers --- 
watch(() => props.modelValue, (newValue) => {
    if (newValue) {
        resetForm();
    }
});

// Add watcher for quantity to update status automatically
watch(() => newMaterial.value.quantity, (newQuantity) => {
    if (newQuantity === 0) {
        newMaterial.value.status = 'Out of Stock';
    } else if (newQuantity <= 10) {
        newMaterial.value.status = 'Low Stock';
    } else {
        newMaterial.value.status = 'In Stock';
    }
});

// --- Functions --- 
function resetForm() {
    newMaterial.value = {
        name: '',
        category: '',
        supplier: '',
        quantity: 0,
        status: 'Out of Stock',
        image: props.defaultImageUrl,
    };
}

function handleSubmit() {
    isProcessing.value = true;
    try {
        // Don't make the API call here, just emit the event
        // Let the parent component (MaterialsPage) handle the API call
        emit('add-material', { ...newMaterial.value });
        emit('update:modelValue', false);
    } catch (error) {
        console.error('Error in add material form submission:', error);
        operationNotifications.add.error('material');
    } finally {
        isProcessing.value = false;
    }
}

function closeDialog() {
    emit('update:modelValue', false);
}

function handleValidationError(message: string) {
    operationNotifications.validation.error(message);
}
</script>

<template>
    <q-dialog :model-value="modelValue" persistent @update:model-value="closeDialog" @hide="resetForm">
        <q-card style="min-width: 400px; max-width: 95vw">
            <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">New Material</div>
                <q-space />
                <q-btn icon="close" flat round dense @click="closeDialog" />
            </q-card-section>

            <q-card-section>
                <MaterialFormCommon
                    v-model:materialData="newMaterial"
                    :categories="categories"
                    :suppliers="suppliers"
                    :defaultImageUrl="defaultImageUrl"
                    :isProcessing="isProcessing"
                    @submit="handleSubmit"
                    @cancel="closeDialog"
                    @validation-error="handleValidationError"
                >
                    <template #submitButton="slotProps">
                        <q-btn unelevated color="primary" label="Add Material" type="submit" :loading="slotProps.loading" :disable="slotProps.disabled" />
                    </template>
                </MaterialFormCommon>
            </q-card-section>
        </q-card>
    </q-dialog>
</template> 