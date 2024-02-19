import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import http from './plugins/http'
import './assets/index.css'


createApp(App)
    .use(router)
    .use(http, { baseURL: 'http://localhost:8888' })
    .use(createPinia())
    .mount('#app')
    