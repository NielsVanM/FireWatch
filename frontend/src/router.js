import Vue from "vue"
import VueRouter from "vue-router"
import Axios from "axios"

import LoginView from "@/views/Login.vue"
import Base from "@/views/Base.vue"
import Dashboard from "@/views/Dashboard.vue"
import DeviceOverView from "@/views/Device.vue"

Vue.use(VueRouter)

var router = new VueRouter({
    mode: "history",
    routes: [
        {
            path: "/",
            name: "base",
            components: {
                root: Base
            },
            meta: {
                protected: true,
            },
            children: [
                {
                    path: "",
                    name: "dashboard",
                    components: {int: Dashboard},
                    meta: {
                        protected: true
                    }
                },{
                    path: "device/",
                    name: "device",
                    commponents: {int: DeviceOverView},
                    meta: {
                        protected: true
                    }
                },{
                    path: "settings/",
                    name: "settings",
                    components: {int: DeviceOverView},
                    meta: {
                        protected: true
                    }
                }
            ]
        },{
            path: "/login/",
            name: "login",
            components: {root: LoginView},
            meta: {
                protected: false
            }
        }
    ]
})


// Auth middleware
router.beforeEach((to, from, next) => {
    // Check if the route is protected
    if (!to.meta.protected) {
        next()
    }
    if (to.name == "login") {
        next()
    }

    // Verify token against backend
    var token = localStorage.getItem("token")
    if (token) {
        Axios
            .post("http://localhost:8000/api/v1/verify/", { "token": token })
            .then((resp) => {
                if (resp.data.success) {
                    // Set as auth header
                    Vue.prototype.$http.defaults.headers.common['Authorization'] = token
                    next()
                } else {
                    // If the token is invalid delete all local data
                    if (resp.data.status_code == "invalid_token") {
                        localStorage.removeItem("token")
                        localStorage.removeItem("user")
                    }
                    next('/login/?next=' + to.path)
                }
            }
            )
            .catch(() => {
                next('/login/?next=' + to.path)
            })
    } else {
        next('/login/?next=' + to.path)
    }
})

export default router