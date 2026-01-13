<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-dialog shadow-lg">
      <div class="modal-content p-4">
        <div class="d-flex justify-content-between align-items-center mb-3">
          <h4 class="mb-0">Il tuo Profilo</h4>
          <button type="button" class="btn-close" @click="$emit('close')" />
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
              title="Cambia Foto"
              @click="$refs.profileInput.click()"
            >
              📷
            </button>
          </div>
        </div>

        <div class="mb-4 text-center">
          <label class="text-muted small">Username</label>
          <div class="fs-4 fw-bold">{{ myUsername }}</div>
          <div class="small text-muted">ID: {{ myId ? String(myId).substring(0,8) + '...' : '...' }}</div>
        </div>
		  <div class="mb-3">
			  <label for="username-input" class="form-label">Username</label>
			  <div class="input-group">
				  <input
					  id="username-input"
					  type="text"
					  class="form-control"
					  v-model="newUsername"
					  placeholder="Nuovo nome..."
				  >
				  <button
					  class="btn btn-primary"
					  type="button"
					  @click="updateUsername"
					  :disabled="newUsername.length < 3 || newUsername === myUsername"
				  >
					  Aggiorna
				  </button>
			  </div>
			  <div class="form-text">Minimo 3 caratteri.</div>
		  </div>

        <button class="btn btn-secondary w-100" @click="$emit('close')">Chiudi</button>

        <input ref="profileInput" type="file" class="d-none" accept="image/*" @change="handleFile">
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: ['myUsername', 'myPhoto', 'myId'],
  emits: ['close', 'update-photo','update:myUsername'],

	data() {
		return {
			// Inizializziamo il campo con il nome attuale
			newUsername: this.myUsername,
			loading: false,
			errormsg: null,
		}
	},

	// Importante: se il nome cambia da fuori, aggiorniamo anche l'input
	watch: {
		myUsername(newVal) {
			this.newUsername = newVal;
		}
	},
  methods: {
	  async updateUsername() {
		  this.loading = true;
		  this.errormsg = null;

		  try {
			  // Logghiamo cosa stiamo facendo (senza usare variabili che non esistono)
			  console.log("Invio richiesta a: /me/name");
			  console.log("Nuovo nome:", this.newUsername);

			  // CHIAMATA CORRETTA
			  let response = await this.$axios.put('/me/name', {
				  name: this.newUsername
			  });

			  // Se arriviamo qui, è andato tutto bene
			  console.log("Successo!");

			  // Aggiorniamo il frontend
			  this.$emit('update:myUsername', this.newUsername);
			  alert("Nome aggiornato con successo!");

			  // Chiudiamo il modale (opzionale, se vuoi che si chiuda da solo)
			  // this.$emit('close');

		  } catch (e) {
			  // Gestione errori specifica
			  if (e.response && e.response.status === 409) {
				  // 409 Conflict -> Messaggio dal backend: "Nome già utilizzato..."
				  this.errorMsg = "Spiacente, questo nome è già in uso. Scegline un altro.";
			  } else if (e.response && e.response.status === 400) {
				  this.errorMsg = "Nome non valido.";
			  } else {
				  this.errorMsg = "Errore del server. Riprova più tardi.";
			  }
		  } finally {
			  this.loading = false;
		  }
	  },
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
