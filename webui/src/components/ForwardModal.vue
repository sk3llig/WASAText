<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-dialog shadow-lg">
      <div class="modal-content p-4">
        <div class="d-flex justify-content-between align-items-center mb-3 border-bottom pb-2">
          <h5 class="mb-0">Inoltra a...</h5>
          <button type="button" class="btn-close" @click="$emit('close')" />
        </div>

        <div class="list-group overflow-auto mb-3" style="max-height: 300px;">
          <div v-if="loading" class="text-center p-3">Caricamento...</div>

          <button
            v-for="conv in conversations"
            v-else
            :key="conv.id"
            class="list-group-item list-group-item-action d-flex align-items-center"
            @click="selectChat(conv)"
          >
            <img :src="conv.photo_url || 'https://placehold.co/40'" class="rounded-circle me-3 object-fit-cover" width="40" height="40">
            <div>
              <div class="fw-bold">{{ conv.name || 'Chat' }}</div>
              <small v-if="conv.is_group" class="text-muted">Gruppo</small>
            </div>
          </button>
        </div>

        <div class="mb-3">
          <label class="form-label small text-muted">Aggiungi un commento (opzionale)</label>
          <input v-model="comment" type="text" class="form-control" placeholder="Scrivi qui...">
        </div>

        <button class="btn btn-secondary w-100" @click="$emit('close')">Annulla</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from '../services/axios';

export default {
  emits: ['close', 'forward-to'],
  data() {
    return {
      conversations: [],
      loading: false,
      comment: ''
    }
  },
  mounted() {
    this.loadConversations();
  },
  methods: {
    async loadConversations() {
      this.loading = true;
      try {
        let res = await axios.get('/me/conversations');
        this.conversations = res.data.conversations || [];
      } catch (e) { console.error(e); }
      this.loading = false;
    },
    selectChat(conv) {
      // Emette l'evento al padre con l'ID della chat scelta e il commento
      this.$emit('forward-to', { targetId: conv.id, comment: this.comment });
    }
  }
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center;
  z-index: 1070; /* Alto z-index */
}
.modal-dialog { background: white; border-radius: 12px; min-width: 350px; max-width: 90%; }
.object-fit-cover { object-fit: cover; }
</style>
