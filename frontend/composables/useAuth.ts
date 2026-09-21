import type { User, Role, AuthResponse } from '~/types'

export const useAuth = () => {
  const config = useRuntimeConfig()
  const router = useRouter()

  // Persistent state using Nuxt useState
  const user = useState<User | null>('auth_user', () => {
    if (import.meta.client) {
      const stored = localStorage.getItem('maint_user')
      if (stored) {
        try {
          return JSON.parse(stored)
        } catch {
          return null
        }
      }
    }
    return null
  })

  const token = useState<string | null>('auth_token', () => {
    if (import.meta.client) {
      return localStorage.getItem('maint_token')
    }
    return null
  })

  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const role = computed<Role | undefined>(() => user.value?.role)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isSupervisor = computed(() => user.value?.role === 'supervisor')
  const isOperator = computed(() => user.value?.role === 'operator')
  const canReview = computed(() => user.value?.role === 'supervisor' || user.value?.role === 'admin')

  const hasRole = (...roles: Role[]) => {
    if (!user.value) return false
    return roles.includes(user.value.role)
  }

  const setAuth = (authData: AuthResponse) => {
    token.value = authData.token
    user.value = authData.user
    if (import.meta.client) {
      localStorage.setItem('maint_token', authData.token)
      localStorage.setItem('maint_user', JSON.stringify(authData.user))
    }
  }

  const clearAuth = () => {
    token.value = null
    user.value = null
    if (import.meta.client) {
      localStorage.removeItem('maint_token')
      localStorage.removeItem('maint_user')
    }
  }

  const login = async (username: string, password: string): Promise<{ success: boolean; error?: string }> => {
    try {
      const res = await $fetch<AuthResponse>(`${config.public.apiBaseUrl}/auth/login`, {
        method: 'POST',
        body: { username, password }
      })
      setAuth(res)
      return { success: true }
    } catch (err: any) {
      const errMsg = err?.data?.error || err?.message || 'Login failed. Please check credentials.'
      return { success: false, error: errMsg }
    }
  }

  const logout = async () => {
    try {
      if (token.value) {
        await $fetch(`${config.public.apiBaseUrl}/auth/logout`, {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token.value}`
          }
        }).catch(() => {})
      }
    } finally {
      clearAuth()
      router.push('/login')
    }
  }

  const fetchMe = async () => {
    if (!token.value) return null
    try {
      const currentUser = await $fetch<User>(`${config.public.apiBaseUrl}/auth/me`, {
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })
      user.value = currentUser
      if (import.meta.client) {
        localStorage.setItem('maint_user', JSON.stringify(currentUser))
      }
      return currentUser
    } catch (err) {
      clearAuth()
      return null
    }
  }

  return {
    user,
    token,
    isAuthenticated,
    role,
    isAdmin,
    isSupervisor,
    isOperator,
    canReview,
    hasRole,
    login,
    logout,
    fetchMe,
    setAuth,
    clearAuth
  }
}

