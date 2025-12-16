<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-dialog shadow-lg">
      <div class="modal-content p-4">

        <div class="d-flex justify-content-between align-items-center mb-3 border-bottom pb-2">
          <h5 class="mb-0">Gestione Gruppo</h5>
          <button type="button" class="btn-close" @click="$emit('close')"></button>
        </div>

        <div class="mb-4">
          <label class="form-label text-muted small fw-bold">NOME GRUPPO</label>
          <div class="input-group">
            <input type="text" class="form-control" v-model="editName" placeholder="Nome gruppo">
            <button class="btn btn-outline-primary" @click="updateName" :disabled="editName === currentName">
              Salva
            </button>
          </div>
        </div>

        <div class="mb-4">
          <label class="form-label text-muted small fw-bold">AGGIUNGI MEMBRO</label>
          <div class="input-group">
            <input type="text" class="form-control" v-model="newMemberName" placeholder="Username (es. Bob)" @keyup.enter="addMember">
            <button class="btn btn-success" @click="addMember">
              +
            </button>
          </div>
        </div>

        <div>
          <label class="form-label text-muted small fw-bold">MEMBRI ({{ members.length }})</label>
          <div class="list-group overflow-auto" style="max-height: 200px;">
            <div v-if="loading" class="text-center p-2">Caricamento...</div>

            <div v-else v-for="u in members" :key="u.id" class="list-group-item d-flex align-items-center">
              <img :src="u.photo_url || 'https://placehold.co/40'" class="rounded-circle me-2 object-fit-cover" width="32" height="32">
              <span class="text-truncate">{{ u.name }}</span>
              </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script>
import axios from '../services/axios';

export default {
  props: ['groupId', 'currentName'],
  emits: ['close', 'name-updated'],
  data() {
    return {
      members: [],
      editName: this.currentName,
      newMemberName: '',
      loading: false
    }
  },
  mounted() {
    this.loadMembers();
  },
  methods: {
    // 1. Carica la lista membri
    async loadMembers() {
      this.loading = true;
      try {
        let res = await axios.get(`/groups/${this.groupId}/members`);
        this.members = res.data.members || [];
      } catch (e) {
        alert("Errore caricamento membri");
      }
      this.loading = false;
    },

    // 2. Aggiorna il nome del gruppo
    async updateName() {
      if (!this.editName.trim()) return;
      try {
        await axios.put(`/groups/${this.groupId}/name`, { name: this.editName });
        this.$emit('name-updated'); // Avvisa il padre di ricaricare
        alert("Nome aggiornato!");
      } catch (e) {
        alert("Errore aggiornamento nome");
      }
    },

    // 3. Aggiungi Membro (Cerca ID dal nome -> Aggiungi ID)
    async addMember() {
      if (!this.newMemberName.trim()) return;

      try {
        // A. Trova l'ID dell'utente (usando il trucco del login)
        let userRes = await axios.post('/session', { name: this.newMemberName });
        let newUserId = userRes.data.identifier;

        // B. Aggiungi al gruppo
        await axios.post(`/groups/${this.groupId}/members`, { user_id: newUserId });

        this.newMemberName = '';
        this.loadMembers(); // Ricarica la lista
      } catch (e) {
        alert("Impossibile aggiungere utente (non trovato o errore server)");
      }
    }
  }
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center;
  z-index: 1060; /* Più alto degli altri modali */
}
.modal-dialog { background: white; border-radius: 12px; min-width: 400px; max-width: 95%; }
.object-fit-cover { object-fit: cover; }
</style>
