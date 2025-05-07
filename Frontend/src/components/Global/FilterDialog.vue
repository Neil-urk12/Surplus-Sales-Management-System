<script setup lang="ts">
import { ref, watch } from 'vue';
import type { PropType } from 'vue';

/**
 * Interface defining a filter option
 */
interface FilterOption {
    key: string;
    label: string;
    options: string[];
    modelValue: string | null;
}

/**
 * Validates that all items in an array are strings
 * @param value The array to validate
 * @returns True if all items are strings, false otherwise
 */
function validateStringArray(value: unknown[]): boolean {
    return value.every(item => typeof item === 'string');
}

/**
 * Validates that a FilterOption has all required fields and correct types
 * @param option FilterOption object to validate
 * @returns True if valid, false otherwise
 */
function validateFilterOption(option: FilterOption): boolean {
    return (
        typeof option.key === 'string' &&
        typeof option.label === 'string' &&
        Array.isArray(option.options) &&
        validateStringArray(option.options) &&
        (option.modelValue === null || typeof option.modelValue === 'string')
    );
}

const props = defineProps({
    modelValue: {
        type: Boolean,
        required: true,
    },
    title: {
        type: String,
        default: 'Filter',
    },
    filterOptions: {
        type: Array as PropType<FilterOption[]>,
        required: true,
        validator: (value: FilterOption[]) => value.every(validateFilterOption)
    },
});

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
    (e: 'apply-filters', filters: Record<string, string | null>): void;
    (e: 'reset-filters'): void;
}>();

// Local state for filters within the dialog
const localFilters = ref<Record<string, string | null>>({});

// Initialize local filters from props
function initLocalFilters() {
    const newFilters: Record<string, string | null> = {};
    props.filterOptions.forEach(option => {
        newFilters[option.key] = option.modelValue;
    });
    localFilters.value = newFilters;
}

// Initialize on component creation
initLocalFilters();

// Watch for external changes to filter options
watch(() => props.filterOptions, () => {
    initLocalFilters();
}, { deep: true });

function apply() {
    emit('apply-filters', localFilters.value);
    closeDialog();
}

function reset() {
    const resetFilters: Record<string, string | null> = {};
    Object.keys(localFilters.value).forEach(key => {
        resetFilters[key] = null;
    });
    localFilters.value = resetFilters;
    emit('reset-filters');
}

function closeDialog() {
    emit('update:modelValue', false);
}
</script>

<template>
    <q-dialog :model-value="modelValue" @update:model-value="closeDialog">
        <q-card style="min-width: 350px">
            <q-card-section class="row items-center">
                <div class="text-h6">{{ title }}</div>
                <q-space />
                <q-btn icon="close" flat round dense @click="closeDialog" />
            </q-card-section>

            <q-card-section class="q-pt-none">
                <div v-for="option in filterOptions" :key="option.key" class="q-mb-md">
                    <q-select v-model="localFilters[option.key]" :options="['All', ...option.options]"
                        :label="option.label" clearable outlined dense />
                </div>
            </q-card-section>

            <q-card-actions align="right" class="text-primary q-pa-md">
                <q-btn flat label="Reset" color="negative" @click="reset" />
                <q-btn flat label="Apply Filters" @click="apply" />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>
