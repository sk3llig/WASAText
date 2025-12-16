import axios from "axios";

const instance = axios.create({
	// L'indirizzo del backend Go.
	baseURL: "http://localhost:3000",
	timeout: 1000 * 5 // 5 secondi timeout
});

// Questo intercettore aggiunge automaticamente il TOKEN a ogni richiesta
instance.interceptors.request.use(
	(config) => {
		const token = localStorage.getItem("token");
		if (token) {
			config.headers["Authorization"] = `Bearer ${token}`;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

export default instance;

