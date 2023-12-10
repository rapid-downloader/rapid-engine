import { getCurrentInstance } from "vue";

export default function Http() {
    const instance = getCurrentInstance()
    if(!instance || !instance.proxy) {
        throw new Error(`Http can only be called inside setup()`)
    }

    return instance.proxy.$http
}