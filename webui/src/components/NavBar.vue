<template>
  <nav class="navbar navbar-light bg-white border-bottom px-4 justify-content-between">
    <div class="d-flex align-items-center">
      <span class="navbar-brand fw-bold text-primary mb-0 h1">
        <i class="bi bi-chat-dots-fill"></i> WASAText
      </span>
    </div>
    
    <div class="d-flex align-items-center">
      <div 
        class="d-flex align-items-center me-3 p-1 rounded hover-bg" 
        style="cursor: pointer;"
        @click="$emit('open-profile')" 
        title="Modifica Profilo"
      >
        <img 
          :src="myPhoto || 'https://placehold.co/50'" 
          class="rounded-circle border me-2 object-fit-cover" 
          width="32" height="32"
        >
        <span class="fw-bold text-dark d-none d-md-block">{{ myUsername }}</span>
      </div>

      <button class="btn btn-outline-danger btn-sm" @click="$emit('logout')">
        Esci
      </button>
    </div>
  </nav>
</template>

<script>
export default {
  props: ['myUsername', 'myPhoto'],
  emits: ['logout', 'open-profile']
}
</script>

<style scoped>
.object-fit-cover { object-fit: cover; }
.hover-bg:hover { background-color: #f8f9fa; }
</style>