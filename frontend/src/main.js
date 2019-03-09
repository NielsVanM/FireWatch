import Vue from 'vue'
import './plugins/vuetify'
import 'vuetify/dist/vuetify.min.css' 
import App from './App.vue'
import router from "./router/router.js"

Vue.config.productionTip = false

new Vue({
  router: router,
  mode: "abstract",
  render: h => h(App),
}).$mount('#app')
