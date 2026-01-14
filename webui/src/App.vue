<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
	<RouterView />
</template>

<style>
/* Reset globale */
html, body {
	height: 100%;
	margin: 0;
	padding: 0;
	overflow: hidden; /* Blocca lo scroll della pagina principale */

	/* --- IL TRUCCO DELLO ZOOM --- */
	zoom: 1; /* Scala tutto al 130% */

	/* Colore di sfondo scuro di sicurezza
	   (Se lo zoom crea spazi vuoti ai bordi, saranno scuri e non bianchi) */
	background-color: #1a1d21;

	/* Migliora la nitidezza del testo quando è ingrandito */
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
}

#app {
	height: 100%;
	width: 100%;
}
</style>
