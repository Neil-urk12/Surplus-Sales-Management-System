<template>
  <div>
    <!-- Regular menu item -->
    <q-item
      v-if="!hasSubmenu"
      clickable
      v-bind="$attrs"
      :to="to"
      class="text-soft-light q-hoverable"
      :active="$route.path === to"
      :class="{ 'text-dark': $route.path === to && !isDark, 'text-white': $route.path === to && isDark }"
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
      :key="title"
      :icon="icon"
      :label="title"
      class="text-soft-light q-hoverable"
      header-class="text-soft-light"
      expand-icon-class="text-soft-light"
    >
      <div>
        <q-item
          v-for="child in children"
          :key="`${child.title}-${child.to}`"
          :to="child.to"
          clickable
          class="text-soft-light q-hoverable"
          :class="{ 'text-dark': $route.path === child.to && !isDark, 'text-white': $route.path === child.to && isDark }"
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

withDefaults(defineProps<menuItemsProps & { isDark: boolean }>(), {
  icon: '',
  to: '/',
  hasSubmenu: false,
  children: () => [],
  isDark: false
})
</script>

<style scoped>
.q-expansion-item :deep(.q-expansion-item__content) {
  padding-left: 16px;
}
</style>
