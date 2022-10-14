import {createRouter, createWebHistory} from "vue-router";

import Home from "@/pages/Home.vue";
import Login from "@/pages/Login.vue";
import Register from "@/pages/Register.vue";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {path: "/", component: Home},
        {path: "/login", component: Login},
        {path: "/register", component: Register},
    ],
});

export default router;
