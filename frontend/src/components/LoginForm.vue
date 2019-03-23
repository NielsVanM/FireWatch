<template>
  <v-container fluid id="loginWrapper">
    <v-layout row justify-center align-center>
      <v-flex xs12 sm6 md3 light>
        <v-form v-on:submit.prevent="Authenticate" id="loginForm">
          <v-alert :value="error" type="warning" dismissible>{{error}}</v-alert>
          <v-layout justify-center align-center id="logoContainer">
            <img src="/images/logo.png" alt="FireWatch Logo">
            <h1>FireWatch</h1>
          </v-layout>
          <v-layout column align-center>
            <v-text-field autofocus v-model="username" label="Username" placeholder="Username"></v-text-field>
            <v-text-field
              v-model="password"
              type="password"
              label="Password"
              placeholder="Password"
            ></v-text-field>
            <v-btn type="submit" color="success">Log In</v-btn>
          </v-layout>
        </v-form>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import axios from "axios";
export default {
  name: "Login",
  data: () => {
    return {
      username: "",
      password: "",
      error: null
    };
  },
  methods: {
    Authenticate: function() {
      axios
        .post(this.$store.state.backendURL + "/api/v1/login/", {
          username: this.username,
          password: this.password
        })
        .then(resp => {
          if (resp.data.success) {
            this.error = null;

            // Set loggedin state
            var data = resp.data.data;
            this.$store.commit("loggedin", {token: data.token, user: data.user});

            // Redirect to next page if we tried to access a page
            if (this.$route.query.next != null) {
              this.$router.push(this.$route.query.next);
              return;
            }
            // Redirect to dash
            this.$router.push("/");
          } else {
            this.error = "Invalid credentials";
          }
        })
        .catch(err => {
          this.error =
            err;
        });
    }
  }
};
</script>

<style scoped>
#loginForm {
  border: 1px solid white;
  padding: 2em;
  border-radius: 10px;
  box-shadow: 0px 0px 10px;
  background-color: white;
}

#loginWrapper {
  width: 100vw;
  height: 100vh;
  background-image: url("/images/login_wallpaper.jpg");
  background-attachment: fixed;
  background-position: center center;
  background-repeat: no-repeat;
  background-size: auto;
}

#logoContainer {
  margin-bottom: 2em;
}

img {
  width: 40px;
  height: auto;
  margin-right: 1em;
}
</style>

