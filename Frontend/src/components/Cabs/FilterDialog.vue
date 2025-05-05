<script setup lang="ts">
import { ref, watch, PropType } from 'vue';
import type { CabStatus } from 'src/types/cabs';

const props = defineProps({
    modelValue: {
        type: Boolean,
        required: true,
    },
    makes: {
        type: Array as PropType<string[]>,
        required: true,
    },
    colors: {
        type: Array as PropType<string[]>,
        required: true,
    },
    statuses: {
        type: Array as PropType<CabStatus[]>,
        required: true,
    },
    // Pass initial values to sync with parent store state
    initialFilterMake: { type: String as PropType<string | null>, default: null },
    initialFilterColor: { type: String as PropType<string | null>, default: null },
    initialFilterStatus: { type: String as PropType<CabStatus | null>, default: null },
});

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'apply-filters', filters: { make: string | null; color: string | null; status: CabStatus | null }): void
    (e: 'reset-filters'): void
}>();

// Local state for filters within the dialog
const localFilterMake = ref<string | null>(props.initialFilterMake);
const localFilterColor = ref<string | null>(props.initialFilterColor);
const localFilterStatus = ref<CabStatus | null>(props.initialFilterStatus);

// Watch props to update local state if dialog is reopened with different initial values
watch(() => [props.initialFilterMake, props.initialFilterColor, props.initialFilterStatus], () => {
    localFilterMake.value = props.initialFilterMake;
    localFilterColor.value = props.initialFilterColor;
    localFilterStatus.value = props.initialFilterStatus;
});

function apply() {
    emit('apply-filters', {
        make: localFilterMake.value,
        color: localFilterColor.value,
        status: localFilterStatus.value,
    });
    closeDialog();
}

function reset() {
    localFilterMake.value = null;
    localFilterColor.value = null;
    localFilterStatus.value = null;
    emit('reset-filters');
    // Keep dialog open after reset or close? Current CabsPage closes.
    // closeDialog(); // Uncomment if reset should also close the dialog
}

function closeDialog() {
    emit('update:modelValue', false);
}

</script>

<template>
    <q-dialog :model-value="modelValue" @update:model-value="closeDialog">
        <q-card style="min-width: 350px">
            <q-card-section class="row items-center">
                <div class="text-h6">Filter Cabs</div>
                <q-space />
                <q-btn icon="close" flat round dense @click="closeDialog" />
            </q-card-section>

            <q-card-section class="q-pt-none">
                <q-select v-model="localFilterMake" :options="makes" label="Make" clearable outlined dense
                    class="q-mb-md" />

                <q-select v-model="localFilterColor" :options="colors" label="Color" clearable outlined dense
                    class="q-mb-md" />

                <q-select v-model="localFilterStatus" :options="statuses" label="Status" clearable outlined dense
                    class="q-mb-md" />
            </q-card-section>

            <q-card-actions align="right" class="text-primary q-pa-md">
                <q-btn flat label="Reset" color="negative" @click="reset" />
                <q-btn flat label="Apply Filters" @click="apply" />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<style scoped lang="sass">
/* Add any specific styles if needed */
</style>
