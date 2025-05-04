<template>
  <q-card bordered class="metrics-card">
    <q-card-section>
      <div class="text-caption text-uppercase">{{ title }}</div>
      <div class="row items-center">
        <q-icon :name="icon" size="sm" class="q-mr-sm" />
      </div>
      <div class="text-h4">{{ value }}</div>
      <div class="text-caption">{{ subtitle }}</div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
defineProps<{
  title: string;
  icon: string;
  value: string | number;
  subtitle: string;
}>();
</script>

<style scoped>
.metrics-card {
  transition: all 0.3s;
}
.metrics-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.12);
}
</style> 