<template>
  <div class="col-md-8 col-lg-9 d-flex flex-column bg-white h-100">
    
    <div v-if="activeConv" class="d-flex flex-column h-100">
      
      <div class="p-3 border-bottom bg-light d-flex justify-content-between align-items-center shadow-sm z-10">
        <div class="d-flex align-items-center">
          <img :src="activeConv.photo_url || 'https://placehold.co/50'" class="rounded-circle me-2 object-fit-cover" width="40" height="40">
          <div>
            <h6 class="mb-0">{{ activeConv.name }}</h6>
            <small class="text-muted" v-if="activeConv.is_group">Gruppo</small>
          </div>
        </div>
        <div v-if="activeConv.is_group">
          <button class="btn btn-outline-primary btn-sm me-2" @click="$emit('open-group-info')">ℹ️ Info</button> 

          <button class="btn btn-outline-secondary btn-sm me-2" @click="$refs.groupInput.click()">📷 Foto</button>
          <button class="btn btn-outline-danger btn-sm" @click="$emit('leave-group')">🚪 Esci</button>
        </div>
      </div>

      <div class="flex-grow-1 p-4 overflow-auto bg-chat" id="chat-window" ref="scrollContainer">
        <div v-if="messages.length === 0" class="text-center mt-5">Inizia a scrivere!</div>

        <div v-for="msg in messages" :key="msg.id" class="d-flex mb-3" :class="isMyMessage(msg) ? 'justify-content-end' : 'justify-content-start'">
          <div class="msg-container d-flex flex-column" :class="isMyMessage(msg) ? 'align-items-end' : 'align-items-start'">
            <div class="msg-bubble p-2 px-3 shadow-sm" :class="isMyMessage(msg) ? 'bg-primary text-white my-msg' : 'bg-white border other-msg'">
              <small v-if="activeConv.is_group && !isMyMessage(msg)" class="d-block fw-bold mb-1 text-secondary">{{ msg.sender.name }}</small>
              
              <div v-if="msg.attachments && msg.attachments.length > 0" class="mb-2">
                <img :src="msg.attachments[0]" class="img-fluid rounded" style="max-height: 200px;">
              </div>
              <div v-if="msg.content">{{ msg.content }}</div>
              <small class="d-block text-end mt-1 opacity-75" style="font-size: 0.65rem;">{{ formatTime(msg.timestamp) }}</small>
              
              <div v-if="msg.reactions && msg.reactions.length > 0" class="mt-1">
                <span v-for="r in msg.reactions" :key="r.user.id" class="badge bg-light text-dark border me-1" @click="$emit('remove-reaction', msg, r)" style="cursor:pointer">
                  {{ r.emoji }}
                </span>
              </div>
            </div>
           <div class="mt-1">
              <button class="btn btn-sm p-0 me-1" @click="$emit('add-reaction', msg, '👍')">👍</button>
              <button class="btn btn-sm p-0 me-1" @click="$emit('add-reaction', msg, '❤️')">❤️</button>
              <button class="btn btn-sm p-0 me-1" @click="$emit('add-reaction', msg, '😂')">😂</button>
              
              <button class="btn btn-sm p-0 text-secondary ms-2" @click="$emit('forward-message', msg)" title="Inoltra">
                ➡️
              </button>

              <button v-if="isMyMessage(msg)" class="btn btn-sm p-0 text-danger ms-2" @click="$emit('delete-message', msg)" title="Elimina">
                🗑️
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="p-3 border-top bg-light">
        <div class="input-group">
          <button class="btn btn-outline-secondary" @click="$refs.fileInput.click()">📷</button>
          <input type="text" class="form-control" placeholder="Scrivi..." v-model="newMessageText" @keyup.enter="send">
          <button class="btn btn-primary" @click="send">Invia</button>
        </div>
      </div>
    </div>

    <div v-else class="d-flex flex-column align-items-center justify-content-center h-100 bg-light text-muted">
      <h3>💬</h3><p>Seleziona una chat</p>
    </div>

    <input type="file" ref="fileInput" class="d-none" accept="image/*" @change="handleFile">
    <input type="file" ref="groupInput" class="d-none" accept="image/*" @change="handleGroupPhoto">
  </div>
</template>

<script>
export default {
  props: ['activeConv', 'messages', 'myUsername'],
  emits: ['send-message', 'send-image', 'update-group-photo', 'leave-group', 'add-reaction', 'remove-reaction', 'delete-message','open-group-info','forward-message'],
  data() {
    return { newMessageText: '' }
  },
  watch: {
    // Quando cambiano i messaggi, scrolla in basso
    messages: {
      handler() { this.scrollToBottom() },
      deep: true
    }
  },
  methods: {
    send() {
      if(!this.newMessageText.trim()) return;
      this.$emit('send-message', this.newMessageText);
      this.newMessageText = '';
    },
    handleFile(e) {
      if(e.target.files[0]) this.$emit('send-image', e.target.files[0]);
    },
    handleGroupPhoto(e) {
      if(e.target.files[0]) this.$emit('update-group-photo', e.target.files[0]);
    },
    isMyMessage(msg) { return msg.sender.name === this.myUsername; },
    formatTime(iso) { 
        if(!iso) return ''; 
        return new Date(iso).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}); 
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const el = document.getElementById('chat-window');
        if (el) el.scrollTop = el.scrollHeight;
      });
    }
  }
}
</script>

<style scoped>
/* Copia qui gli stili .bg-chat, .msg-bubble ecc. dal vecchio HomeView */
.bg-chat { background-color: #f0f2f5; }
.msg-bubble { max-width: 75%; border-radius: 15px; }
.my-msg { border-top-right-radius: 0; }
.other-msg { border-top-left-radius: 0; }
.object-fit-cover { object-fit: cover; }
</style>