<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-dialog shadow-lg">
      <div class="modal-content p-4">
        <h4 class="mb-3">Nuova Conversazione</h4>
        
        <ul class="nav nav-pills nav-fill mb-4 bg-light p-1 rounded">
          <li class="nav-item">
            <a class="nav-link" :class="{active: !isGroup}" href="#" style="cursor:pointer" @click="isGroup=false">Privata</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" :class="{active: isGroup}" href="#" style="cursor:pointer" @click="isGroup=true">Gruppo</a>
          </li>
        </ul>

        <div v-if="!isGroup">
          <div class="mb-3">
            <label class="form-label">Nome Utente</label>
            <input v-model="targetName" type="text" class="form-control" placeholder="Es. Bob" @keyup.enter="confirm">
            <div class="form-text">Cerca un utente per nome.</div>
          </div>
        </div>

        <div v-else>
          <div class="mb-3">
            <label class="form-label">Nome del Gruppo</label>
            <input v-model="groupName" type="text" class="form-control" placeholder="Es. Calcetto" @keyup.enter="confirm">
          </div>
          <div class="mb-3">
            <label class="form-label">Aggiungi un Membro (Username)</label>
            <input v-model="targetName" type="text" class="form-control" placeholder="Es. Alice">
          </div>
        </div>

        <div class="d-flex gap-2 mt-4">
          <button class="btn btn-light flex-grow-1" @click="$emit('close')">Annulla</button>
          <button class="btn btn-primary flex-grow-1" :disabled="!isValid" @click="confirm">
            {{ isGroup ? 'Crea Gruppo' : 'Inizia Chat' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  emits: ['close', 'create-private', 'create-group'],
  data() {
    return {
      isGroup: false,
      targetName: '',
      groupName: ''
    }
  },
  computed: {
    isValid() {
      if (!this.isGroup) return this.targetName.length > 0;
      return this.groupName.length > 0 && this.targetName.length > 0;
    }
  },
  methods: {
    confirm() {
      if (!this.isValid) return;
      
      if (this.isGroup) {
        this.$emit('create-group', { groupName: this.groupName, memberName: this.targetName });
      } else {
        this.$emit('create-private', this.targetName);
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
  z-index: 1050;
}
.modal-dialog { background: white; border-radius: 12px; min-width: 350px; max-width: 90%; }
</style>