<template>
  <v-container>
    <v-layout wrap justify-space-around>
        <v-layout column justify-space-around>

      <v-form id="login-form">
        <v-text-field label="Username"></v-text-field>
      </v-form>
        </v-layout>
    </v-layout>
  </v-container>
</template>

<script>
export default {};
</script>

<style scoped>
#login-form {
  width: 300px;
  height: auto;
}
</style>

