import Vue from 'vue'
import './plugins/vuetify'
// import 'vuetify/dist/vuetify.min.css'
import Axios from "axios"

import App from './App.vue'
import router from "./router.js"
import store from "./store.js"

// Set up AXIOS
Vue.prototype.$http = Axios;
Vue.prototype.$http.defaults.headers.common['Authorization'] = localStorage.getItem("token")

Vue.config.productionTip = false

new Vue({
  router: router,
  store: store,
  render: h => h(App),
}).$mount('#app')
