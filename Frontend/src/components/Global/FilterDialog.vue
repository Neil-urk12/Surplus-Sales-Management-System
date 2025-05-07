<script setup lang="ts">
import { ref, watch } from 'vue';
import type { PropType } from 'vue';

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
        type: Array as PropType<{
            key: string;
            label: string;
            options: string[];
            modelValue: string | null;
        }[]>,
        required: true,
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
