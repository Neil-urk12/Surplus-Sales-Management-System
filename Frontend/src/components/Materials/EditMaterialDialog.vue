<script setup lang="ts">
import { ref, watch } from 'vue';
import type { PropType } from 'vue';
import type { MaterialRow, NewMaterialInput } from 'src/types/materials';
import { operationNotifications } from '../../utils/notifications';
import MaterialFormCommon from './MaterialFormCommon.vue';

const props = defineProps({
    modelValue: {
        type: Boolean,
        required: true,
    },
    materialData: {
        type: Object as PropType<MaterialRow | null>,
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
    (e: 'update-material', materialData: NewMaterialInput): void
}>();

// --- Local State --- 
const localMaterialData = ref<NewMaterialInput>({
    name: '',
    category: '',
    supplier: '',
    quantity: 0,
    status: 'Out of Stock',
    image: props.defaultImageUrl,
});

const isProcessing = ref(false);

// Watch for changes in materialData
watch(() => props.materialData, (newMaterial) => {
    console.log('materialData changed:', newMaterial);
    if (newMaterial) {
        // Reset local state based on the new material data
        const initialImage = newMaterial.image || props.defaultImageUrl;
        
        localMaterialData.value = {
            name: newMaterial.name,
            category: newMaterial.category,
            supplier: newMaterial.supplier,
            quantity: newMaterial.quantity,
            status: newMaterial.status,
            image: initialImage
        };
    }
}, { immediate: true });

// Add watcher for quantity to update status automatically
watch(() => localMaterialData.value.quantity, (newQuantity) => {
    if (newQuantity === 0) {
        localMaterialData.value.status = 'Out of Stock';
    } else if (newQuantity <= 10) {
        localMaterialData.value.status = 'Low Stock';
    } else {
        localMaterialData.value.status = 'In Stock';
    }
});

// --- Functions --- 
function resetForm() {
    console.log('Resetting form in EditMaterialDialog');
    if (props.materialData) {
        localMaterialData.value = {
            name: props.materialData.name,
            category: props.materialData.category,
            supplier: props.materialData.supplier,
            quantity: props.materialData.quantity,
            status: props.materialData.status,
            image: props.materialData.image || props.defaultImageUrl,
        };
    } else {
        localMaterialData.value = {
            name: '',
            category: '',
            supplier: '',
            quantity: 0,
            status: 'Out of Stock',
            image: props.defaultImageUrl,
        };
    }
}

function handleSubmit() {
    isProcessing.value = true;
    try {
        // Don't make the API call here, just emit the event
        // Let the parent component handle the API call
        emit('update-material', { ...localMaterialData.value });
        emit('update:modelValue', false);
    } catch (error) {
        console.error('Error in edit material form submission:', error);
        operationNotifications.update.error('material');
    } finally {
        isProcessing.value = false;
    }
}

function closeDialog() {
    console.log('closeDialog called in EditMaterialDialog');
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
                <div class="text-h6">Edit Material</div>
                <q-space />
                <q-btn icon="close" flat round dense @click="closeDialog" />
            </q-card-section>

            <q-card-section>
                <MaterialFormCommon
                    v-model:materialData="localMaterialData"
                    :categories="categories"
                    :suppliers="suppliers"
                    :defaultImageUrl="defaultImageUrl"
                    :isProcessing="isProcessing"
                    @submit="handleSubmit"
                    @cancel="closeDialog"
                    @validation-error="handleValidationError"
                >
                    <template #submitButton="slotProps">
                        <q-btn unelevated color="primary" label="Update Material" type="submit" :loading="slotProps.loading" :disable="slotProps.disabled" />
                    </template>
                </MaterialFormCommon>
            </q-card-section>
        </q-card>
    </q-dialog>
</template> 