import { createRouter, createWebHashHistory } from "vue-router";
import Home from "./module/home/Home.vue";
import Log from "./module/log/Log.vue";
import Download from "./module/download/Download.vue";

const router = createRouter({
    history: createWebHashHistory('/'),
    routes: [
        { path: '/', component: Home },
        { path: '/log', component: Log },
        { path: '/download/:id', component: Download },
    ]
})

export default router