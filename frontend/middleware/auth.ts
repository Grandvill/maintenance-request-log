export default defineNuxtRouteMiddleware((to, from) => {
  const { isAuthenticated, user } = useAuth()

  // Allow navigation to login page
  if (to.path === '/login') {
    if (isAuthenticated.value) {
      return navigateTo('/requests')
    }
    return
  }

  // If user is not authenticated, redirect to login
  if (!isAuthenticated.value) {
    return navigateTo('/login')
  }

  // Protect /users route: Admin ONLY
  if (to.path.startsWith('/users')) {
    if (user.value?.role !== 'admin') {
      return navigateTo('/requests')
    }
  }
})

