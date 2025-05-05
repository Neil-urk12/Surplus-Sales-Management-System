<script setup lang="ts">
import { computed, ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const tab = ref('email')
const showModal = ref(false)
const recoveryEmail = ref('')
const recoveryPhone = ref('')

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

  if (errors.email || errors.password) {
    showShake.value = true
    setTimeout(() => showShake.value = false, 500)
    return
  }

  isSubmitting.value = true

  try {
    const success = await authStore.login(form)

    if (success) {
      await router.push('/')
    } else {
      loginError.value = 'Invalid email or password. Please try again.'
      showShake.value = true
      setTimeout(() => showShake.value = false, 500)
    }
  } catch (error) {
    console.error('Login error:', error)
    loginError.value = 'An unexpected error occurred. Please try again later.'
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
          <a href="#" class="forgot-password" @click.prevent="showModal = true">Forgot password?</a>
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
  </div>
</template>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  line-height: 1.6;
  padding: 20px;
}

.login-container {
  padding: 2.5rem;
  border-radius: 12px;
  border: 1px solid;
  width: 100%;
  max-width: 500px;
  transition: all 0.3s ease;
}

.body--light .login-page {
  background: #ffffff;
}

.body--light .login-container {
  background: #ffffff;
  border: 1px solid #e0e0e0;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
}

.body--light input {
  background: #ffffff;
  border: 1px solid #e0e0e0;
  color: #000000;
}

.body--light input::placeholder {
  color: rgba(0, 0, 0, 0.4);
}

.body--light input:focus {
  border-color: #000000;
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.1);
}

.body--light h1 {
  color: #000000;
}

.body--light button {
  background-color: #18181b;
  color: #ffffff;
  border: none;
  font-weight: 500;
}

.body--light button:hover {
  background-color: #000000;
}

.body--light button:disabled {
  background-color: rgba(0, 0, 0, 0.5);
  cursor: not-allowed;
}

.body--light .forgot-password {
  color: rgba(0, 0, 0, 0.7);
}

.body--light .forgot-password:hover {
  color: #000000;
}

.body--dark .login-page {
  background: #000000;
}

.body--dark .login-container {
  background: #121212;
  border: 1px solid #2d2d2d;
}

.body--dark input {
  background: #1d1d1d;
  border: 1px solid #2d2d2d;
  color: #ffffff;
}

.body--dark input::placeholder {
  color: rgba(255, 255, 255, 0.6);
}

.body--dark input:focus {
  border-color: #ffffff;
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.1);
}

.body--dark h1 {
  color: #ffffff;
}

.body--dark button {
  background: #ffffff;
  color: #000000;
  border: none;
  font-weight: 500;
}

.body--dark button:hover {
  background: #f5f5f5;
}

.body--dark button:disabled {
  background: rgba(255, 255, 255, 0.5);
  cursor: not-allowed;
}

.body--dark .forgot-password {
  color: rgba(255, 255, 255, 0.7);
}

.body--dark .forgot-password:hover {
  color: #ffffff;
}

h1 {
  text-align: center;
  margin-bottom: 2rem;
  font-size: 1.75rem;
  font-weight: 600;
  color: inherit;
}

input {
  width: 100%;
  padding: 0.9rem;
  border: 1px solid;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.error-message {
  display: block;
  margin-top: 0.5rem;
  font-size: 0.8rem;
  color: var(--q-negative);
}

input.error {
  border-color: var(--q-negative);
}

.error-message.general-error {
  background-color: var(--q-negative);
  color: white;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 15px;
  text-align: center;
  opacity: 0.9;
}

.recovery-card {
  min-width: 300px;
  max-width: 90vw;
  border-radius: 12px;
  overflow: hidden;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
}

.close-icon {
  cursor: pointer;
  font-size: 1.5rem;
  opacity: 0.8;
  transition: opacity 0.3s ease;
}

.close-icon:hover {
  opacity: 1;
}

.content-section {
  padding: 24px;
  text-align: center;
}

.success-message {
  padding: 24px;
  text-align: center;
}

/* Dark theme styles */
.body--dark .recovery-card {
  background: #121212 !important;
  color: #ffffff !important;
  border: 1px solid #2d2d2d;
}

.body--dark .header-section {
  background: #1d1d1d !important;
  border-bottom: 1px solid #2d2d2d;
}

.body--dark .text-h6,
.body--dark .text-body1,
.body--dark .text-caption {
  color: #ffffff !important;
}

.body--dark .q-field__native,
.body--dark .q-field__prefix,
.body--dark .q-field__suffix {
  color: #ffffff !important;
}

.body--dark .q-field__marginal {
  color: rgba(255, 255, 255, 0.7) !important;
}

.body--dark .success-message {
  background: #1d1d1d;
  border-top: 1px solid #2d2d2d;
  color: #4CAF50;
}

.body--dark a {
  color: #ffffff;
}

/* Light theme styles */
.body--light .recovery-card {
  background: #ffffff !important;
  color: #000000 !important;
  border: 1px solid #e0e0e0;
}

.body--light .header-section {
  background: #f5f5f5 !important;
  border-bottom: 1px solid #e0e0e0;
}

.body--light .text-h6,
.body--light .text-body1,
.body--light .text-caption {
  color: #000000 !important;
}

.body--light .q-field__native,
.body--light .q-field__prefix,
.body--light .q-field__suffix {
  color: #000000 !important;
}

.body--light .q-field__marginal {
  color: rgba(0, 0, 0, 0.7) !important;
}

.body--light .success-message {
  background: #f5f5f5;
  border-top: 1px solid #e0e0e0;
  color: #4CAF50;
}

.body--light a {
  color: #000000;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid var(--q-separator-dark);
  border-radius: 50%;
  border-top-color: var(--q-primary);
  animation: spin 1s ease-in-out infinite;
}

button:hover {
  opacity: 0.9;
}

.forgot-password:hover {
  opacity: 0.8;
}

.shake-animation {
  animation: shake 0.5s;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-5px); }
  40%, 80% { transform: translateX(5px); }
}

.logo {
  text-align: center;
  margin-bottom: 1.5rem;
}

.logo svg {
  width: 48px;
  height: 48px;
  margin-bottom: 1rem;
  color: var(--primary);
}

.input-group {
  margin-bottom: 1.5rem;
  position: relative;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--text);
}

.forgot-password {
  display: inline;
  float: right;
  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
  font-size: 0.85rem;
  color: #666;
  text-decoration: none;
}

button {
  width: 100%;
  padding: 1rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.divider {
  display: flex;
  align-items: center;
  margin: 1.5rem 0;
  color: var(--q-primary);
  opacity: 0.7;
  font-size: 0.9rem;
}

.divider::before, .divider::after {
  content: "";
  flex: 1;
  border-bottom: 1px solid var(--q-primary);
  opacity: 0.3;
}

.divider::before {
  margin-right: 1rem;
}

.divider::after {
  margin-left: 1rem;
}

.signup-link {
  text-align: center;
  margin-top: 1.5rem;
  font-size: 0.95rem;
}

.signup-link a {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
}

.illustration-container {
  display: flex;
  justify-content: center;
  margin: 20px 0;
}

.illustration-container img {
  width: 120px;
  height: 120px;
}

.text-body1 {
  margin: 16px 0;
  font-size: 1rem;
  line-height: 1.5;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.text-center .text-caption {
  color: #666;
}

.text-center a {
  color: #888 !important;
  text-decoration: none;
}

.text-center a:hover {
  color: #aaa !important;
}

@media (max-width: 600px) {
  .recovery-card {
    min-width: unset;
    width: 95vw;
    margin: 10px;
  }

  .content-section {
    padding: 16px;
  }

  .text-body1 {
    font-size: 0.95rem;
  }

  .illustration-container img {
    width: 100px;
    height: 100px;
  }
}

.body--dark .logo q-icon {
  color: #ffffff !important;
}

.body--light .logo q-icon {
  color: #000000 !important;
}
</style>
