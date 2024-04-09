import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
    history: createWebHashHistory('/'),
    routes: [
        { path: '/', component: () => import("./module/home/Home.vue") },
        { path: '/log', component: () => import("./module/log/Log.vue") },
        { path: '/download/:id', component: () => import("./module/download/Download.vue") },
        { path: '/setting', component: () => import("./module/setting/Setting.vue") },
    ]
})

export default router