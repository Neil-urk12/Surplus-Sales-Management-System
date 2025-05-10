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

const LOAD_TIMEOUT = 5000 // 5 seconds timeout

onMounted(() => {
  const siteKey = import.meta.env.VITE_TURNSTILE_SITE_KEY
  if (!siteKey) {
    console.error('VITE_TURNSTILE_SITE_KEY environment variable is not defined')
    emit('error')
    return
  }

  // Set up timeout to check if Turnstile loads
  const timeoutId = setTimeout(() => {
    if (!window.turnstile) {
      console.error('Turnstile script failed to load within timeout period')
      emit('error')
    }
  }, LOAD_TIMEOUT)

  // Function to initialize Turnstile once it's available
  const initializeTurnstile = () => {
    if (!widgetEl.value || !window.turnstile) {
      return
    }

    if (typeof window.turnstile.render !== 'function') {
      console.error('Turnstile render function is not available')
      emit('error')
      return
    }

    // Clear the timeout since Turnstile has loaded
    clearTimeout(timeoutId)

    window.turnstile.render(widgetEl.value, {
      sitekey: siteKey,
      callback: (token: string) => emit('verify', token),
      'error-callback': ()    => emit('error'),
      'expired-callback': ()  => emit('expired'),
      theme: 'light',
    })
  }

  // Check immediately in case Turnstile is already loaded
  initializeTurnstile()

  // Set up an observer to watch for Turnstile to become available
  const observer = new MutationObserver(() => {
    if (window.turnstile) {
      initializeTurnstile()
      observer.disconnect()
    }
  })

  // Start observing the document for changes
  observer.observe(document, {
    childList: true,
    subtree: true
  })
})
</script>
