<template>
	<nav class="navbar navbar-dark bg-dark-header px-4 justify-content-between sticky-top border-bottom-dark">

		<div class="d-flex align-items-center">
			<a class="navbar-brand d-flex align-items-center m-0 clean-brand" href="#">
				<img src="@/assets/logo.png" alt="WasaText" class="d-inline-block align-text-top navbar-logo me-2">
			</a>
		</div>

		<div class="d-flex align-items-center gap-3">

			<div
				class="user-pill d-flex align-items-center p-1 pe-3 rounded-pill"
				style="cursor: pointer;"
				title="Modifica Profilo"
				@click="$emit('open-profile')"
			>
				<img
					:src="myPhoto || 'https://placehold.co/50'"
					class="rounded-circle object-fit-cover avatar-border"
					width="38" height="38"
				>
				<span class="fw-bold ms-2 d-none d-md-block text-white">
            {{ myUsername }}
        </span>
			</div>

			<button
				class="btn btn-sm btn-logout-dark d-flex align-items-center"
				@click="$emit('logout')"
				title="Esci"
			>
				<span class="me-1">🚪</span> <span class="d-none d-sm-inline">Esci</span>
			</button>
		</div>
	</nav>
</template>

<script>
export default {
	props: ['myUsername', 'myPhoto'],
	emits: ['logout', 'open-profile']
}
</script>

<style scoped>
/* --- STILI NAVBAR --- */

.navbar {
	height: 80px;
	z-index: 1020;
}

.bg-dark-header {
	background-color: #1a1d21;
}

.border-bottom-dark {
	border-bottom: 1px solid #2c3036;
}

/* --- FIX LOGO: RIMUOVE IL QUADRATO GRIGIO --- */
.navbar-brand.clean-brand {
	background: transparent !important; /* Forza la trasparenza */
	box-shadow: none !important;        /* Rimuove ombre/aloni */
	border: none !important;            /* Rimuove bordi */
	padding: 0 !important;              /* Rimuove spazi extra */
	outline: none !important;           /* Rimuove il bordo di selezione */
}

/* Rimuove sfondo anche al passaggio del mouse o al click */
.navbar-brand.clean-brand:hover,
.navbar-brand.clean-brand:focus,
.navbar-brand.clean-brand:active {
	background-color: transparent !important;
}

.navbar-logo {
	height: 55px;
	width: auto;
	object-fit: contain;
	/* Assicurati che l'immagine non abbia bordi suoi */
	border: none;
	outline: none;
}
/* ------------------------------------------- */

.user-pill {
	transition: all 0.2s ease;
	border: 1px solid transparent;
	background-color: rgba(255,255,255, 0.05);
}
.user-pill:hover {
	background-color: rgba(255,255,255, 0.1);
	border-color: rgba(255,255,255, 0.2);
}

.avatar-border {
	border: 2px solid #3282b8;
	padding: 2px;
	background-color: #1a1d21;
}

.btn-logout-dark {
	color: #aebac6;
	border: 1px solid #2c3036;
	background-color: transparent;
	transition: all 0.2s;
	border-radius: 8px;
	padding: 0.4rem 0.8rem;
}
.btn-logout-dark:hover {
	background-color: rgba(220, 53, 69, 0.1);
	color: #ff6b6b;
	border-color: #ff6b6b;
}
.btn-logout-dark:focus {
	box-shadow: none; /* Rimuove alone blu sul bottone logout */
}

.object-fit-cover { object-fit: cover; }
</style>
