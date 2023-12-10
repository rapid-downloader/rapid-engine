import axios, { AxiosInstance } from "axios"
import { App } from "vue"

export interface HttpOption {
    baseURL: string,
}

declare module '@vue/runtime-core' {
    export interface ComponentCustomProperties {
        $http: AxiosInstance
    }
}

export default {
    install(app: App, option: HttpOption) {
        const http = axios.create({
            baseURL: option.baseURL || '/api',
            headers: {
                Authorization: `Bearer ${window.localStorage.getItem('jwt')}`
            }
        })

        app.provide('http', http)
        app.config.globalProperties.$http = http
    }
}