import { Http } from "@/composable"
import { BatchDownload, UpdateDownload } from "./types"
import { client } from "wailsjs/go/models"

export default function Entries() {

    const http = Http()

    async function all(): Promise<Record<string, client.Download>> {
        try {
            const res = await http.get<client.Download[]>('/entries')
            if (res.status !== 200) {
                return {}
            }

            const data: Record<string, client.Download> = {}

            for (const entry of res.data) {
                data[entry.id] = entry
            }

            return data
        } catch (error) {
            return {}
        }
    }

    async function updateAll(entries: Record<string, client.Download>): Promise<boolean> {
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
            // TODO: add notification
            return false
        }
    }

    async function update(entry: client.Download) {
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

            const res = await http.put(`/entries/${entry.id}`, req)
            return res.status === 200
        } catch (error) {
            // TODO: add notification
            return false
        }
    }

    return { all, updateAll, update }
}