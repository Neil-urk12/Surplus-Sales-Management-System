<template>
  <q-item
    clickable
    v-bind="$attrs"
    :to="to"
    class="text-soft-light q-hoverable"
    :active="$route.path === to"
    :exact="exact"
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
import type { menuItemsProps } from '../types/menu-items'

withDefaults(defineProps<menuItemsProps>(), {
  icon: '',
  to: '/',
})
</script>
