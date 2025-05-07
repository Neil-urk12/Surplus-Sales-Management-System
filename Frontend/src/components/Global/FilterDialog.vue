<script setup lang="ts">
import { ref, watch } from 'vue';
import type { PropType } from 'vue';

/**
 * Interface defining filter data format
 */
interface FilterItemData {
    label: string;
    options: string[];
    value: string | null;
}

interface FilterData {
    [key: string]: FilterItemData;
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
    filterData: {
        type: Object as PropType<FilterData>,
        required: true
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
    Object.keys(props.filterData).forEach(key => {
        newFilters[key] = props.filterData[key].value;
    });
    localFilters.value = newFilters;
}

// Initialize on component creation
initLocalFilters();

// Watch for external changes to filter data
watch(() => props.filterData, () => {
    initLocalFilters();
}, {
    deep: true,
    immediate: true,
    key: 'filter-dialog-data-changes'
});

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
                <div v-for="(filterItem, key) in filterData" :key="key" class="q-mb-md">
                    <q-select v-model="localFilters[key]" :options="['All', ...filterItem.options]"
                        :label="filterItem.label" clearable outlined dense />
                </div>
            </q-card-section>

            <q-card-actions align="right" class="text-primary q-pa-md">
                <q-btn flat label="Reset" color="negative" @click="reset" />
                <q-btn flat label="Apply Filters" @click="apply" />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>
