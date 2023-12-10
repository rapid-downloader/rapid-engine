import { getCurrentInstance } from "vue";

export default function() {
    const instance = getCurrentInstance()
    if(!instance || !instance.proxy) {
        throw new Error(`Websocket can only be called inside setup()`)
    }

    return instance.proxy.$socket
}