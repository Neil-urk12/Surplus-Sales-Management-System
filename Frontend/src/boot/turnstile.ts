// src/boot/turnstile.t
import { boot } from 'quasar/wrappers'

export default boot(() => {
  const script = document.createElement('script')
  script.src   = 'https://challenges.cloudflare.com/turnstile/v0/api.js'
  script.async = true
  script.defer = true
  document.head.appendChild(script)
})
