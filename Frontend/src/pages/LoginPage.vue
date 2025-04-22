<template>
  <div class="login-page">
    <div class="login-container" :class="{ 'shake-animation': showShake }">
      <div class="logo">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
          <circle cx="8.5" cy="7" r="4"></circle>
          <line x1="20" y1="8" x2="20" y2="14"></line>
          <line x1="23" y1="11" x2="17" y2="11"></line>
        </svg>
      </div>

      <h1>Welcome back</h1>

      <form @submit.prevent="handleSubmit">
        <div class="input-group">
          <label for="email">Email</label>
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
          <label for="password">Password</label>
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
          <span v-if="!isSubmitting">Sign In</span>
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
            <a href="#" class="text-primary" @click.prevent="showModal = false">Sign in</a>
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

<script setup lang="ts">
import {computed } from 'vue'
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
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

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

  if (errors.email || errors.password) {
    showShake.value = true
    setTimeout(() => showShake.value = false, 500)
    return
  }

  isSubmitting.value = true

  try {
    await new Promise(resolve => setTimeout(resolve, 1500))

    await router.push('/app') // ✅ Fixed
  } catch (error) {
    console.error('Login error:', error)
  } finally {
    isSubmitting.value = false
  }
}


</script>
<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  color: var(--text);
  line-height: 1.6;
  padding: 1rem;
}

.login-container {
  background: rgba(255, 255, 255, 0.76);
  padding: 2.5rem;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  width: 100%;
  max-width: 500px;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.login-container:hover {
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.12);
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

h1 {
  text-align: center;
  margin-bottom: 2rem;
  font-size: 1.75rem;
  font-weight: 600;
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

input {
  width: 100%;
  padding: 0.9rem;
  border: 1px solid var(--gray);
  border-radius: 8px;
  font-size: 1rem;
  transition: border 0.3s ease;
}

input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(67, 97, 238, 0.2);
}

input.error {
  border-color: var(--error);
}

.error-message {
  display: block;
  margin-top: 0.5rem;
  font-size: 0.8rem;
  color: var(--error);
}

.forgot-password {
  display: block;
  text-align: right;
  margin-top: 0.5rem;
  font-size: 0.85rem;
  color: var(--primary);
  text-decoration: none;
}

button {
  width: 100%;
  padding: 1rem;
  background-color: var(--primary);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

button:hover {
  background-color: var(--primary-dark);
}

button:disabled {
  background-color: var(--dark-gray);
  cursor: not-allowed;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.social-login {
  background: white;
  color: var(--text);
  border: 1px solid var(--gray);
}

.social-login:hover {
  background: var(--light-gray);
}

.divider {
  display: flex;
  align-items: center;
  margin: 1.5rem 0;
  color: var(--dark-gray);
  font-size: 0.9rem;
}

.divider::before, .divider::after {
  content: "";
  flex: 1;
  border-bottom: 1px solid var(--gray);
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

.recovery-card {
  min-width: 400px;
  max-width: 90vw;
  border-radius: 12px;
  overflow: hidden;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.close-icon {
  cursor: pointer;
  font-size: 1.5rem;
  color: #6c757d;
}

.close-icon:hover {
  color: #495057;
}

.content-section {
  padding: 24px;
}

.illustration-container {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.success-message {
  background-color: #f8f9fa;
  text-align: center;
  padding: 24px;
  border-top: 1px solid #e9ecef;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 600px) {
  .recovery-card {
    min-width: 90vw;
  }

  .content-section {
    padding: 16px;
  }

  .login-container {
    padding: 1.5rem;
  }
}
</style>
