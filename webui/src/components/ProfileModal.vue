<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-dialog shadow-lg">
      <div class="modal-content p-4">
        <div class="d-flex justify-content-between align-items-center mb-3">
          <h4 class="mb-0">Il tuo Profilo</h4>
          <button type="button" class="btn-close" @click="$emit('close')"></button>
        </div>

        <div class="text-center my-4">
          <div class="position-relative d-inline-block">
            <img 
              :src="myPhoto || 'https://placehold.co/150'" 
              class="rounded-circle border object-fit-cover" 
              width="120" height="120"
            >
            <button 
              class="btn btn-sm btn-primary position-absolute bottom-0 end-0 rounded-circle" 
              style="width: 35px; height: 35px;"
              @click="$refs.profileInput.click()"
              title="Cambia Foto"
            >
              📷
            </button>
          </div>
        </div>

        <div class="mb-4 text-center">
          <label class="text-muted small">Username</label>
          <div class="fs-4 fw-bold">{{ myUsername }}</div>
          <div class="small text-muted">ID: {{ myId ? myId.substring(0,8) + '...' : '...' }}</div>
        </div>

        <button class="btn btn-secondary w-100" @click="$emit('close')">Chiudi</button>
        
        <input type="file" ref="profileInput" class="d-none" accept="image/*" @change="handleFile">
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: ['myUsername', 'myPhoto', 'myId'],
  emits: ['close', 'update-photo'],
  methods: {
    handleFile(e) {
      if (e.target.files[0]) {
        this.$emit('update-photo', e.target.files[0]);
      }
    }
  }
}
</script>

<style scoped>
/* Stesso stile backdrop degli altri */
.modal-backdrop {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center;
  z-index: 1050;
}
.modal-dialog { background: white; border-radius: 12px; min-width: 350px; max-width: 90%; }
.object-fit-cover { object-fit: cover; }
</style>