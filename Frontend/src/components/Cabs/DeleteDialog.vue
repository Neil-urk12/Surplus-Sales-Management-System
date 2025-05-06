<script setup lang="ts">
import GlobalDeleteDialog from 'src/components/Global/DeleteDialog.vue';

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
        v-bind="$props"
        @update:modelValue="(val) => emit('update:modelValue', val)"
        @confirm-delete="() => emit('confirm-delete')"
    />
</template>
