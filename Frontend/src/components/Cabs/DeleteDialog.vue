<script setup lang="ts">
import { defineAsyncComponent } from 'vue';
import type { PropType } from 'vue';
const GlobalDeleteDialog = defineAsyncComponent(() => import('src/components/Global/DeleteDialog.vue'));

defineProps({
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
    confirmText: {
        type: String,
        default: 'Delete',
    },
    cancelText: {
        type: String,
        default: 'Cancel',
    },
    iconName: {
        type: String, 
        default: 'warning',
    },
    onConfirmDelete: {
        type: Function as PropType<() => void>,
        required: true,
        description: 'Callback function that will be called when deletion is confirmed'
    }
});

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'confirm-delete'): void
}>();
</script>

<template>
    <GlobalDeleteDialog
        :model-value="modelValue"
        :item-type="itemType"
        :item-name="itemName"
        :confirm-text="confirmText"
        :cancel-text="cancelText"
        :icon-name="iconName"
        :on-confirm-delete="onConfirmDelete"
        @update:modelValue="(val) => emit('update:modelValue', val)"
        @confirm-delete="() => {
            emit('confirm-delete');
            onConfirmDelete();
        }"
    />
</template>
