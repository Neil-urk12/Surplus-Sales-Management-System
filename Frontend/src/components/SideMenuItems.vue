<template>
  <div>
    <!-- Regular menu item -->
    <q-item
      v-if="!isDropdown"
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

    <!-- Dropdown menu item -->
    <q-expansion-item
      v-else
      :icon="icon"
      :label="title"
      class="text-soft-light q-hoverable"
      header-class="text-soft-light"
      expand-icon-class="text-soft-light"
      dense
    >
      <div>
        <q-item
          v-for="child in children"
          :key="child.title"
          :to="child.to"
          clickable
          class="text-soft-light q-hoverable"
        >
          <q-item-section avatar v-if="child.icon">
            <q-icon :name="child.icon" />
          </q-item-section>

          <q-item-section>
            <q-item-label>{{ child.title }}</q-item-label>
          </q-item-section>
        </q-item>
      </div>
    </q-expansion-item>
  </div>
</template>

<script setup lang="ts">
import type { menuItemsProps } from '../types/menu-items'

withDefaults(defineProps<menuItemsProps>(), {
  icon: '',
  to: '/',
  isDropdown: false,
  children: () => []
})
</script>

<style scoped>
.q-expansion-item :deep(.q-expansion-item__content) {
  padding-left: 16px;
}
</style>
