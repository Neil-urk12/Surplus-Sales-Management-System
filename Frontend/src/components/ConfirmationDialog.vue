<script setup lang="ts">
interface Props {
  modelValue: boolean;
  title?: string;
  message?: string;
}

withDefaults(defineProps<Props>(), {
  title: 'Confirm Action',
  message: 'Are you sure you want to proceed?',
});

const emit = defineEmits(['update:modelValue', 'confirm']);

const closeDialog = () => {
  emit('update:modelValue', false);
};

const confirmAction = () => {
  emit('confirm');
  closeDialog();
};
</script>

<template>
  <q-dialog :model-value="modelValue" @update:model-value="closeDialog" persistent>
    <q-card style="min-width: 350px">
      <q-card-section>
        <div class="text-h6">{{ title }}</div>
      </q-card-section>

      <q-card-section class="q-pt-none">
        {{ message }}
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Cancel" color="primary" @click="closeDialog" />
        <q-btn flat label="Confirm" color="negative" @click="confirmAction" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
