import { http } from "@/plugins/http"

export default function() {

    async function get(date: string): Promise<string[]> {
        try {
            const res = await http.get(`/logs/${date}`)
            return res.data
        } catch (error) {
            return []
        }
    }

    return { get }
}