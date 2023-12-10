export type Sort = 'date' | 'name' | 'size'

export interface UpdateDownload {
    url: string
    provider: string
    resumable: boolean
    progress: number
    expired: boolean
    downloadedChunks: number[]
    timeLeft: number
    speed: number
    status: string
}

export interface BatchDownload {
    ids: string[]
    payload: UpdateDownload[]
}