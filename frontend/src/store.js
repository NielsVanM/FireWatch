import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
    state: {
        authenticated: false,
        user: {},
        token: "",
        backendURL: "http://localhost:8000"
    },
    mutations: {
        loggedin (state, load) {
            // Set local storage
            localStorage.setItem("token", load.token)
            localStorage.setItem("user", JSON.stringify(load.user))

            // Set state
            state.token = load.token
            state.user = load.user
            state.authenticated = true
        }
    },
    getters: {
        isAuthenticated () {
            return localStorage.getItem("token") != "";
        },
        getUser () {
            return JSON.parse(localStorage.getItem("user"))
        },
        getToken () {
            return localStorage.getItem("token")
        }
    }
})