import { http } from "@/plugins/http";
import { Setting } from "./types";

export default function() {

    async function get(): Promise<Setting | undefined> {
        try {
            const res = await http.get('/settings')
            return res.data
        } catch (error) {
            console.error(error);
        }
    }

    async function update(payload: Setting): Promise<void> {
        try {
            const res = await http.put('/settings', payload)
            return res.data
        } catch (error) {
            console.error(error);
        }
    }

    return {
        get,
        update
    }
}