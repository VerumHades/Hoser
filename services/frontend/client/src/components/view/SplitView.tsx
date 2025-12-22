import { type ReactNode, useState, useRef, useEffect } from "react"

interface SplitViewProps {
    leftPane: ReactNode
    rightPane: ReactNode
    minLeftWidth?: number
    minRightWidth?: number
}

export default function SplitView({
    leftPane,
    rightPane,
    minLeftWidth = 240,
    minRightWidth = 240
}: SplitViewProps) {
    const containerRef = useRef<HTMLDivElement>(null)
    const [dividerDragging, setDividerDragging] = useState(false)
    const [leftWidth, setLeftWidth] = useState<number | null>(480)
    const [showRightOnMobile, setShowRightOnMobile] = useState(false)

    // Drag handling
    useEffect(() => {
        const handleMouseMove = (e: MouseEvent) => {
            if (!dividerDragging || !containerRef.current) return
            const rect = containerRef.current.getBoundingClientRect()
            let newWidth = e.clientX - rect.left
            if (newWidth < minLeftWidth) newWidth = minLeftWidth
            if (rect.width - newWidth < minRightWidth) newWidth = rect.width - minRightWidth
            setLeftWidth(newWidth)
        }

        const handleMouseUp = () => setDividerDragging(false)

        window.addEventListener("mousemove", handleMouseMove)
        window.addEventListener("mouseup", handleMouseUp)
        return () => {
            window.removeEventListener("mousemove", handleMouseMove)
            window.removeEventListener("mouseup", handleMouseUp)
        }
    }, [dividerDragging, minLeftWidth, minRightWidth])

    return (
        <div ref={containerRef} className="flex flex-col md:flex-row h-full w-full">
            {/* Left pane */}
            <div
                className={`
                    flex flex-col border-r border-gray-200 p-3
                    ${showRightOnMobile ? "hidden md:flex" : "flex"}
                `}
                style={{ width: leftWidth ? `${leftWidth}px` : "auto", minWidth: minLeftWidth }}
            >
                {leftPane}
            </div>

            {/* Divider */}
            <div
                className="hidden md:flex w-1 cursor-col-resize bg-gray-200 hover:bg-gray-300 transition-colors"
                onMouseDown={() => setDividerDragging(true)}
            />

            {/* Right pane */}
            <div
                className={`
                    flex-1 overflow-auto
                    ${!showRightOnMobile ? "hidden md:flex" : "flex w-full"}
                `}
            >
                <div className="flex flex-col h-full">
                    {/* Mobile back button */}
                    <div className="md:hidden flex p-2 bg-gray-50 border-b border-gray-200">
                        <button
                            className="text-blue-600 font-semibold"
                            onClick={() => setShowRightOnMobile(false)}
                        >
                            ← Back
                        </button>
                    </div>
                    {rightPane}
                </div>
            </div>
        </div>
    )
}
