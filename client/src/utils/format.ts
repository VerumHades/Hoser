
export function formatBytes(bytes?: number): string {
    if (bytes == null) return "0 B";
    const thresh = 1024;
    if (bytes < thresh) return bytes + " B";
    const units = ["KB", "MB", "GB", "TB"];
    let u = -1;
    let b = bytes;
    do {
        b /= thresh;
        u++;
    } while (b >= thresh && u < units.length - 1);
    return b.toFixed(1) + " " + units[u];
}