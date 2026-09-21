<script setup lang="ts">
import { 
  Plus, Users, Shield, UserCheck, UserX, 
  Edit, Lock, AlertCircle, Loader2 
} from 'lucide-vue-next'
import type { User, Role, CreateUserPayload, UpdateUserPayload } from '~/types'

definePageMeta({
  middleware: 'auth'
})

const { user: currentUser } = useAuth()
const { apiFetch } = useApi()

const users = ref<User[]>([])
const isLoading = ref(true)
const errorMessage = ref('')

// Create / Edit User Modal state
const modal = reactive({
  isOpen: false,
  isEdit: false,
  userId: '',
  username: '',
  fullName: '',
  role: 'operator' as Role,
  password: '',
  isSubmitting: false,
  error: ''
})

const fetchUsers = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    users.value = await apiFetch<User[]>('/users')
  } catch (err: any) {
    errorMessage.value = err?.data?.error || 'Failed to fetch users'
  } finally {
    isLoading.value = false
  }
}

const openCreateModal = () => {
  modal.isEdit = false
  modal.userId = ''
  modal.username = ''
  modal.fullName = ''
  modal.role = 'operator'
  modal.password = ''
  modal.error = ''
  modal.isOpen = true
}

const openEditModal = (u: User) => {
  modal.isEdit = true
  modal.userId = u.id
  modal.username = u.username
  modal.fullName = u.full_name
  modal.role = u.role
  modal.password = ''
  modal.error = ''
  modal.isOpen = true
}

const handleSaveUser = async () => {
  if (!modal.fullName.trim()) {
    modal.error = 'Full name is required.'
    return
  }

  if (!modal.isEdit) {
    if (!modal.username.trim()) {
      modal.error = 'Username is required.'
      return
    }
    if (!modal.password || modal.password.length < 6) {
      modal.error = 'Password must be at least 6 characters.'
      return
    }
  }

  modal.isSubmitting = true
  modal.error = ''

  try {
    if (modal.isEdit) {
      const payload: UpdateUserPayload = {
        full_name: modal.fullName.trim(),
        role: modal.role
      }
      if (modal.password && modal.password.trim() !== '') {
        payload.password = modal.password.trim()
      }
      await apiFetch<User>(`/users/${modal.userId}`, {
        method: 'PUT',
        body: payload
      })
    } else {
      const payload: CreateUserPayload = {
        username: modal.username.trim(),
        full_name: modal.fullName.trim(),
        role: modal.role,
        password: modal.password
      }
      await apiFetch<User>('/users', {
        method: 'POST',
        body: payload
      })
    }

    modal.isOpen = false
    await fetchUsers()
  } catch (err: any) {
    modal.error = err?.data?.error || 'Failed to save user'
  } finally {
    modal.isSubmitting = false
  }
}

// Toggle active/inactive status
const handleToggleStatus = async (u: User) => {
  if (u.id === currentUser.value?.id) {
    alert('You cannot deactivate your own admin account.')
    return
  }

  const actionText = u.is_active ? 'deactivate' : 'activate'
  if (!confirm(`Are you sure you want to ${actionText} user "${u.username}"?`)) {
    return
  }

  try {
    await apiFetch<User>(`/users/${u.id}/status`, {
      method: 'PATCH',
      body: { is_active: !u.is_active }
    })
    await fetchUsers()
  } catch (err: any) {
    alert(err?.data?.error || 'Failed to update user status')
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div>
    <!-- Header -->
    <div class="sm:flex sm:items-center sm:justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2.5">
          <Users class="w-6 h-6 text-purple-600" />
          User Management
        </h1>
        <p class="mt-1 text-sm text-slate-500">
          Create, edit, and deactivate factory personnel accounts (Admin exclusive).
        </p>
      </div>

      <div class="mt-4 sm:mt-0">
        <button
          @click="openCreateModal"
          class="inline-flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white text-sm font-medium rounded-lg shadow-sm transition-colors"
        >
          <Plus class="w-4 h-4" />
          Create New User
        </button>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="errorMessage" class="mb-6 p-4 rounded-lg bg-rose-50 border border-rose-200 text-rose-700 text-sm">
      {{ errorMessage }}
    </div>

    <!-- Users Table -->
    <div class="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
      <div v-if="isLoading" class="py-16 text-center text-slate-500">
        <Loader2 class="w-8 h-8 animate-spin mx-auto mb-2 text-purple-600" />
        <p class="text-sm">Loading users list...</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-200 text-left text-sm">
          <thead class="bg-slate-50 text-slate-600 font-semibold uppercase text-[11px] tracking-wider">
            <tr>
              <th scope="col" class="py-3.5 pl-4 pr-3 sm:pl-6">Username</th>
              <th scope="col" class="px-3 py-3.5">Full Name</th>
              <th scope="col" class="px-3 py-3.5">Role</th>
              <th scope="col" class="px-3 py-3.5">Status</th>
              <th scope="col" class="px-3 py-3.5">Created At</th>
              <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 bg-white">
            <tr v-for="u in users" :key="u.id" class="hover:bg-slate-50/80 transition-colors">
              <!-- Username -->
              <td class="whitespace-nowrap py-4 pl-4 pr-3 sm:pl-6 font-bold text-slate-900 font-mono">
                {{ u.username }}
                <span v-if="u.id === currentUser?.id" class="ml-1 text-xs text-purple-600 font-sans font-normal">(You)</span>
              </td>

              <!-- Full Name -->
              <td class="whitespace-nowrap px-3 py-4 text-slate-700">
                {{ u.full_name }}
              </td>

              <!-- Role -->
              <td class="whitespace-nowrap px-3 py-4">
                <span
                  class="px-2 py-0.5 rounded text-xs font-semibold uppercase tracking-wider border"
                  :class="{
                    'bg-purple-100 text-purple-800 border-purple-200': u.role === 'admin',
                    'bg-blue-100 text-blue-800 border-blue-200': u.role === 'supervisor',
                    'bg-amber-100 text-amber-800 border-amber-200': u.role === 'operator'
                  }"
                >
                  {{ u.role }}
                </span>
              </td>

              <!-- Status -->
              <td class="whitespace-nowrap px-3 py-4">
                <span
                  class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                  :class="u.is_active ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-slate-100 text-slate-500 border-slate-200'"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="u.is_active ? 'bg-emerald-500' : 'bg-slate-400'"></span>
                  {{ u.is_active ? 'Active' : 'Deactivated' }}
                </span>
              </td>

              <!-- Created At -->
              <td class="whitespace-nowrap px-3 py-4 text-slate-500 text-xs">
                {{ formatDate(u.created_at) }}
              </td>

              <!-- Actions -->
              <td class="whitespace-nowrap py-4 pl-3 pr-4 sm:pr-6 text-right space-x-2">
                <!-- Edit button -->
                <button
                  @click="openEditModal(u)"
                  class="inline-flex items-center p-1.5 text-blue-600 hover:text-blue-800 rounded hover:bg-blue-50"
                  title="Edit User"
                >
                  <Edit class="w-4 h-4" />
                </button>

                <!-- Deactivate/Activate button -->
                <button
                  v-if="u.id !== currentUser?.id"
                  @click="handleToggleStatus(u)"
                  class="inline-flex items-center p-1.5 rounded"
                  :class="u.is_active ? 'text-rose-600 hover:text-rose-800 hover:bg-rose-50' : 'text-emerald-600 hover:text-emerald-800 hover:bg-emerald-50'"
                  :title="u.is_active ? 'Deactivate User' : 'Activate User'"
                >
                  <UserX v-if="u.is_active" class="w-4 h-4" />
                  <UserCheck v-else class="w-4 h-4" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create / Edit User Modal -->
    <div
      v-if="modal.isOpen"
      class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/50 flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-xl border border-slate-200">
        <h3 class="text-lg font-bold text-slate-900 mb-4">
          {{ modal.isEdit ? `Edit User: ${modal.username}` : 'Create New Personnel Account' }}
        </h3>

        <!-- Error Banner -->
        <div v-if="modal.error" class="mb-4 p-3 rounded-lg bg-rose-50 border border-rose-200 text-xs text-rose-700 font-medium">
          {{ modal.error }}
        </div>

        <form class="space-y-4" @submit.prevent="handleSaveUser">
          <!-- Username (only editable when creating) -->
          <div v-if="!modal.isEdit">
            <label class="block text-xs font-semibold text-slate-700 uppercase mb-1">
              Username <span class="text-rose-500">*</span>
            </label>
            <input
              v-model="modal.username"
              type="text"
              required
              placeholder="e.g. operator3, supervisor2"
              class="w-full px-3 py-2 text-sm border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
            />
          </div>

          <!-- Full Name -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase mb-1">
              Full Name <span class="text-rose-500">*</span>
            </label>
            <input
              v-model="modal.fullName"
              type="text"
              required
              placeholder="e.g. Ahmad Fauzi"
              class="w-full px-3 py-2 text-sm border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
            />
          </div>

          <!-- Role -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase mb-1">
              System Role <span class="text-rose-500">*</span>
            </label>
            <select
              v-model="modal.role"
              class="w-full px-3 py-2 text-sm border border-slate-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-purple-500"
            >
              <option value="operator">Operator (Raise & view own requests)</option>
              <option value="supervisor">Supervisor (View all & Approve/Reject)</option>
              <option value="admin">Admin (Full Access & User Management)</option>
            </select>
          </div>

          <!-- Password -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase mb-1">
              Password <span v-if="!modal.isEdit" class="text-rose-500">*</span>
              <span v-else class="text-slate-400 font-normal lowercase">(leave blank to keep current)</span>
            </label>
            <input
              v-model="modal.password"
              type="password"
              :required="!modal.isEdit"
              placeholder="••••••••"
              class="w-full px-3 py-2 text-sm border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
            />
          </div>

          <div class="flex justify-end gap-2 pt-4 border-t border-slate-200">
            <button
              type="button"
              @click="modal.isOpen = false"
              class="px-4 py-2 text-sm font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-lg"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="modal.isSubmitting"
              class="px-4 py-2 text-sm font-medium text-white bg-purple-600 hover:bg-purple-700 rounded-lg shadow-sm disabled:opacity-50"
            >
              {{ modal.isSubmitting ? 'Saving...' : 'Save User' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

