import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { oneLight } from "react-syntax-highlighter/dist/cjs/styles/prism";

interface GitHubMarkdownProps {
    /** Raw markdown content string */
    content: string;
}

/**
 * A full-featured Markdown viewer that emulates the GitHub display style.
 * Includes GFM support and high-fidelity syntax highlighting.
 */
export default function MarkdownViewer({ content }: GitHubMarkdownProps) {
    return (
        <article className="prose prose-slate max-w-none 
            prose-headings:border-b prose-headings:pb-2 prose-headings:border-slate-200
            prose-pre:bg-transparent prose-pre:p-0
            prose-table:border prose-table:rounded-xl
            prose-th:bg-slate-50 prose-th:p-3
            prose-td:p-3 prose-td:border-t">
            
            <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                components={{
                    // Custom renderer for code blocks
                    code({ node, inline, className, children, ...props }: any) {
                        const match = /language-(\w+)/.exec(className || "");
                        
                        return !inline && match ? (
                            <div className="rounded-md overflow-hidden border border-slate-200 my-4">
                                <div className="bg-slate-50 px-4 py-2 text-xs font-mono text-slate-500 border-b border-slate-200 flex justify-between">
                                    <span>{match[1].toUpperCase()}</span>
                                </div>
                                <SyntaxHighlighter
                                    style={oneLight}
                                    language={match[1]}
                                    PreTag="div"
                                    customStyle={{
                                        margin: 0,
                                        padding: "1rem",
                                        fontSize: "0.875rem",
                                        backgroundColor: "#f8fafc",
                                    }}
                                    {...props}
                                >
                                    {String(children).replace(/\n$/, "")}
                                </SyntaxHighlighter>
                            </div>
                        ) : (
                            <code className="bg-slate-100 px-1.5 py-0.5 rounded text-pink-600 font-mono text-sm" {...props}>
                                {children}
                            </code>
                        );
                    },
                }}
            >
                {content}
            </ReactMarkdown>
        </article>
    );
}