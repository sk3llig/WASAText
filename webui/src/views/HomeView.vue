<template>
  <div class="container-fluid vh-100 d-flex flex-column overflow-hidden p-0">
    <GroupInfoModal
      v-if="showGroupInfoModal && activeConv"
      :group-id="activeConv.id"
      :current-name="activeConv.name"
      @close="showGroupInfoModal = false"
      @name-updated="loadConversations"
    />

    <NavBar
      :my-username="myUsername"
      :my-photo="myPhoto"
      @logout="logout"
      @open-profile="showProfileModal = true"
    />

    <div class="row flex-grow-1 g-0 overflow-hidden">
      <ChatSidebar
        :conversations="conversations"
        :active-conv-id="activeConv ? activeConv.id : null"
        @select-chat="selectConversation"
        @open-new-chat="showChatModal = true"
        @delete-chat="deleteConversation"
      />

      <ChatWindow
        :active-conv="activeConv"
        :messages="messages"
        :my-username="myUsername"
        @send-message="sendMessage"
        @send-image="handleImageUpload"
        @update-group-photo="handleGroupPhotoUpdate"
        @leave-group="leaveGroup"
        @add-reaction="addReaction"
        @remove-reaction="removeReaction"
        @delete-message="deleteMessage"
        @open-group-info="showGroupInfoModal = true"
        @forward-message="openForwardModal"
		@send-composite-message="handleCompositeMessage"
		@mark-as-read="markAsReadAPI"
      />
    </div>

	  <ProfileModal
		  v-if="showProfileModal"
		  :my-username="myUsername"
		  :my-id="myId"
		  :my-photo="myPhoto"
		  @close="showProfileModal = false"
		  @update:myUsername="myUsername = $event"
		  @update-photo="uploadPhoto"
	  />

	  <NewChatModal
		  v-if="showChatModal"
		  :error-msg="chatModalError"    @close="showChatModal = false"
		  @create-private="startPrivateChat"
		  @create-group="startGroupChat"
	  />

    <ForwardModal
      v-if="showForwardModal"
      @close="showForwardModal = false"
      @forward-to="handleForward"
    />
  </div>
</template>

<script>
import axios from '../services/axios';
// Importazione Componenti
import NavBar from '../components/NavBar.vue';
import ChatSidebar from '../components/ChatSidebar.vue';
import ChatWindow from '../components/ChatWindow.vue';
import ProfileModal from '../components/ProfileModal.vue';
import NewChatModal from '../components/NewChatModal.vue';
import GroupInfoModal from '../components/GroupInfoModal.vue';
import ForwardModal from '../components/ForwardModal.vue';

export default {
  components: { NavBar, ChatSidebar, ChatWindow, ProfileModal, NewChatModal, GroupInfoModal ,ForwardModal},
  data() {
    return {
      myUsername: localStorage.getItem('username'),
      myId: localStorage.getItem('token'), // Nel simplified login l'ID è il token
      myPhoto: null,
      conversations: [],
      messages: [],
      activeConv: null,
      pollInterval: null,

      // Stati Modali
      showProfileModal: false,
      showChatModal: false,
      showGroupInfoModal: false,
      showForwardModal: false,
		msgToForward: null,
		chatModalError: '',
    }
  },
  async mounted() {
    await this.loadMyProfile();
    await this.loadConversations();
    this.pollInterval = setInterval(() => {
      this.loadConversations();
      if (this.activeConv) this.loadMessages();
    }, 2000);
  },
  beforeUnmount() { clearInterval(this.pollInterval); },

  methods: {
    // --- LOGICA DATI ---
	  // In HomeView.vue
		  async markAsReadAPI(conversation) {
			  console.log("Lettura implicita gestita dal backend al prossimo refresh");
		  },
		  // ...

    async loadMyProfile() {
      try {
        let res = await axios.get('/me');
        this.myPhoto = res.data.photo_url;
      } catch (e) {}
    },
    async loadConversations() {
      try {
        let res = await axios.get('/me/conversations');
        this.conversations = res.data.conversations || [];
      } catch (e) { console.error(e); }
    },
    async selectConversation(conv) {
      this.activeConv = conv;
      await this.loadMessages();
    },
    async loadMessages() {
      if (!this.activeConv) return;
      try {
        let res = await axios.get(`/conversations/${this.activeConv.id}/messages`);
        this.messages = (res.data.messages || []).reverse();
      } catch (e) {}
    },

    // --- AZIONI UTENTE ---
    async sendMessage(text) {
      try {
        await axios.post(`/conversations/${this.activeConv.id}/messages`, { content: text, attachments: [] });
        this.loadMessages();
      } catch(e) { alert("Errore invio"); }
    },

	  // Dentro methods di HomeView.vue

	  async uploadPhoto(file) {
		  if (!file) return;

		  try {
			  console.log("Avvio caricamento file:", file.name);

			  // 1. Creiamo il payload Multipart
			  let formData = new FormData();

			  // --- PUNTO CRUCIALE ---
			  // Il tuo codice Go dice: r.FormFile("file")
			  // Quindi qui DOBBIAMO usare "file", non "image" o "photo".
			  formData.append("file", file);
			  // ----------------------

			  // 2. Upload fisico dell'immagine (POST /images)
			  // Questa chiamata salva il file su disco e ci restituisce l'URL
			  const uploadResponse = await this.$axios.post("/images", formData, {
				  headers: {
					  'Content-Type': 'multipart/form-data'
				  }
			  });

			  // 3. Recuperiamo l'URL generato dal backend
			  // Il tuo Go risponde con la struct ImageMetadata { ID, URL, ContentType }
			  // Solitamente JSON trasforma i campi in minuscolo o come definiti nei tag.
			  // Proviamo 'url' (standard) o 'URL'.
			  const newPhotoUrl = uploadResponse.data.url || uploadResponse.data.URL;

			  if (!newPhotoUrl) {
				  throw new Error("Il server non ha restituito l'URL dell'immagine.");
			  }

			  console.log("Immagine salvata. URL:", newPhotoUrl);

			  // 4. Aggiorniamo il profilo utente con il nuovo URL (PUT /me/photo)
			  // Il tuo backend Go si aspetta un JSON con SetPhotoRequest
			  await this.$axios.put("/me/photo", {
				  photo_url: newPhotoUrl
			  });

			  // 5. Aggiorniamo l'interfaccia locale subito
			  if (this.myProfile) {
				  this.myProfile.photo_url = newPhotoUrl;
			  }
			  this.myPhoto = newPhotoUrl + '?t=' + new Date().getTime();
			  alert("Foto aggiornata con successo!");

		  } catch (e) {
			  console.error("Errore durante l'upload:", e);

			  // Gestione errori basata sui messaggi del tuo Backend Go
			  if (e.response) {
				  if (e.response.status === 400) {
					  alert("Errore: " + (e.response.data || "Formato non valido o file troppo grande (Max 10MB, solo JPG/PNG)"));
				  } else if (e.response.status === 401) {
					  alert("Sessione scaduta. Esegui il login.");
				  } else {
					  alert("Errore del server: " + e.response.status);
				  }
			  } else {
				  alert("Errore di connessione.");
			  }
		  }
	  },
    // --- GESTIONE CHAT ---
	  async handleCompositeMessage({ text, file }) {
		  console.log("Invio messaggio combinato...", text, file);

		  try {
			  if (!this.activeConv || !this.activeConv.id) {
				  console.error("Nessuna conversazione attiva!");
				  return;
			  }

			  // 1. UPLOAD FOTO (Se c'è)
			  let uploadedAttachments = []; // Deve essere un array!
			  let messageType = "text";     // Default

			  if (file) {
				  messageType = "photo"; // Se c'è un file, il tipo diventa photo

				  let formData = new FormData();
				  formData.append("file", file);

				  // Carichiamo l'immagine
				  const uploadRes = await this.$axios.post("/images", formData, {
					  headers: { 'Content-Type': 'multipart/form-data' }
				  });

				  // Recuperiamo l'URL.
				  // Nota: adattalo se il backend risponde con .url o .URL
				  const url = uploadRes.data.url || uploadRes.data.URL;

				  // La struct vuole un array di stringhe in "attachments"
				  uploadedAttachments.push(url);
			  }

			  // 2. PREPARAZIONE PAYLOAD (Esattamente come vuole la struct Go)
			  const payload = {
				  conversation_id: this.activeConv.id,  // Corretto da struct: `json:"conversation_id"`
				  content: text || "",                  // Corretto da struct: `json:"content"`
				  attachments: uploadedAttachments,     // Corretto da struct: `json:"attachments"` ([]string)
				  type: messageType                     // Corretto da struct: `json:"type"` ("text" o "photo")
			  };

			  console.log("Payload finale:", payload);

			  // 3. INVIO MESSAGGIO
			  // IMPORTANTE: Se ti dà ancora 404, verifica l'URL nel punto 4 qui sotto
			  await this.$axios.post(`/conversations/${this.activeConv.id}/messages`, payload);
			  console.log("Messaggio inviato con successo!");

			  // Aggiorna la chat (chiama la tua funzione di fetch o resetta)
			  // this.fetchMessages();

		  } catch (e) {
			  console.error("Errore invio:", e);
			  alert("Errore durante l'invio del messaggio.");
		  }
	  },
	  async findUserId(name) {
		  // 1. Resettiamo l'errore ogni volta che cerchiamo
		  this.chatModalError = '';

		  try {
			  // Usiamo axios per cercare l'utente (endpoint /session o /users a seconda del tuo backend)
			  let res = await axios.post('/session', { name: name });

			  // Controllo se sono io
			  if (res.data.identifier === this.myId) {
				  this.chatModalError = "Non puoi chattare con te stesso!";
				  return null;
			  }

			  return res.data.identifier;

		  } catch (e) {
			  // 2. GESTIONE SILENZIOSA (Niente console rossa per errori 404)
			  if (e.response && (e.response.status === 404 || e.response.status === 400)) {
				  this.chatModalError = "⚠️ Utente '" + name + "' non trovato.";
			  } else {
				  // Stampiamo in console solo se è un errore di server grave
				  console.error(e);
				  this.chatModalError = "Errore di connessione.";
			  }
			  return null;
		  }
	  },

    async startPrivateChat(username) {
      let uid = await this.findUserId(username);
      if (!uid) return;
      this.createConvAPI([uid], false, null);
    },

    async startGroupChat({ groupName, memberName }) {
      let uid = await this.findUserId(memberName);
      if (!uid) return;
      this.createConvAPI([uid], true, groupName);
    },

    async createConvAPI(ids, isGroup, name) {
      try {
        let res = await axios.post('/conversations', {
          participant_ids: ids, is_group: isGroup, group_name: name
        });
        this.showChatModal = false;
        await this.loadConversations();
        // Seleziona la nuova chat
        let newConv = this.conversations.find(c => c.id === res.data.id);
        if (newConv) this.selectConversation(newConv);
      } catch (e) { alert("Errore creazione chat"); }
    },

    async deleteConversation(conv) {
      if (!confirm("Eliminare chat?")) return;
      try { await axios.delete(`/conversations/${conv.id}`); this.activeConv = null; this.loadConversations(); } catch (e) {}
    },
    async leaveGroup() {
      if (!confirm("Uscire dal gruppo?")) return;
      try { await axios.delete(`/groups/${this.activeConv.id}/members/me`); this.activeConv = null; this.loadConversations(); } catch (e) {}
    },

    // --- MESSAGGI ---
    async deleteMessage(msg) {
      if (!confirm("Eliminare messaggio?")) return;
      try { await axios.delete(`/conversations/${this.activeConv.id}/messages/${msg.id}`); this.loadMessages(); } catch (e) { alert("Non puoi eliminare questo messaggio."); }
    },
    async addReaction(msg, emoji) {
      try { await axios.post(`/conversations/${this.activeConv.id}/messages/${msg.id}/reactions`, { emoji }); this.loadMessages(); } catch (e) {}
    },


    async removeReaction(msg, reaction) {
      try { await axios.delete(`/conversations/${this.activeConv.id}/messages/${msg.id}/reactions`, { params: { emoji: reaction.emoji } }); this.loadMessages(); } catch (e) {}
    },

    // --- UPLOAD ---
    async uploadFileGeneric(file) {
      let formData = new FormData(); formData.append('file', file);
      let res = await axios.post('/images', formData);
      return res.data.url;
    },
    async handleImageUpload(file) {
      try {
        let url = await this.uploadFileGeneric(file);
        await axios.post(`/conversations/${this.activeConv.id}/messages`, { content: '', attachments: [url] });
        this.loadMessages();
      } catch (e) { alert("Errore upload"); }
    },
    async handleProfileUpdate(file) {
		try {
			// Carica la foto sul server e ottieni il nuovo URL
			let url = await this.uploadFileGeneric(file);

			// Aggiorna il DB del backend
			await axios.put('/me/photo', { photo_url: url });

			// Vue vedrà il cambiamento e aggiornerà subito l'immagine nella NavBar.
			this.myPhoto = url;

			// Chiudi il modale
			this.showProfileModal = false;

		} catch (e) {
			console.error(e);
			alert("Errore aggiornamento foto");
		}
	},
    async handleGroupPhotoUpdate(file) {
      try {
        let url = await this.uploadFileGeneric(file);
        await axios.put(`/groups/${this.activeConv.id}/photo`, { photo_url: url });
        this.activeConv.photo_url = url; // Update locale immediato
        this.loadConversations();
      } catch (e) { alert("Errore foto gruppo"); }
    },

    logout() { localStorage.clear(); this.$router.push('/login'); },

    // 1. Apre il modale quando clicchi la freccia nella ChatWindow
    openForwardModal(msg) {
      this.msgToForward = msg;
      this.showForwardModal = true;
    },

    // 2. Esegue l'inoltro vero e proprio
    async handleForward({ targetId, comment }) {
      if (!this.msgToForward) return;

      try {
        await axios.post(`/conversations/${this.activeConv.id}/messages/${this.msgToForward.id}/forward`, {
          target_conversation_id: targetId,
          comment: comment
        });

        this.showForwardModal = false;
        this.msgToForward = null;
        alert("Messaggio inoltrato!");

        // Se hai inoltrato alla chat che stavi guardando, ricarica i messaggi
        if (this.activeConv.id === targetId) {
          this.loadMessages();
        }
      } catch (e) {
        alert("Errore inoltro: " + (e.response?.data?.message || e.message));
      }
    },
}
  }

</script>
