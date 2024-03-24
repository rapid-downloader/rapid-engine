import { http } from "@/plugins/http"
import { BatchDownload, Download, UpdateDownload } from "./types"

export default function Entries() {

    async function all(page: number = 1): Promise<Record<string, Download>> {
        try {
            const res = await http.get<Download[]>(`/entries?page=${page}`)
            if (res.status !== 200) {
                return {}
            }

            const data: Record<string, Download> = {}

            for (const entry of res.data) {
                data[entry.id] = entry
            }

            return data
        } catch (error) {
            console.error(error);
            return {}
        }
    }

    async function updateAll(entries: Record<string, Download>): Promise<boolean> {
        try {
            const ids: string[] = []
            const payload: UpdateDownload[] = []

            for (const [id, entry] of Object.entries(entries)) {
                ids.push(id)
                payload.push(entry)
            }

            const req: BatchDownload = {
                ids: ids,
                payload: payload
            }

            const res = await http.put('/entries', req)
            return res.status === 200
        } catch (error) {
            console.error(error);
            return false
        }
    }

    async function update(entry: Download) {
        try {
            const req: UpdateDownload = {
                url: entry.url,
                provider: entry.provider,
                resumable: entry.resumable,
                progress: entry.progress,
                expired: entry.expired,
                downloadedChunks: entry.downloadedChunks,
                timeLeft: entry.timeLeft,
                speed: entry.speed,
                status: entry.status,
            }

            console.log('update', req);
            

            const res = await http.put(`/entries/${entry.id}`, req)
            return res.status === 200
        } catch (error) {
            console.error(error);
            return false
        }
    }

    async function deleteEntry(id: string, fromDisk: boolean) {
        try {
            await http.delete(`/entries/${id}?fromDisk=${fromDisk}`)
        } catch (error) {
            console.error(error);
        }
    }

    return { all, updateAll, update, deleteEntry }
}