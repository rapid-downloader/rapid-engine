import { Http } from "@/composable"
import { isAxiosError } from "axios"
import { client } from '@/../wailsjs/go/models'
import { Fetch, Download } from '@/../wailsjs/go/main/App'

export interface Downloader {
    fetch(req: client.Request): Promise<client.Download | undefined>
    download(id: string): Promise<boolean>
}

export function HttpDownloader(): Downloader {

    const http = Http()

    async function fetch(req: client.Request): Promise<client.Download | undefined> {
        try {
            const res = await http.post<client.Download>('/fetch', req)
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

    async function download(id: string): Promise<boolean> {
        try {
            const res = await http.get(`/gui/download/${id}`)
            return res.status === 200
            
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

export function BindingDownloader(): Downloader {

    async function fetch(req: client.Request): Promise<client.Download | undefined> {
        try {
            return await Fetch(req)
        } catch (error) {
            console.error(error)
        }
    }


    async function download(id: string): Promise<boolean> {
        try {
            await Download(id)
            return true
        } catch (error) {
            console.log(error);
            return false
        }
    }

    return {
        fetch,
        download
    }
}

export default function useDownloader(provider: 'http' | 'binding'): Downloader {
    return provider === 'http' ? HttpDownloader() : BindingDownloader()
}