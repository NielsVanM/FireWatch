<template>
  <v-app :dark="this.$store.getters.isAuthenticated && this.$route.name != 'login'">
    <router-view name="root" />
  </v-app>
</template>
