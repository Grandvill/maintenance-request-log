export const useApi = () => {
  const config = useRuntimeConfig()
  const { token, logout } = useAuth()

  const apiFetch = async <T>(endpoint: string, options: any = {}): Promise<T> => {
    const url = `${config.public.apiBaseUrl}${endpoint}`

    const headers: Record<string, string> = {
      Accept: 'application/json',
      ...options.headers
    }

    if (token.value) {
      headers.Authorization = `Bearer ${token.value}`
    }

    try {
      return await $fetch<T>(url, {
        ...options,
        headers
      })
    } catch (err: any) {
      // If 401 Unauthorized, automatically log out
      if (err?.status === 401 || err?.statusCode === 401) {
        await logout()
      }
      throw err
    }
  }

  return {
    apiFetch
  }
}

