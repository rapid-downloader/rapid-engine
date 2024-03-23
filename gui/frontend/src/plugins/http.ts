import { toast } from "@/components/ui/sonner"
import axios, { AxiosInstance, isAxiosError } from "axios"
import Refresh from "@/components/ui/icon/Refresh.vue"
import { App } from "vue"

export interface HttpOption {
    baseURL: string,
}

declare module '@vue/runtime-core' {
    export interface ComponentCustomProperties {
        $http: AxiosInstance
    }
}

export const http = axios.create({
    baseURL: 'http://localhost:8888',
})

export default {
    install(app: App) {
        const onError = (err: any) => {
            if (isAxiosError(err)) {
                const res = err.response
                if (res) {
                    toast('Error', {
                        description: res.data.message,
                        type: 'destructive',
                        dismissible: true
                    })

                    return
                }

                if (err.code === 'ERR_NETWORK') {
                    toast('Error', {
                        description: 'Rapid engine is down. Please restart the engine',
                        type: 'destructive',
                        dismissible: true,
                        action: {
                            label: 'Restart',
                            icon: Refresh,
                            onClick: () => {
                                // TODO: restart engine

                            }
                        }
                    })

                    return
                }
                
                
                toast('Error', {
                    description: err.message,
                    type: 'destructive',
                    dismissible: true
                })

                return
            }

            toast('Error', {
                description: `${err}`,
                type: 'destructive',
                dismissible: true
            })
        }

        http.interceptors.request.use(undefined, onError)
        http.interceptors.response.use(undefined, onError)

        app.provide('http', http)
        app.config.globalProperties.$http = http
    }
}