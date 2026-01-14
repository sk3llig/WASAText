<template>
	<div class="col-md-4 col-lg-3 d-flex flex-column sidebar-dark border-end-dark h-100">

		<div class="p-3 border-bottom-dark flex-shrink-0">
			<button class="btn btn-brand w-100 shadow-sm fw-bold" @click="$emit('open-new-chat')">
				+ Nuova Chat
			</button>
		</div>

		<div class="flex-grow-1 overflow-auto custom-scrollbar">

			<div v-if="conversations.length === 0" class="text-center mt-4 text-muted-dark small">
				Nessuna conversazione attiva.
			</div>

			<div
				v-for="conv in conversations"
				:key="conv.id"
				class="p-3 border-bottom-dark position-relative chat-item"
				:class="{ 'active-chat': activeConvId === conv.id }"
				style="cursor: pointer;"
				@click="$emit('select-chat', conv)"
			>
				<div class="d-flex align-items-center">

					<div class="position-relative me-3">
						<img :src="conv.photo_url || 'https://placehold.co/50'" class="rounded-circle object-fit-cover border-dark-profile" width="50" height="50">
					</div>

					<div class="flex-grow-1 overflow-hidden">
						<div class="d-flex justify-content-between align-items-baseline mb-1">
							<h6 class="mb-0 text-white text-truncate fw-bold pe-2" style="font-size: 0.95rem;">
								{{ conv.name || 'Chat' }}
							</h6>

							<small
								:class="conv.unread_count > 0 ? 'text-accent fw-bold' : 'text-muted-light'"
								class="timestamp-meta"
							>
								{{ formatSmartTime(conv.last_message_time) }}
							</small>
						</div>

						<div class="d-flex justify-content-between align-items-center">
							<p class="mb-0 text-truncate pe-4"
							   :class="conv.unread_count > 0 ? 'text-white fw-bold' : 'text-muted-dark'"
							   style="font-size: 0.85rem;">

								<span v-if="conv.last_message && isMe(conv.last_message_sender_id)" class="fw-normal opacity-75">Tu: </span>

								<span v-if="conv.last_message">{{ conv.last_message }}</span>
								<span v-else class="fst-italic opacity-50">Nessun messaggio</span>
							</p>

							<span v-if="conv.unread_count > 0" class="badge rounded-pill bg-accent shadow-sm">
                        {{ conv.unread_count }}
                      </span>
						</div>
					</div>
				</div>

				<button
					class="btn btn-sm btn-icon-danger position-absolute end-0 top-50 translate-middle-y me-2 delete-btn"
					title="Elimina conversazione"
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
	props: ['conversations', 'activeConvId', 'myId'],
	emits: ['open-new-chat', 'select-chat', 'delete-chat'],

	methods: {
		isMe(senderId) {
			if (!this.myId || !senderId) return false;
			return String(senderId) === String(this.myId);
		},

		formatSmartTime(isoDate) {
			if (!isoDate || String(isoDate).startsWith('0001')) return '';

			const date = new Date(isoDate);
			if (isNaN(date.getTime())) return '';

			const now = new Date();
			if (date.toDateString() === now.toDateString()) {
				return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
			}
			const yesterday = new Date(now);
			yesterday.setDate(now.getDate() - 1);
			if (date.toDateString() === yesterday.toDateString()) {
				return 'Ieri';
			}
			return date.toLocaleDateString([], { day: '2-digit', month: '2-digit', year: '2-digit' });
		}
	}
}
</script>

<style scoped>
/* --- TEMA SCURO --- */
.sidebar-dark { background-color: #1a1d21; color: white; }
.border-end-dark { border-right: 1px solid #2c3036; }
.border-bottom-dark { border-bottom: 1px solid #2c3036; }
.border-dark-profile { border: 2px solid #2c3036; }

/* ELEMENTO CHAT */
.chat-item {
	transition: background-color 0.2s;
	border-left: 4px solid transparent;
}
.chat-item:hover { background-color: #2c3036; }
.active-chat { background-color: #25282e; border-left: 4px solid #3282b8; }

/* FIX ORARIO E CESTINO */
.timestamp-meta {
	font-size: 0.75rem;
	margin-right: 28px; /* Crea lo spazio esatto per il cestino */
	white-space: nowrap;
}

/* TESTI */
.text-muted-light { color: #aebac6; }
.text-muted-dark { color: #8b9bb4; }
.text-accent { color: #3282b8; }
.bg-accent { background-color: #3282b8; color: white; }

.badge { font-size: 0.7rem; padding: 0.35em 0.6em; }

.btn-brand {
	background-color: #3282b8;
	color: white;
	border: none;
	transition: transform 0.1s;
}
.btn-brand:hover { background-color: #0f4c75; transform: scale(1.02); }

/* CESTINO */
.delete-btn {
	display: none;
	background: transparent;
	border: none;
	font-size: 1.1rem;
	z-index: 5;
}
.chat-item:hover .delete-btn { display: block; }

/* SCROLLBAR */
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: #1a1d21; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #383d45; border-radius: 3px; }

.object-fit-cover { object-fit: cover; }
</style>
