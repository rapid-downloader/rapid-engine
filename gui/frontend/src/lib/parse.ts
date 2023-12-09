import { useDateFormat } from "@vueuse/core"
import { Status } from "@/module/home/types"

export function parseSize(size: number): string {
    const KB = 1024
    const MB = KB * KB
    const GB = MB * KB

    if (size < KB) return (size / KB).toFixed(2) + " KB"
    if (size > KB && size < MB) return (size / KB).toFixed(2) + " KB"
    if (size > KB && size < GB) return (size / MB).toFixed(2) + " MB"
    if (size > GB) return (size / GB).toFixed(2) + " GB"

    return '0 KB'
}

export function parseElapsedTime(start: Date): string {
    const begin = new Date(start)
    const end = new Date()

    const elapsed = new Date(end.getTime() - begin.getTime())
    let s = elapsed.getSeconds() % 60
    let m = elapsed.getMinutes() % 60
    let h = ((elapsed.getMinutes() / 60) % 24)

    return `${h < 10 ? '0' + h.toFixed(0) : h.toFixed(0)}h : ${m < 10 ? '0' + m : m}m : ${s < 10 ? '0' + s : s}s`
}

export function parseDate(date: Date | string): string {
    return useDateFormat(date, 'MMMM DD, YYYY HH:mm').value
}

export function parseTimeleft(timeLeft: number): string {
    if (timeLeft === 0) {
        return '00:00:00'
    }

    const hours = Math.floor(timeLeft / 3600);
    const minutes = Math.floor((timeLeft % 3600) / 60);
    const remainingSeconds = timeLeft % 60;

    const h = hours.toFixed(0).padStart(2, '0');
    const m = minutes.toFixed(0).padStart(2, '0');
    const s = remainingSeconds.toFixed(0).padStart(2, '0');

    return `${h}:${m}:${s}`;
}

export function statusColor(status: Status): string {
    const records: Record<Status, string> = {
        'Completed': 'text-success',
        'Downloading': 'text-info',
        'Failed': 'text-destructive',
        'Stoped': 'text-warning',
        'Paused': '',
        'Queued': 'text-info'
    }

    return records[status]
}