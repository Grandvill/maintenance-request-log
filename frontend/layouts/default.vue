<script setup lang="ts">
import { Wrench, Plus, Users, LogOut, User as UserIcon, Shield } from 'lucide-vue-next'

const { user, isAdmin, logout } = useAuth()

const handleLogout = async () => {
  await logout()
}

const roleBadgeClass = computed(() => {
  switch (user.value?.role) {
    case 'admin':
      return 'bg-purple-100 text-purple-800 border-purple-200'
    case 'supervisor':
      return 'bg-blue-100 text-blue-800 border-blue-200'
    case 'operator':
    default:
      return 'bg-amber-100 text-amber-800 border-amber-200'
  }
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-slate-50">
    <!-- Navbar Header -->
    <header class="bg-white border-b border-slate-200 sticky top-0 z-30 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <!-- Left side: Brand Logo & Navigation -->
          <div class="flex items-center space-x-8">
            <NuxtLink to="/requests" class="flex items-center gap-2.5 font-bold text-slate-800 tracking-tight text-lg hover:opacity-90">
              <div class="w-9 h-9 rounded-lg bg-factory-600 flex items-center justify-center text-white shadow-sm">
                <Wrench class="w-5 h-5" />
              </div>
              <div>
                <span class="text-factory-700">HIROSE</span>
                <span class="text-slate-700 font-medium text-sm block -mt-1">Maintenance Log</span>
              </div>
            </NuxtLink>

            <nav class="hidden md:flex items-center space-x-1">
              <NuxtLink
                to="/requests"
                class="px-3 py-2 rounded-md text-sm font-medium transition-colors"
                :class="$route.path === '/requests' ? 'bg-slate-100 text-slate-900 font-semibold' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'"
              >
                All Requests
              </NuxtLink>

              <NuxtLink
                to="/requests/create"
                class="px-3 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-1.5"
                :class="$route.path === '/requests/create' ? 'bg-slate-100 text-slate-900 font-semibold' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'"
              >
                <Plus class="w-4 h-4 text-factory-600" />
                New Request
              </NuxtLink>

              <!-- Admin Only Navigation Item -->
              <NuxtLink
                v-if="isAdmin"
                to="/users"
                class="px-3 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-1.5"
                :class="$route.path === '/users' ? 'bg-slate-100 text-slate-900 font-semibold' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'"
              >
                <Users class="w-4 h-4 text-purple-600" />
                User Management
              </NuxtLink>
            </nav>
          </div>

          <!-- Right side: User Profile & Logout -->
          <div class="flex items-center space-x-4">
            <div v-if="user" class="hidden sm:flex items-center gap-3 pr-3 border-r border-slate-200">
              <div class="w-8 h-8 rounded-full bg-slate-200 flex items-center justify-center text-slate-600">
                <UserIcon class="w-4 h-4" />
              </div>
              <div class="text-right">
                <div class="text-sm font-semibold text-slate-800 leading-tight">
                  {{ user.full_name || user.username }}
                </div>
                <div class="flex items-center justify-end gap-1 mt-0.5">
                  <span
                    class="px-1.5 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider border"
                    :class="roleBadgeClass"
                  >
                    {{ user.role }}
                  </span>
                </div>
              </div>
            </div>

            <button
              @click="handleLogout"
              class="inline-flex items-center gap-1.5 text-sm font-medium text-slate-600 hover:text-rose-600 px-3 py-2 rounded-md hover:bg-rose-50 transition-colors"
              title="Logout"
            >
              <LogOut class="w-4 h-4" />
              <span class="hidden sm:inline">Logout</span>
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Page Content -->
    <main class="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <slot />
    </main>

    <!-- Footer -->
    <footer class="bg-white border-t border-slate-200 py-4 text-center text-xs text-slate-500">
      PT. Hirose Electric Indonesia — Factory Maintenance Request System
    </footer>
  </div>
</template>

