import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import http from './plugins/http'
import ws from './plugins/ws'
import './assets/index.css'

createApp(App)
    .use(router)
    .use(http, { baseURL: 'http://localhost:9999' })
    .use(createPinia())
    .use(ws, { baseURL: 'ws://localhost:9999/ws/gui' })
    .mount('#app')
    