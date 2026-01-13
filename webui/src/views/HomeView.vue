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
      @close="showChatModal = false"
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

	  async uploadPhoto(file) {
		  console.log("TODO: Implementare caricamento foto", file);
		  alert("Funzionalità caricamento foto non ancora implementata!");
	  },

    // --- GESTIONE CHAT ---
    async findUserId(name) {
      try {
        let res = await axios.post('/session', { name: name });
        if (res.data.identifier === this.myId) {
          alert("Non puoi chattare con te stesso!"); return null;
        }
        return res.data.identifier;
      } catch (e) { alert("Utente " + name + " non trovato."); return null; }
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
