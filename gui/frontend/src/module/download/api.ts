import { Http } from "@/composable"
import { isAxiosError } from "axios"
import { Download, Request } from "../home/types"

export interface Downloader {
    fetch(req: Request): Promise<Download | undefined>
    download(id: string): void
}

export function Downloader(): Downloader {

    const http = Http()

    async function fetch(req: Request): Promise<Download | undefined> {
        try {
            const res = await http.post<Download>('/fetch', req)
            if (res.status === 200) {
                return res.data
            }
            
        } catch (error) {
            if (isAxiosError(error)) {
                // TODO: add notification
                console.log(error.response?.data.message);
            }
        }
    }

    async function download(id: string) {
        try {
            await http.get(`/gui/download/${id}`)
            
        } catch (error) {
            if (isAxiosError(error)) {
                // TODO: add notification
                console.log(error.response?.data.message);
            }

            return false
        }
    }

    return { fetch, download }
}