<template>
  <q-page class="flex flex-center unauthorized-page">
    <div class="text-center">
      <q-icon name="warning" size="xl" color="negative" />
      <div class="text-h4 q-mt-md" :style="{ color: 'var(--foreground)' }">Unauthorized Access</div>
      <div class="text-subtitle1 q-mt-sm" :style="{ color: 'var(--foreground)' }">You do not have permission to view this page.</div>
      <q-btn
        class="q-mt-xl"
        color="primary"
        unelevated
        to="/"
        label="Go Home"
        no-caps
      />
    </div>
  </q-page>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'UnauthorizedPage',
});
</script>

<style lang="scss" scoped>
.unauthorized-page {
  background-color: var(--background);
}
</style>