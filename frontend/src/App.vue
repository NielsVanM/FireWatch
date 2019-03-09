<template>
  <v-app id="inspire" dark>
    <!-- Drawer -->
    <v-navigation-drawer clipped app :mini-variant="drawer">
      <v-list>
        <v-list-tile v-for="item in items" :key="" active-class="white--text">
          <v-list-tile-action>
            <v-icon medium>{{item.icon}}</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>{{item.title}}</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <!-- Toolbar -->
    <v-toolbar app fixed clipped-left>
      <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <img src="/images/logo.png" class="fw-logo">
      <v-toolbar-title>FireWatch</v-toolbar-title>
    </v-toolbar>
    <!-- Page Content -->
    <v-content>
      <v-container fluid>
        <router-view></router-view>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
export default {
  data: () => ({
    drawer: null,
    items: [
      {title: "Dashboard", icon: "dashboard", to: "/"},
      {title: "Devices", icon: "memory", to: "/device/"},
      {title: "Settings", icon: "settings", to: "/settings/"},
    ]
  }),
  props: {
    source: String
  }
};
</script>

<style scoped>
.fw-logo {
  height: 2em;
  width: auto;
  margin-left: 1em;
  margin-right: -.5em;
}
</style>
