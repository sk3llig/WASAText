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

			loading: false,
			addingMember: false,
			errormsg: null
		}
	},

	computed: {
		isValid() {
			if (!this.isGroup) {
				return this.targetName.length > 0;
			} else {
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

		async addMemberToList() {
			const val = this.memberInput.trim();
			if (!val) return;

			const myUsername = localStorage.getItem("username");
			if (val === myUsername) {
				this.errormsg = "Non puoi aggiungere te stesso.";
				return;
			}

			if (this.groupMembers.includes(val)) {
				this.errormsg = "Utente già presente in lista.";
				return;
			}

			this.addingMember = true;
			this.errormsg = null;

			try {
				// Nota: Assicurati che il backend accetti GET /users/{username}
				// Se il tuo backend richiede POST /session per cercare l'ID, va adattato qui.
				await this.$axios.get("/users/" + val);

				this.groupMembers.push(val);
				this.memberInput = '';

			} catch (e) {
				// GESTIONE SILENZIOSA ERRORE GRUPPO
				if (e.response && e.response.status === 404) {
					// Utente non esiste: Mostriamo messaggio ma NON stampiamo console.error
					this.errormsg = `⚠️ L'utente '${val}' non esiste.`;
				} else {
					// Altri errori (es. server down): Stampiamo l'errore reale
					console.error("Errore aggiunta membro:", e);
					this.errormsg = "Errore di connessione.";
				}
			} finally {
				this.addingMember = false;
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
					payload.participant_ids = [this.targetName];
				}

				// Chiamata per creare la conversazione
				await this.$axios.post("/conversations", payload);

				this.$emit('chat-created');
				this.$emit('close');

			} catch (e) {

				// Se l'errore è 400 (Bad Request) o 404 (Not Found)
				if (e.response && (e.response.status === 400 || e.response.status === 404)) {

					// Messaggio personalizzato invece del JSON brutto
					if (!this.isGroup) {
						this.errormsg = `⚠️ Utente '${this.targetName}' non trovato.`;
					} else {
						this.errormsg = "⚠️ Uno o più utenti non sono validi.";
					}

					// NON facciamo console.error(e) qui!
					// In console vedrai solo la riga rossa "POST ... 404" del browser (che è normale)

				} else {
					// Errore imprevisto (es. 500 Server Error)
					console.error("Errore creazione chat:", e);
					this.errormsg = "Errore del server. Riprova più tardi.";
				}
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
