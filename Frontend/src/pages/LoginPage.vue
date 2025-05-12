<script setup lang="ts">
import { computed, ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import * as TurnstileCaptchaComponent from '../components/Global/TurnstileCaptcha.vue'
import { useQuasar } from 'quasar'

// Define Turnstile interface for TypeScript
interface TurnstileOptions {
  sitekey: string;
  callback: (token: string) => void;
  'error-callback': () => void;
  'expired-callback': () => void;
  theme?: 'light' | 'dark' | 'auto';
}

interface Turnstile {
  render: (element: HTMLElement, options: TurnstileOptions) => string;
}

// Extend Window interface
declare global {
  interface Window {
    turnstile?: Turnstile;
  }
}

const tab = ref('email')
const showModal = ref(false)
const recoveryEmail = ref('')
const recoveryPhone = ref('')
const showCaptchaDialog = ref(false)
const captchaToken = ref<string | null>(null)

const isSendingRecovery = ref(false)
const recoverySent = ref(false)

const resetRecoveryForm = () => {
  recoveryEmail.value = ''
  recoveryPhone.value = ''
  recoverySent.value = false
  tab.value = 'email'
}

const handleRecoveryRequest = async () => {
  isSendingRecovery.value = true

  try {
    await new Promise(resolve => setTimeout(resolve, 1500))
    recoverySent.value = true

    setTimeout(() => {
      showModal.value = false
    }, 3000)
  } catch (error) {
    console.error('Recovery error:', error)
  } finally {
    isSendingRecovery.value = false
  }
}

const isRecoveryInputValid = computed(() => {
  return tab.value === 'email'
    ? recoveryEmail.value.trim() !== ''
    : recoveryPhone.value.trim().length >= 10
})

const router = useRouter()
const authStore = useAuthStore()
const $q = useQuasar()

const form = reactive({
  username: '',
  password: ''
})

const isPageLoading = ref(true)

const errors = reactive({
  username: '',
  password: '',
})

const isSubmitting = ref(false)
const showShake = ref(false)
const loginError = ref('')
const isCaptchaLoading = ref(false)
const isLoginValidated = ref(false)

onMounted(async () => {
  // Show page loader
  isPageLoading.value = true

  try {
    if (authStore.isAuthenticated) await router.push('/')

    // Load Turnstile script if not already loaded
    if (!window.turnstile) {
      const script = document.createElement('script')
      script.src = 'https://challenges.cloudflare.com/turnstile/v0/api.js'
      script.async = true
      script.defer = true
      document.head.appendChild(script)
    }

    // Simulate a minimum loading time for better UX
    await new Promise(resolve => setTimeout(resolve, 800))
  } catch (error) {
    console.error('Error during page initialization:', error)
  } finally {
    // Hide page loader
    isPageLoading.value = false
  }
})

const validateUsername = () => {
  if (!form.username) {
    errors.username = 'Username or email is required'
  } else {
    // Allow either email or username format
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/

    if (!emailRegex.test(form.username) && !usernameRegex.test(form.username)) {
      errors.username = 'Please enter a valid username or email'
    } else {
      errors.username = ''
    }
  }
}

const validatePassword = () => {
  if (!form.password) {
    errors.password = 'Password is required'
  } else if (form.password.length < 6) {
    errors.password = 'Password must be at least 6 characters'
  } else {
    errors.password = ''
  }
}

const handleCaptchaVerified = (captchaValue: string) => {
  captchaToken.value = captchaValue
  showCaptchaDialog.value = false
  void proceedLogin()
}

const handleCaptchaError = () => {
  captchaToken.value = null
  showCaptchaDialog.value = false
}

const proceedLogin = async () => {
  isSubmitting.value = true
  // Show loading overlay during login
  $q.loading.show({
    message: 'Signing in...',
    spinnerColor: 'primary',
    spinnerSize: 80
  })

  try {
    const result = await authStore.login(form)
    if (result.success) {
      isLoginValidated.value = true
      captchaToken.value = null
      await router.push('/')
    } else {
      captchaToken.value = null
      form.password = ''
      let msg = result.message || 'Invalid username/email or password. Please try again.'
      if (msg === 'Invalid credentials') {
        msg = 'Invalid username/email or password.'
      }
      loginError.value = msg
      showShake.value = true
      setTimeout(() => showShake.value = false, 500)
    }
  } catch (error) {
    console.error('Login error:', error)
    loginError.value = 'An unexpected error occurred. Please try again later.'
    showShake.value = true
    setTimeout(() => showShake.value = false, 500)
  } finally {
    isSubmitting.value = false
    // Hide loading overlay
    $q.loading.hide()
  }
}

const handleSubmit = () => {
  validateUsername()
  validatePassword()
  loginError.value = ''

  if (errors.username || errors.password) {
    showShake.value = true
    setTimeout(() => showShake.value = false, 500)
    return
  }

  // Open CAPTCHA dialog instead of proceeding
  showCaptchaDialog.value = true
}

// Watch for when both email and password are filled to set captcha loading
watch(
  () => [form.username, form.password],
  ([username, password], [oldUsername, oldPassword]) => {
    if (username && password && (!oldUsername || !oldPassword)) {
      isCaptchaLoading.value = true
    }
    // Optionally, if either is cleared, reset loading state
    if (!username || !password) {
      isCaptchaLoading.value = false
    }
  }
)
</script>

<template>
  <div class="login-page">
      <div class="login-container" :class="{ 'shake-animation': showShake }">
      <!-- Spinner overlay when page is loading -->
      <div v-if="isPageLoading" class="spinner-overlay">
        <q-spinner-dots color="primary" size="80px" />
      </div>

      <!-- Actual content -->
      <div class="logo">
        <q-icon name="account_circle" size="48px" />
      </div>

      <h1>Welcome back</h1>

      <form @submit.prevent="handleSubmit">
        <div v-if="loginError" class="error-message general-error">
          <q-icon name="error" class="error-icon" />
          {{ loginError }}
        </div>

        <div class="input-group">
          <input v-model="form.username" id="username" placeholder="Enter your username or email"
            @blur="validateUsername" :class="{ 'error': errors.username }">
          <span class="error-message" v-if="errors.username">{{ errors.username }}</span>
        </div>

        <div class="input-group">
          <input v-model="form.password" type="password" id="password" placeholder="Enter your password"
            @blur="validatePassword" :class="{ 'error': errors.password }">
          <span class="error-message" v-if="errors.password">{{ errors.password }}</span>
          <div class="col flex content-center justify-end">
            <a href="#" class="forgot-password" @click.prevent="showModal = true">Forgot password?</a>
          </div>
        </div>

        <button type="submit" :disabled="isSubmitting">
          <span v-if="!isSubmitting" class="text-weight-medium">Sign In</span>
          <span v-else class="spinner"></span>
        </button>

        <div class="divider"></div>
      </form>
    </div>

    <q-dialog v-model="showModal" @hide="resetRecoveryForm">
      <q-card class="recovery-card">
        <q-card-section class="header-section">
          <div class="text-h6">Password Recovery</div>
          <q-icon name="close" class="close-icon" @click="showModal = false" />
        </q-card-section>

        <q-tabs v-model="tab" dense class="bg-grey-1 text-dark" active-color="primary" indicator-color="primary"
          align="justify">
          <q-tab name="email" icon="mail" label="Email" />
          <q-tab name="phone" icon="phone" label="Phone" />
        </q-tabs>

        <q-separator />

        <q-card-section class="content-section">
          <div class="illustration-container">
            <q-img src="https://cdn-icons-png.flaticon.com/512/6195/6195699.png" alt="Password Recovery"
              style="width: 120px; height: 120px" />
          </div>

          <div class="text-center text-body1 q-mb-md">
            <template v-if="tab === 'email'">
              Enter your email address and we'll send you a link to reset your password
            </template>
            <template v-else>
              Enter your phone number and we'll send you a verification code
            </template>
          </div>

          <q-input v-if="tab === 'email'" v-model="recoveryEmail" type="email" label="Email Address" outlined dense
            class="q-mb-md"
            :rules="[val => !!val || 'Email is required', val => /.+@.+\..+/.test(val) || 'Invalid email']">
            <template v-slot:prepend>
              <q-icon name="mail" />
            </template>
          </q-input>

          <q-input v-else v-model="recoveryPhone" type="tel" label="Phone Number" mask="(###) ###-####" fill-mask
            outlined dense class="q-mb-md"
            :rules="[val => val.replace(/\D/g, '').length === 10 || 'Valid phone number required']">
            <template v-slot:prepend>
              <q-icon name="phone" />
            </template>
          </q-input>

          <q-btn label="Send Recovery Link" color="primary" class="full-width q-mt-md" :loading="isSendingRecovery"
            :disable="!isRecoveryInputValid" @click="handleRecoveryRequest" />

          <div class="text-center text-caption q-mt-md">
            Remember your password?
            <a href="#" @click.prevent="showModal = false">Sign in</a>
          </div>
        </q-card-section>

        <transition name="fade">
          <q-card-section v-if="recoverySent" class="success-message">
            <q-icon name="check_circle" color="positive" size="lg" />
            <div class="text-body1 q-mt-sm">
              Recovery email sent successfully!
            </div>
            <div class="text-caption">
              Check your inbox and follow the instructions
            </div>
          </q-card-section>
        </transition>
      </q-card>
    </q-dialog>

    <q-dialog v-model="showCaptchaDialog" persistent>
      <q-card>
        <q-card-section>
          <div class="text-h6">CAPTCHA Verification</div>
        </q-card-section>
        <q-card-section>
          <TurnstileCaptchaComponent.default @verify="handleCaptchaVerified" @error="handleCaptchaError"
            @expired="handleCaptchaError" />
        </q-card-section>
      </q-card>
    </q-dialog>
  </div>
</template>

<style>
@import '../assets/styles/LoginPage.base.css';
@import '../assets/styles/LoginPage.light.css';
@import '../assets/styles/LoginPage.dark.css';

/* Keep only animation-related styles here */

/* Spinner overlay styles */
.spinner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.26);
  z-index: 10;
}

.shake-animation {
  animation: shake 0.5s;
}

@keyframes shake {
  0%,
  100% {
    transform: translateX(0);
  }

  20%,
  60% {
    transform: translateX(-5px);
  }

  40%,
  80% {
    transform: translateX(5px);
  }
}

.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid var(--q-separator-dark);
  border-radius: 50%;
  border-top-color: var(--q-primary);
  animation: spin 1s ease-in-out infinite;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.captcha-container {
  margin: 15px 0;
  display: flex;
  justify-content: center;
}

.general-error {
  background-color: var(--q-negative);
  color: white;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 15px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.error-icon {
  font-size: 20px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.captcha-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px;
  background-color: var(--q-grey-1);
  border-radius: 4px;
  min-height: 100px;
}
</style>
