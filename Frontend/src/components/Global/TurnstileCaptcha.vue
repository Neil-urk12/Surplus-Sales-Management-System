<!-- src/components/TurnstileCaptcha.vue -->
<template>
  <div ref="widgetEl"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, defineEmits } from 'vue'

const widgetEl = ref<HTMLElement|null>(null)
const emit = defineEmits<{
  (e: 'verify', token: string): void
  (e: 'error'): void
  (e: 'expired'): void
}>()

onMounted(() => {
  if (!widgetEl.value || !window.turnstile) {
    console.error('Turnstile script not loaded yet')
    emit('error')
    return
  }

  window.turnstile.render(widgetEl.value, {
    sitekey: import.meta.env.VITE_TURNSTILE_SITE_KEY!,
    callback: (token: string) => emit('verify', token),
    'error-callback': ()    => emit('error'),
    'expired-callback': ()  => emit('expired'),
    theme: 'light',
  })
})
</script>
