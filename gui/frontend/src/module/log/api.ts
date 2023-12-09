import Http from "@/composable/http"

export default function() {

    const http = Http()

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