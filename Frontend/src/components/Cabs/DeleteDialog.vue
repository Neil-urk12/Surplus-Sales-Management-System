<script setup lang="ts">
import { defineAsyncComponent } from 'vue';

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
        @update:modelValue="(val) => emit('update:modelValue', val)"
        @confirm-delete="() => emit('confirm-delete')"
    />
</template>
