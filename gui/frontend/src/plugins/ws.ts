import { App } from "vue"

export interface WebSocketOption {
    baseURL: string,
    onopen?: Function
    onclose?: Function
    onerror?: Function
}

type Callback = (event: Event) => void
type Message = (data: any) => void

interface Socket {
    socket?: WebSocket
    connect: (callback?: Callback) => void
    disconnect: (callback?: Callback) => void
    send: (message: string | ArrayBufferLike | Blob | ArrayBufferView) => void
    onmessage: Message
}

declare module '@vue/runtime-core' {
    export interface ComponentCustomProperties {
        $socket: Socket,
    }
}

export default {
    install(app: App, option: WebSocketOption) {
        let ws: WebSocket | undefined = undefined

        function _connect(callback?: Callback) {
            if (ws?.readyState === ws?.OPEN && ws) return

            ws = new WebSocket(option.baseURL)
            ws.onopen = event => {
                console.log('web socket connected')
                if (callback) callback(event)
            };
        }

        function connect(callback?: Callback) {
            _connect(callback)
        }

        function close() {
            try {
                if (ws?.readyState !== ws?.CLOSED && ws) ws.close()
                console.log('web socket closed');
            } catch (error) {
                if (error instanceof Error) {
                    console.log('Error closing web socket:', error.message);
                }
            }
        }

        function _disconnect(callback?: Callback) {
            if (ws) {
                ws.onclose = (event) => {
                    if (callback) callback(event)
                }
            }

            close()
        }

        function disconnect(callback?: Callback) {
            _disconnect(callback)
        }

        function send(message: string | ArrayBufferLike | Blob | ArrayBufferView) {
            if (!ws) return
            ws.send(message)
        }

        function onmessage<T>(fn: Message) {
            if (ws) {
                ws.onmessage = event => {
                    const data = JSON.parse(event.data)
                    if (data) fn(data)
                }
            }
        }

        const socket: Socket = {
            socket: ws,
            connect: connect,
            disconnect: disconnect,
            onmessage: onmessage,
            send: send
        }

        app.provide('ws', socket)
        app.config.globalProperties.$socket = socket
    }
}