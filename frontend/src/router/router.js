import Vue from "vue"
import VueRouter from "vue-router"

import Dashboard from "@/components/dashboard.vue"

Vue.use(VueRouter)

export default new VueRouter({
    routes: [
        {path: "/", component: Dashboard}
    ]
})