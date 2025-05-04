<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';

defineProps({
  modelValue: Boolean,
  disable: Boolean
});
const emit = defineEmits(['update:modelValue', 'hide', 'update']);

function onHide() {
  emit('update:modelValue', false);
  emit('hide');
}
function onUpdate() {
  emit('update');
}
</script>

<template>
  <q-dialog :modelValue="modelValue" persistent @hide="onHide" @update:modelValue="$emit('update:modelValue', $event)">
    <q-card style="min-width: 400px; max-width: 95vw">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">Edit Material</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>
      <q-card-section>
        <q-form @submit.prevent="onUpdate" class="q-gutter-sm">
          <slot />
        </q-form>
      </q-card-section>
      <q-card-actions align="right" class="q-pa-md">
        <q-btn flat label="Cancel" @click="onHide" />
        <q-btn unelevated color="primary" label="Update Material" @click="onUpdate" :disable="disable" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>