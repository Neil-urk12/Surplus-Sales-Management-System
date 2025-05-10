<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';

defineProps({
  modelValue: {
    type: Boolean,
    required: true
  },
  title: {
    type: String,
    required: true
  },
  message: {
    type: String,
    required: true
  },
  confirmButtonLabel: {
    type: String,
    default: 'Confirm'
  },
  confirmButtonColor: {
    type: String,
    default: 'primary'
  },
  icon: {
    type: String,
    default: 'warning'
  },
  iconColor: {
    type: String,
    default: 'negative'
  }
});

const emit = defineEmits(['update:modelValue', 'confirm']);

function onConfirm() {
  emit('confirm');
}

function onHide() {
  emit('update:modelValue', false);
}
</script>

<template>
  <q-dialog :model-value="modelValue" persistent @hide="onHide">
    <q-card>
      <q-card-section class="row items-center">
        <q-avatar :icon="icon" :color="iconColor" text-color="white" />
        <span class="q-ml-sm text-h6">{{ title }}</span>
      </q-card-section>

      <q-card-section>
        {{ message }}
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Cancel" @click="onHide" />
        <q-btn flat :label="confirmButtonLabel" :color="confirmButtonColor" @click="onConfirm" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
