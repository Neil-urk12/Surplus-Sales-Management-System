<script setup lang="ts">
import { computed, ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import * as TurnstileCaptchaComponent from '../components/Global/TurnstileCaptcha.vue'

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
const token = ref<string|null>(null)
const error = ref(false)
const isCaptchad = ref(false)
const showCaptchaDialog = ref(false)

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

const form = reactive({
  email: '',
  password: ''
})

const errors = reactive({
  email: '',
  password: '',
})

const isSubmitting = ref(false)
const showShake = ref(false)
const loginError = ref('')

onMounted( async () => {
  if (authStore.isAuthenticated) await router.push('/')

  // Load Turnstile script if not already loaded
  if (!window.turnstile) {
    const script = document.createElement('script')
    script.src = 'https://challenges.cloudflare.com/turnstile/v0/api.js'
    script.async = true
    script.defer = true
    document.head.appendChild(script)
  }
})

const validateEmail = () => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!form.email) {
    errors.email = 'Email is required'
  } else if (!emailRegex.test(form.email)) {
    errors.email = 'Please enter a valid email'
  } else {
    errors.email = ''
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

const handleSubmit = async () => {
  validateEmail()
  validatePassword()
  loginError.value = ''

  // First check if fields are valid
  if (errors.email || errors.password) {
    showShake.value = true
    setTimeout(() => showShake.value = false, 500)
    return
  }

  // Only after fields are valid, check for captcha
  if (!token.value) {
    showCaptchaDialog.value = true
    isCaptchad.value = true
    return
  }

  isSubmitting.value = true

  try {
    const result = await authStore.login(form)

    if (result.success) {
      await router.push('/')
    } else {
      let msg = result.message || 'Invalid email or password. Please try again.'
      if (msg === 'Invalid credentials') {
        msg = 'Invalid email or password.'
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
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-container" :class="{ 'shake-animation': showShake }">
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
          <input
            v-model="form.email"
            id="email"
            placeholder="Enter your email"
            @blur="validateEmail"
            :class="{ 'error': errors.email }"
          >
          <span class="error-message" v-if="errors.email">{{ errors.email }}</span>
        </div>

        <div class="input-group">
          <input
            v-model="form.password"
            type="password"
            id="password"
            placeholder="Enter your password"
            @blur="validatePassword"
            :class="{ 'error': errors.password }"
          >
          <span class="error-message" v-if="errors.password">{{ errors.password }}</span>
          <div class="row">
            <div class="captcha-container col" v-if="form.email && form.password">
              <TurnstileCaptchaComponent.default
                @verify="token = $event"
                @error="error = true"
                @expired="token = null"
              />
              <div class="col flex content-center justify-end">
                <a href="#" class="forgot-password" @click.prevent="showModal = true">Forgot password?</a>
              </div>
            </div>
            <div class="col" v-else>
              <div class="flex content-center justify-end">
                <a href="#" class="forgot-password" @click.prevent="showModal = true">Forgot password?</a>
              </div>
            </div>
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
          <q-icon
            name="close"
            class="close-icon"
            @click="showModal = false"
          />
        </q-card-section>

        <q-tabs
          v-model="tab"
          dense
          class="bg-grey-1 text-dark"
          active-color="primary"
          indicator-color="primary"
          align="justify"
        >
          <q-tab name="email" icon="mail" label="Email" />
          <q-tab name="phone" icon="phone" label="Phone" />
        </q-tabs>

        <q-separator />

        <q-card-section class="content-section">
          <div class="illustration-container">
            <q-img
              src="https://cdn-icons-png.flaticon.com/512/6195/6195699.png"
              alt="Password Recovery"
              style="width: 120px; height: 120px"
            />
          </div>

          <div class="text-center text-body1 q-mb-md">
            <template v-if="tab === 'email'">
              Enter your email address and we'll send you a link to reset your password
            </template>
            <template v-else>
              Enter your phone number and we'll send you a verification code
            </template>
          </div>

          <q-input
            v-if="tab === 'email'"
            v-model="recoveryEmail"
            type="email"
            label="Email Address"
            outlined
            dense
            class="q-mb-md"
            :rules="[val => !!val || 'Email is required', val => /.+@.+\..+/.test(val) || 'Invalid email']"
          >
            <template v-slot:prepend>
              <q-icon name="mail" />
            </template>
          </q-input>

          <q-input
            v-else
            v-model="recoveryPhone"
            type="tel"
            label="Phone Number"
            mask="(###) ###-####"
            fill-mask
            outlined
            dense
            class="q-mb-md"
            :rules="[val => val.replace(/\D/g,'').length === 10 || 'Valid phone number required']"
          >
            <template v-slot:prepend>
              <q-icon name="phone" />
            </template>
          </q-input>

          <q-btn
            label="Send Recovery Link"
            color="primary"
            class="full-width q-mt-md"
            :loading="isSendingRecovery"
            :disable="!isRecoveryInputValid"
            @click="handleRecoveryRequest"
          />

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

    <q-dialog v-model="showCaptchaDialog">
      <q-card >
        <q-card-section>
          <div class="text-h6">CAPTCHA Verification required</div>
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
.shake-animation {
  animation: shake 0.5s;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-5px); }
  40%, 80% { transform: translateX(5px); }
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
  to { transform: rotate(360deg); }
}
</style>
