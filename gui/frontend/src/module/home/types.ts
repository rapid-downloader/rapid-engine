export type Type = 'Document' | 'Audio' | 'Video' | 'Image' | 'Compressed' | 'Other'

export type Status = 'Completed' | 'Stoped' | 'Downloading' | 'Queued' | 'Failed' | 'Paused'

export type Sort = 'date' | 'name' | 'size'

export interface Download {
    id: string
    name: string
    url: string
    provider: string
    size: number
    type: Type
    chunklen: number
    resumable: boolean
    progress: number
    expired: boolean
    downloadedChunks: number[]
    timeLeft: number
    speed: number
    status: Status
    date: Date
}

export interface UpdateDownload {
    url: string
    provider: string
    resumable: boolean
    progress: number
    expired: boolean
    downloadedChunks: number[]
    timeLeft: number
    speed: number
    status: Status
}

export interface BatchDownload {
    ids: string[]
    payload: UpdateDownload[]
}

export interface Request {
    url: string
    provider: string
    client: string
    mimeType?: string
    userAgent?: string
    cookies?: Cookie
}

export interface Cookie {
    name: string
    value: string
    path: string
    domain: string
    expirationDate: Date
    secure: boolean
    httpOnly: boolean
    sameSite: string
}