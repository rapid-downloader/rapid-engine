import { http } from "@/plugins/http"
import { Download, Request } from "../home/types"

export interface Downloader {
    fetch(req: Request): Promise<Download | undefined>
    download(id: string): void
    stop(id: string): void
    pause(id: string): void
    resume(id: string): void
    restart(id: string): void
}

export function Downloader(): Downloader {

    async function fetch(req: Request): Promise<Download | undefined> {
        try {
            const res = await http.post<Download>('/fetch', req)
            if (res.status === 200) return res.data
        } catch (error) {
            console.error(error);
        }
    }

    async function download(id: string) {
        try {
            await http.get(`/gui/download/${id}`)
        } catch (error) {
            console.error(error);
        }
    }

    async function stop(id: string): Promise<void> {
        try {
            await http.put(`/stop/${id}`)
        } catch (error) {
            console.error(error);
        }
    }

    async function pause(id: string): Promise<void> {
        try {
            await http.put(`/pause/${id}`)
        } catch (error) {
            console.error(error);
        }
    } 

    async function resume(id: string): Promise<void> {
        try {
            await http.put(`/gui/resume/${id}`)
        } catch (error) {
            console.error(error);
        }
    } 

    async function restart(id: string): Promise<void> {
        try {
            await http.put(`/gui/restart/${id}`)
        } catch (error) {
            console.error(error);
        }
    } 

    return { fetch, download, stop, pause, resume, restart }
}