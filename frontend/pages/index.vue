<script setup lang="ts">
definePageMeta({
  middleware: [
    function () {
      const { isAuthenticated } = useAuth()
      if (isAuthenticated.value) {
        return navigateTo('/requests')
      }
      return navigateTo('/login')
    }
  ]
})
</script>

<template>
  <div class="min-h-[50vh] flex items-center justify-center">
    <div class="text-center text-slate-500 text-sm">
      Redirecting to portal...
    </div>
  </div>
</template>

