<template>
	<div class="col-md-8 col-lg-9 d-flex flex-column bg-white h-100 position-relative">

		<div v-if="activeConv" class="d-flex flex-column h-100">
			<div class="p-3 border-bottom bg-light d-flex justify-content-between align-items-center shadow-sm z-10">
				<div class="d-flex align-items-center">
					<img :src="activeConv.photo_url || 'https://placehold.co/50'" class="rounded-circle me-2 object-fit-cover" width="40" height="40">
					<div>
						<h6 class="mb-0">{{ activeConv.name }}</h6>
						<small v-if="activeConv.is_group" class="text-muted">Gruppo</small>
					</div>
				</div>
				<div v-if="activeConv.is_group">
					<button class="btn btn-outline-primary btn-sm me-2" @click="$emit('open-group-info')">ℹ️ Info</button>
					<button class="btn btn-outline-secondary btn-sm me-2" @click="$refs.groupInput.click()">📷 Foto</button>
					<button class="btn btn-outline-danger btn-sm" @click="$emit('leave-group')">🚪 Esci</button>
				</div>
			</div>
			<div id="chat-window" ref="scrollContainer" class="flex-grow-1 p-4 overflow-auto bg-chat" @click="openPickerId = null">
				<div v-if="messages.length === 0" class="text-center mt-5">Inizia a scrivere!</div>

				<div v-for="msg in sortedMessages" :key="msg.id" class="d-flex mb-3" :class="isMyMessage(msg) ? 'justify-content-end' : 'justify-content-start'">
					<div class="msg-container d-flex flex-column position-relative" :class="isMyMessage(msg) ? 'align-items-end' : 'align-items-start'">

						<div class="msg-bubble p-2 px-3 shadow-sm" :class="isMyMessage(msg) ? 'bg-primary text-white my-msg' : 'bg-white border other-msg'">
							<small v-if="activeConv.is_group && !isMyMessage(msg)" class="d-block fw-bold mb-1 text-secondary">{{ msg.sender.name }}</small>

							<div v-if="msg.attachments && msg.attachments.length > 0" class="mb-2">
								<img :src="msg.attachments[0]" class="img-fluid rounded" style="max-height: 200px;">
							</div>
							<div v-if="msg.content">{{ msg.content }}</div>
							<div class="d-flex justify-content-end align-items-center mt-1">
								<small class="opacity-75 me-1" style="font-size: 0.65rem;">
									{{ formatTime(msg.timestamp) }}
								</small>

								<div v-if="isMyMessage(msg)" class="lh-1">
        <span v-if="msg.status === 'read'" class="text-info fw-bold" style="font-size: 0.8rem;">
            ✓✓
        </span>
									<span v-else class="text-white-50 fw-bold" style="font-size: 0.8rem;">
            ✓✓
        </span>
								</div>
							</div>
							<div v-if="msg.reactions && msg.reactions.length > 0" class="reaction-bar mt-2">
								<div class="badge bg-light text-dark border rounded-pill shadow-sm d-inline-flex align-items-center px-2 py-1"
									 style="cursor: pointer;"
									 @click.stop="openReactionDetails(msg)">
                             <span v-for="emoji in getUniqueEmojis(msg.reactions)" :key="emoji" class="me-1">
                                {{ emoji }}
                             </span>
									<span class="small ms-1 text-muted">{{ msg.reactions.length }}</span>
								</div>
							</div>
						</div>

						<div class="mt-1 d-flex align-items-center position-relative">
							<button class="btn btn-sm text-secondary p-1 rounded-circle hover-bg me-1"
									title="Aggiungi reazione"
									@click.stop="togglePicker(msg.id)">
								☺
							</button>

							<div v-if="openPickerId === msg.id"
								 class="emoji-picker shadow-lg p-2 rounded bg-white border"
								 :class="isMyMessage(msg) ? 'picker-right' : 'picker-left'">
								<div class="d-flex gap-2">
									<button v-for="emoji in availableEmojis" :key="emoji"
											class="btn btn-sm btn-light p-1 fs-5"
											@click="addReaction(msg, emoji)">
										{{ emoji }}
									</button>
								</div>
							</div>

							<button class="btn btn-sm p-0 text-secondary ms-2" title="Inoltra" @click="$emit('forward-message', msg)">➡️</button>
							<button v-if="isMyMessage(msg)" class="btn btn-sm p-0 text-danger ms-2" title="Elimina" @click="$emit('delete-message', msg)">🗑️</button>
						</div>
					</div>
				</div>
			</div>

			<div class="p-3 border-top bg-light">
				<div v-if="previewUrl" class="mb-2 position-relative d-inline-block">
					<img :src="previewUrl" class="rounded border" style="height: 80px; width: auto; object-fit: cover;">
					<button @click="removeAttachment" class="btn btn-danger btn-sm position-absolute top-0 end-0 rounded-circle d-flex align-items-center justify-content-center" style="width: 20px; height: 20px; font-size: 12px; transform: translate(30%, -30%);">✕</button>
				</div>
				<div class="input-group">
					<button class="btn btn-outline-secondary" @click="$refs.fileInput.click()">📷</button>
					<input v-model="newMessageText" type="text" class="form-control" placeholder="Scrivi..." @keyup.enter="send">
					<button class="btn btn-primary" @click="send">Invia</button>
				</div>
			</div>
		</div>

		<div v-else class="d-flex flex-column align-items-center justify-content-center h-100 bg-light text-muted">
			<h3>💬</h3><p>Seleziona una chat</p>
		</div>

		<input ref="fileInput" type="file" class="d-none" accept="image/*" @click="e => e.target.value = null" @change="handleFile">
		<input ref="groupInput" type="file" class="d-none" accept="image/*" @change="handleGroupPhoto">

		<div v-if="viewingReactionsMsg" class="modal-backdrop" @click.self="closeReactionDetails">

			<div class="modal-dialog bg-white rounded shadow-lg p-0 overflow-hidden"
				 style="max-width: 350px; width: 100%; pointer-events: auto;"
				 @click.stop>

				<div class="modal-header p-3 border-bottom d-flex justify-content-between align-items-center bg-light">
					<h6 class="m-0 fw-bold text-dark">Reazioni</h6>
					<button class="btn-close small" @click="closeReactionDetails"></button>
				</div>

				<div class="modal-body p-0" style="max-height: 400px; overflow-y: auto;">
					<div v-for="(items, emoji) in groupedReactions" :key="emoji" class="p-3 border-bottom">
						<div class="fs-5 mb-2">{{ emoji }} <span class="fs-6 text-muted">({{ items.length }})</span></div>

						<div v-for="(item, index) in items" :key="index" class="d-flex align-items-center mb-2">

							<img src="https://placehold.co/30" class="rounded-circle me-2 object-fit-cover border" width="30" height="30">

							<span class="small ms-2 fw-bold text-dark">
                         {{ item.user ? item.user.name : 'Sconosciuto' }}
                      </span>

							<button
								v-if="item.user && item.user.name === myUsername"
								class="btn btn-light text-danger ms-auto d-flex align-items-center justify-content-center shadow-sm"
								style="cursor: pointer; width: 32px; height: 32px; border-radius: 50%; transition: all 0.2s;"
								title="Elimina la tua reazione"
								@click.stop="deleteReactionInModal(item)"
								onmouseover="this.style.backgroundColor='#ffeaea'"
								onmouseout="this.style.backgroundColor='#f8f9fa'">
								<span style="font-size: 14px;">🗑️</span>
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
		```



	</div>
</template>

<script>
export default {
	props: ['activeConv', 'messages', 'myUsername'],
	emits: ['mark-as-read','send-message', 'send-image', 'send-composite-message', 'update-group-photo', 'leave-group', 'add-reaction', 'remove-reaction', 'delete-message','open-group-info','forward-message'],

	data() {
		return {
			newMessageText: '',
			selectedFile: null,
			previewUrl: null,
			openPickerId: null,
			viewingReactionsMsg: null,
			availableEmojis: ['👍', '❤️', '😂', '😮', '😢', '🔥', '🎉', '💩'],
			initialScrollDone: false
		}
	},

	computed: {
		// --- MODIFICATO PER LA STRUTTURA SEMPLICE ---
		groupedReactions() {
			if (!this.viewingReactionsMsg || !this.viewingReactionsMsg.reactions) return {};

			const groups = {};
			this.viewingReactionsMsg.reactions.forEach(r => {
				// r adesso è: { emoji: "👍", user_name: "Mario" }
				if (!groups[r.emoji]) groups[r.emoji] = [];
				groups[r.emoji].push(r);
			});
			return groups;
		},
		sortedMessages() {

			return [...this.messages].sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));
		}
	},

	watch: {
		messages: {
			handler(newMessages, oldMessages) {
				// 1. Scroll automatico (già presente)
				this.handleSmartScroll(newMessages, oldMessages);

				// 2. AGGIUNTA FONDAMENTALE: Aggiorna il modale se è aperto!
				if (this.viewingReactionsMsg) {
					// Cerchiamo la versione aggiornata del messaggio corrente
					const updatedMsg = newMessages.find(m => m.id === this.viewingReactionsMsg.id);

					if (updatedMsg) {
						// Aggiorniamo i dati del modale (così la reazione sparisce visivamente)
						this.viewingReactionsMsg = updatedMsg;
					} else {
						// Se il messaggio non esiste più (es. cancellato), chiudiamo il modale
						this.closeReactionDetails();
					}

				}
				this.triggerReadReceipt();
			},
			deep: true
		},
		// ... lascia activeConv.id invariato ...
		'activeConv.id'() {
			this.initialScrollDone = false;
			this.$nextTick(() => this.forceScrollToBottom());
			this.triggerReadReceipt();
		}

	},

	mounted() {
		this.forceScrollToBottom();
		this.triggerReadReceipt();
	},

	methods: {triggerReadReceipt() {
			// Se la chat è aperta e ci sono messaggi non miei,
			// dico al padre (HomeView) di segnarli come letti.
			if (this.activeConv && this.messages.length > 0) {
				this.$emit('mark-as-read', this.activeConv);
			}
		},
		handleSmartScroll(newMsgs, oldMsgs) {
			const el = this.$refs.scrollContainer;
			if (!el) return;
			const isAtBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 100;
			const lastMsg = newMsgs[newMsgs.length - 1];
			const isMyNewMessage = lastMsg && this.isMyMessage(lastMsg) && (!oldMsgs || oldMsgs.length !== newMsgs.length);

			this.$nextTick(() => {
				if (!this.initialScrollDone || isAtBottom || isMyNewMessage) {
					this.forceScrollToBottom();
					this.initialScrollDone = true;
				}
			});
		},
		deleteReactionInModal(reactionItem) {
			// viewingReactionsMsg è il messaggio che stiamo guardando nel modale
			// reactionItem contiene l'emoji e l'utente
			console.log("🗑️ Click cestino per:", reactionItem.emoji);
			this.$emit('remove-reaction', this.viewingReactionsMsg, reactionItem);
		},
		forceScrollToBottom() {
			const el = this.$refs.scrollContainer;
			if (el) el.scrollTop = el.scrollHeight;
		},

		getUniqueEmojis(reactions) {
			if(!reactions) return [];
			return [...new Set(reactions.map(r => r.emoji))].slice(0,3);
		},

		// --- METODI STANDARD ---
		send() {
			if(!this.newMessageText.trim() && !this.selectedFile) return;
			this.$emit('send-composite-message', {
				text: this.newMessageText,
				file: this.selectedFile
			});
			this.newMessageText = '';
			this.removeAttachment();
			this.$nextTick(() => this.forceScrollToBottom());
		},

		handleFile(e) {
			const file = e.target.files[0];
			if (!file) return;
			this.selectedFile = file;
			this.previewUrl = URL.createObjectURL(file);
		},
		removeAttachment() {
			this.selectedFile = null;
			if (this.previewUrl) {
				URL.revokeObjectURL(this.previewUrl);
				this.previewUrl = null;
			}
		},
		handleGroupPhoto(e) {
			if(e.target.files[0]) this.$emit('update-group-photo', e.target.files[0]);
		},
		isMyMessage(msg) { return msg.sender.name === this.myUsername; },
		formatTime(iso) {
			if(!iso) return '';
			return new Date(iso).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
		},

		togglePicker(msgId) {
			this.openPickerId = (this.openPickerId === msgId) ? null : msgId;
		},

		addReaction(msg, emoji) {
			this.$emit('add-reaction', msg, emoji);
			this.openPickerId = null;
		},

		openReactionDetails(msg) {
			this.viewingReactionsMsg = msg;
		},

		closeReactionDetails() {
			this.viewingReactionsMsg = null;
		}
	}
}
</script>

<style scoped>
.bg-chat { background-color: #f0f2f5; }
.msg-bubble { max-width: 75%; border-radius: 15px; }
.my-msg { border-top-right-radius: 0; }
.other-msg { border-top-left-radius: 0; }
.object-fit-cover { object-fit: cover; }
.hover-bg:hover { background-color: #e9ecef; }

.emoji-picker {
	position: absolute;
	bottom: 100%;
	z-index: 1000;
	min-width: 280px;
	animation: fadeIn 0.2s ease-out;
}

.picker-right { right: 0; left: auto; }
.picker-left { left: 0; right: auto; }

.modal-backdrop {
	position: fixed;
	top: 0; left: 0; width: 100%; height: 100%;
	background: rgba(0,0,0,0.5);
	z-index: 2000;
	display: flex;
	align-items: center;
	justify-content: center;
}

@keyframes fadeIn {
	from { opacity: 0; transform: translateY(10px); }
	to { opacity: 1; transform: translateY(0); }
}
</style>
