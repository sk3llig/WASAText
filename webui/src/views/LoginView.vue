<template>
  <div class="d-flex justify-content-center align-items-center vh-100 bg-light">
    <div class="card shadow p-4" style="width: 400px;">
      <h2 class="text-center mb-4">WASAText Login</h2>

      <div class="mb-3">
        <label for="username" class="form-label">Scegli uno Username</label>
        <input
          type="text"
          id="username"
          class="form-control"
          v-model="username"
          placeholder="Es. Alice"
          @keyup.enter="doLogin"
        >
      </div>

      <div v-if="errorMsg" class="alert alert-danger">
        {{ errorMsg }}
      </div>

      <button
        class="btn btn-primary w-100"
        @click="doLogin"
        :disabled="loading"
      >
        <span v-if="loading" class="spinner-border spinner-border-sm"></span>
        <span v-else>Entra</span>
      </button>
    </div>
  </div>
</template>

<script>
import axios from '../services/axios';

export default {
  data() {
    return {
      username: '',
      loading: false,
      errorMsg: ''
    }
  },
  methods: {
    async doLogin() {
      if (this.username.length < 3) {
        this.errorMsg = "Username troppo corto (min 3 caratteri)";
        return;
      }

      this.loading = true;
      this.errorMsg = '';

      try {
        // 1. Chiamata al Backend
        let response = await axios.post('/session', {
          name: this.username
        });

        // 2. Salva il token e l'username
        localStorage.setItem('token', response.data.identifier);
        localStorage.setItem('username', this.username);

        // 3. Vai alla Home
        this.$router.push('/');

      } catch (e) {
        this.errorMsg = e.response ? e.response.data.message : "Errore di connessione";
      } finally {
        this.loading = false;
      }
    }
  }
}
</script>
