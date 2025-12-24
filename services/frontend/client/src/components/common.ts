export interface HasChildren {
    children?: React.ReactNode
}

export interface Clickable {
    onClick?: () => void
}

export interface HasClassname {
    className?: string
}

export function formatBytes(byteCount?: number): string {
    if (byteCount === undefined || byteCount === null) {
        return "—"
    }

    if (byteCount < 1024) {
        return `${byteCount} B`
    }

    const units = ["KB", "MB", "GB", "TB", "PB"]
    let unitIndex = -1
    let value = byteCount

    do {
        value = value / 1024
        unitIndex++
    } while (value >= 1024 && unitIndex < units.length - 1)

    return `${value.toFixed(1)} ${units[unitIndex]}`
}