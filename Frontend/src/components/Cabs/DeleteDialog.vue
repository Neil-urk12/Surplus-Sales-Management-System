<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps({
    modelValue: {
        type: Boolean,
        required: true,
    },
    itemType: {
        type: String,
        default: 'item',
    },
    itemName: {
        type: String,
        default: 'this item',
    },
});

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'confirm-delete'): void
}>();

const capitalizedItemType = computed(() => {
    return props.itemType.charAt(0).toUpperCase() + props.itemType.slice(1);
});

function confirm() {
    emit('confirm-delete');
    closeDialog(); // Usually close after confirm
}

function closeDialog() {
    emit('update:modelValue', false);
}

</script>

<template>
    <q-dialog :model-value="modelValue" persistent @update:model-value="closeDialog">
        <q-card>
            <q-card-section class="row items-center">
                <q-avatar icon="warning" color="negative" text-color="white" />
                <span class="q-ml-sm text-h6">Delete {{ capitalizedItemType }}</span>
            </q-card-section>

            <q-card-section>
                Are you sure you want to delete <span class="text-weight-bold">{{ itemName }}</span>? This action cannot
                be undone.
            </q-card-section>

            <q-card-actions align="right" class="q-pa-md">
                <q-btn flat label="Cancel" @click="closeDialog" />
                <q-btn flat label="Delete" color="negative" @click="confirm" />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<style scoped lang="sass">
/* Add any specific styles if needed */
</style>
