import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import http from './plugins/http'
import './assets/index.css'

createApp(App)
    .use(router)
    .use(http)
    .use(createPinia())
    .mount('#app')
    