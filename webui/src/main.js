import { createApp } from 'vue' 
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';

// Componenti Globali (utili per mostrare errori o caricamenti ovunque)
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'

// Stili CSS
import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)

// Configurazione Globale

app.config.globalProperties.$axios = axios;

// Registrazione Componenti
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);

app.use(router)
app.mount('#app')
