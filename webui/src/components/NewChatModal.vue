<template>
	<div class="modal-backdrop" @click.self="$emit('close')">
		<div class="modal-dialog shadow-lg">
			<div class="modal-content p-4">
				<h4 class="mb-3">Nuova Conversazione</h4>

				<ul class="nav nav-pills nav-fill mb-4 bg-light p-1 rounded">
					<li class="nav-item">
						<a class="nav-link" :class="{active: !isGroup}" href="#" @click="setMode(false)">Privata</a>
					</li>
					<li class="nav-item">
						<a class="nav-link" :class="{active: isGroup}" href="#" @click="setMode(true)">Gruppo</a>
					</li>
				</ul>

				<div v-if="errormsg" class="alert alert-danger py-2 small">
					{{ errormsg }}
				</div>

				<div v-if="!isGroup">
					<div class="mb-3">
						<label class="form-label">Username (o ID Utente)</label>
						<input
							v-model="targetName"
							type="text"
							class="form-control"
							placeholder="Inserisci username..."
							@keyup.enter="confirm"
							@input="errormsg = null"
						>
					</div>
				</div>

				<div v-else>
					<div class="mb-3">
						<label class="form-label">Nome del Gruppo</label>
						<input
							v-model="groupName"
							type="text"
							class="form-control"
							placeholder="Es. Calcetto"
							@input="errormsg = null"
						>
					</div>

					<div class="mb-3">
						<label class="form-label">Aggiungi Partecipanti</label>
						<div class="input-group">
							<input
								v-model="memberInput"
								type="text"
								class="form-control"
								placeholder="Username membro..."
								@keyup.enter="addMemberToList"
								@input="errormsg = null"
							>
							<button
								class="btn btn-outline-secondary"
								type="button"
								@click="addMemberToList"
								:disabled="!memberInput"
							>
								Aggiungi
							</button>
						</div>
						<div class="form-text">
							Aggiungi tutti gli utenti che vuoi, poi clicca su Crea Gruppo.
						</div>
					</div>

					<div class="mb-3" v-if="groupMembers.length > 0">
						<label class="form-label small text-muted">Partecipanti da invitare:</label>
						<div class="d-flex flex-wrap gap-2">
                        <span
							v-for="(member, index) in groupMembers"
							:key="index"
							class="badge bg-primary d-flex align-items-center gap-2 p-2"
						>
                            {{ member }}
                            <span
								@click="removeMemberFromList(index)"
								style="cursor:pointer; font-weight:bold"
								title="Rimuovi"
							>&times;</span>
                        </span>
						</div>
					</div>
				</div>

				<div class="d-flex gap-2 mt-4">
					<button class="btn btn-light flex-grow-1" @click="$emit('close')" :disabled="loading">Annulla</button>

					<button
						class="btn btn-primary flex-grow-1"
						:disabled="!isValid || loading"
						@click="confirm"
					>
						<span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
						{{ isGroup ? 'Crea Gruppo' : 'Inizia Chat' }}
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	emits: ['close', 'chat-created'],

	data() {
		return {
			isGroup: false,
			targetName: '',
			groupName: '',
			memberInput: '',
			groupMembers: [],

			loading: false,       // Caricamento del tasto finale "Crea"
			addingMember: false,  // <--- NUOVO: Spinner per il tasto "Aggiungi"
			errormsg: null
		}
	},

	computed: {
		isValid() {
			if (!this.isGroup) {
				return this.targetName.length > 0;
			} else {
				// Valido se ho un nome gruppo e almeno 1 membro in lista
				return this.groupName.length > 0 && this.groupMembers.length > 0;
			}
		}
	},

	methods: {
		setMode(groupMode) {
			this.isGroup = groupMode;
			this.errormsg = null;
			this.targetName = '';
			this.memberInput = '';
			this.groupMembers = [];
		},

		// --- QUESTA È LA FUNZIONE CHE FA IL CONTROLLO ---
		async addMemberToList() {
			const val = this.memberInput.trim();
			if (!val) return;

			// 1. Controllo: Non posso aggiungere me stesso
			// (Assumiamo che tu abbia l'ID o username corrente salvato da qualche parte, es. localStorage)
			const myUsername = localStorage.getItem("username"); // Esempio
			if (val === myUsername) {
				this.errormsg = "Non puoi aggiungere te stesso.";
				return;
			}

			// 2. Controllo Duplicati Locale (è già in lista?)
			if (this.groupMembers.includes(val)) {
				this.errormsg = "Utente già presente in lista.";
				return;
			}

			// 3. VERIFICA ONLINE (Chiamata al Backend)
			this.addingMember = true;
			this.errormsg = null;

			try {
				// Chiamo l'API. Qui uso l'URL standard REST.
				// Se il tuo endpoint è diverso, cambialo qui.
				// Ci aspettiamo che restituisca 200 OK se esiste, 404 se no.
				await this.$axios.get("/users/" + val);

				// Se arriviamo qui, l'utente ESISTE. Lo aggiungiamo.
				this.groupMembers.push(val);
				this.memberInput = ''; // Pulisco il campo

			} catch (e) {
				if (e.response && e.response.status === 404) {
					this.errormsg = `L'utente '${val}' non esiste.`;
				} else {
					console.error(e);
					this.errormsg = "Errore di connessione.";
				}
			} finally {
				this.addingMember = false; // Fermo lo spinner
			}
		},

		removeMemberFromList(index) {
			this.groupMembers.splice(index, 1);
		},

		async confirm() {
			if (!this.isValid) return;
			this.loading = true;
			this.errormsg = null;

			try {
				let payload = {
					is_group: this.isGroup,
				};

				if (this.isGroup) {
					payload.group_name = this.groupName;
					payload.participant_ids = this.groupMembers;
				} else {
					// Per la chat privata faccio un controllo rapido anche qui se vuoi
					// oppure lascio che fallisca se non esiste
					payload.participant_ids = [this.targetName];
				}

				await this.$axios.post("/conversations", payload);

				this.$emit('chat-created');
				this.$emit('close');
			} catch (e) {
				this.errormsg = e.response?.data || "Errore durante la creazione.";
			} finally {
				this.loading = false;
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
.modal-dialog { background: white; border-radius: 12px; min-width: 400px; max-width: 90%; }
</style>
