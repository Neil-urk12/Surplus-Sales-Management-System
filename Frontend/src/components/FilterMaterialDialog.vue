<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';

defineProps({ modelValue: Boolean });
const emit = defineEmits(['update:modelValue']);

function onHide() {
  emit('update:modelValue', false);
}
function onUpdate(val: boolean) {
  emit('update:modelValue', val);
}
</script>

<template>
  <q-dialog :model-value="modelValue" @update:model-value="onUpdate" @hide="onHide">
    <q-card style="min-width: 350px">
      <q-card-section class="row items-center">
        <div class="text-h6">Filter Materials</div>
        <q-space />
        <q-btn icon="close" flat round dense @click="onHide" />
      </q-card-section>
      <q-card-section class="q-pt-none">
        <slot />
      </q-card-section>
      <q-card-actions align="right" class="text-primary">
        <slot name="actions" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>