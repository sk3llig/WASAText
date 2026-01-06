<template>
  <div class="col-md-4 col-lg-3 border-end d-flex flex-column bg-light h-100">
    <div class="p-3 border-bottom bg-white">
      <button class="btn btn-primary w-100 shadow-sm" @click="$emit('open-new-chat')">
        + Nuova Conversazione
      </button>
    </div>

    <div class="list-group list-group-flush flex-grow-1 overflow-auto">
      <div 
        v-for="conv in conversations" 
        :key="conv.id"
        class="list-group-item list-group-item-action p-3 border-bottom position-relative"
        :class="{ 'active-chat': activeConvId === conv.id }"
        style="cursor: pointer;"
        @click="$emit('select-chat', conv)"
      >
        <div class="d-flex align-items-center">
          <img :src="conv.photo_url || 'https://placehold.co/50'" class="rounded-circle me-3 object-fit-cover border" width="50" height="50">
          <div class="flex-grow-1 overflow-hidden" style="margin-right: 25px;">
            <div class="d-flex justify-content-between align-items-center">
              <h6 class="mb-0 text-truncate">{{ conv.name || 'Chat' }}</h6>
              <small class="text-muted" style="font-size: 0.7rem;">{{ formatTime(conv.last_message_time) }}</small>
            </div>
            <small class="text-muted d-block text-truncate">
              {{ conv.last_message || '...' }}
            </small>
          </div>
        </div>
        
        <button 
          class="btn btn-link text-danger position-absolute end-0 top-50 translate-middle-y p-2 delete-btn" 
          title="Elimina" 
          @click.stop="$emit('delete-chat', conv)"
        >
          🗑️
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  // I dati che arrivano dal padre (HomeView)
  props: ['conversations', 'activeConvId'],
  // Gli eventi che questo componente può lanciare verso il padre
  emits: ['open-new-chat', 'select-chat', 'delete-chat'],
  methods: {
    formatTime(iso) { 
        if(!iso) return ''; 
        return new Date(iso).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}); 
    }
  }
}
</script>

<style scoped>
.active-chat { background-color: #e9ecef; border-left: 4px solid #0d6efd; }
.delete-btn { display: none; }
.list-group-item:hover .delete-btn { display: block; }
.object-fit-cover { object-fit: cover; }
</style>