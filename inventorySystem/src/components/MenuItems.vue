<template>
  <q-item
    clickable
    tag="a"
    target="_blank"
  >
    <q-item-section
      v-if="icon"
      avatar
    >
      <q-icon :name="icon" />
    </q-item-section>

    <q-item-section>
      <q-item-label>{{ title }}</q-item-label>
    </q-item-section>
  </q-item>
</template>

<script setup lang="ts">
export interface menuItemsProps {
  title: string;
  icon?: string;
};

withDefaults(defineProps<menuItemsProps>(), {
  caption: '',
  icon: '',
});
</script>
