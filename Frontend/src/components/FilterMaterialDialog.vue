<script setup lang="ts">
import { ref, watch } from 'vue';
import type { PropType } from 'vue';
import type { MaterialStatus } from 'src/types/materials';

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
  statuses: {
    type: Array as PropType<MaterialStatus[]>,
    required: true,
  },
  // Pass initial values to sync with parent store state
  initialFilterCategory: {
    type: String as PropType<string | null>,
    default: null
  },
  initialFilterSupplier: {
    type: String as PropType<string | null>,
    default: null
  },
  initialFilterStatus: {
    type: String as PropType<MaterialStatus | null>,
    default: null
  },
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'apply-filters', filters: { category: string | null; supplier: string | null; status: MaterialStatus | null }): void
  (e: 'reset-filters'): void
}>();

// Local state for filters within the dialog
const localFilterCategory = ref<string | null>(props.initialFilterCategory);
const localFilterSupplier = ref<string | null>(props.initialFilterSupplier);
const localFilterStatus = ref<MaterialStatus | null>(props.initialFilterStatus);

// Watch props to update local state if dialog is reopened with different initial values
watch(() => [props.initialFilterCategory, props.initialFilterSupplier, props.initialFilterStatus], () => {
  localFilterCategory.value = props.initialFilterCategory;
  localFilterSupplier.value = props.initialFilterSupplier;
  localFilterStatus.value = props.initialFilterStatus;
});

function apply() {
  emit('apply-filters', {
    category: localFilterCategory.value,
    supplier: localFilterSupplier.value,
    status: localFilterStatus.value,
  });
  closeDialog();
}

function reset() {
  localFilterCategory.value = null;
  localFilterSupplier.value = null;
  localFilterStatus.value = null;
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
        <div class="text-h6">Filter Materials</div>
        <q-space />
        <q-btn icon="close" flat round dense @click="closeDialog" />
      </q-card-section>

      <q-card-section class="q-pt-none">
        <q-select v-model="localFilterCategory" :options="['All', ...categories]" label="Category" clearable outlined
          dense class="q-mb-md" />

        <q-select v-model="localFilterSupplier" :options="['All', ...suppliers]" label="Supplier" clearable outlined
          dense class="q-mb-md" />

        <q-select v-model="localFilterStatus" :options="['All', ...statuses]" label="Status" clearable outlined dense
          class="q-mb-md" />
      </q-card-section>

      <q-card-actions align="right" class="text-primary q-pa-md">
        <q-btn flat label="Reset" color="negative" @click="reset" />
        <q-btn flat label="Apply Filters" @click="apply" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
