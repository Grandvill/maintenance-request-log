<script setup lang="ts">
import { Wrench, Shield, AlertCircle, Loader2 } from 'lucide-vue-next'

definePageMeta({
  layout: false // Standalone layout for login
})

const { login, isAuthenticated } = useAuth()
const router = useRouter()

// If already logged in, redirect to requests
onMounted(() => {
  if (isAuthenticated.value) {
    router.replace('/requests')
  }
})

const form = reactive({
  username: '',
  password: ''
})

const isLoading = ref(false)
const errorMessage = ref('')

const handleSubmit = async () => {
  if (!form.username || !form.password) {
    errorMessage.value = 'Please enter both username and password.'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  const result = await login(form.username, form.password)
  isLoading.value = false

  if (result.success) {
    router.push('/requests')
  } else {
    errorMessage.value = result.error || 'Invalid credentials'
  }
}

// Quick fill for testing
const fillCredentials = (user: string, pass: string) => {
  form.username = user
  form.password = pass
  errorMessage.value = ''
}
</script>

<template>
  <div class="min-h-screen bg-slate-100 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md text-center">
      <div class="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-factory-600 text-white shadow-lg mb-4">
        <Wrench class="w-8 h-8" />
      </div>
      <h2 class="text-2xl font-extrabold text-slate-900 tracking-tight">
        PT. HIROSE ELECTRIC INDONESIA
      </h2>
      <p class="mt-1 text-sm text-slate-600">
        Factory Maintenance Request Log Portal
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="bg-white py-8 px-6 shadow-md rounded-xl sm:px-10 border border-slate-200">
        <!-- Error Alert -->
        <div
          v-if="errorMessage"
          class="mb-6 p-4 rounded-lg bg-rose-50 border border-rose-200 flex items-start gap-3"
        >
          <AlertCircle class="w-5 h-5 text-rose-600 shrink-0 mt-0.5" />
          <div class="text-sm text-rose-700 font-medium leading-tight">
            {{ errorMessage }}
          </div>
        </div>

        <form class="space-y-5" @submit.prevent="handleSubmit">
          <div>
            <label for="username" class="block text-sm font-medium text-slate-700">
              Username
            </label>
            <div class="mt-1">
              <input
                id="username"
                v-model="form.username"
                name="username"
                type="text"
                required
                autocomplete="username"
                placeholder="e.g. operator1, supervisor, admin"
                class="appearance-none block w-full px-3 py-2.5 border border-slate-300 rounded-lg shadow-sm placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-factory-500 focus:border-factory-500 sm:text-sm"
              />
            </div>
          </div>

          <div>
            <label for="password" class="block text-sm font-medium text-slate-700">
              Password
            </label>
            <div class="mt-1">
              <input
                id="password"
                v-model="form.password"
                name="password"
                type="password"
                required
                autocomplete="current-password"
                placeholder="••••••••"
                class="appearance-none block w-full px-3 py-2.5 border border-slate-300 rounded-lg shadow-sm placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-factory-500 focus:border-factory-500 sm:text-sm"
              />
            </div>
          </div>

          <div>
            <button
              type="submit"
              :disabled="isLoading"
              class="w-full flex justify-center py-2.5 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-factory-600 hover:bg-factory-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-factory-500 disabled:opacity-50 transition-colors"
            >
              <Loader2 v-if="isLoading" class="w-5 h-5 animate-spin mr-2" />
              <span>{{ isLoading ? 'Signing in...' : 'Sign In' }}</span>
            </button>
          </div>
        </form>

        <!-- Reviewer Quick Login Helpers -->
        <div class="mt-8 pt-6 border-t border-slate-200">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 text-center">
            Quick Fill Seed Credentials (for testing)
          </p>
          <div class="grid grid-cols-3 gap-2">
            <button
              type="button"
              @click="fillCredentials('operator1', 'operator123')"
              class="px-2 py-1.5 border border-amber-200 bg-amber-50 hover:bg-amber-100 rounded text-xs font-medium text-amber-800 transition-colors text-center"
            >
              Operator 1
            </button>
            <button
              type="button"
              @click="fillCredentials('supervisor', 'supervisor123')"
              class="px-2 py-1.5 border border-blue-200 bg-blue-50 hover:bg-blue-100 rounded text-xs font-medium text-blue-800 transition-colors text-center"
            >
              Supervisor
            </button>
            <button
              type="button"
              @click="fillCredentials('admin', 'admin123')"
              class="px-2 py-1.5 border border-purple-200 bg-purple-50 hover:bg-purple-100 rounded text-xs font-medium text-purple-800 transition-colors text-center"
            >
              Admin
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

