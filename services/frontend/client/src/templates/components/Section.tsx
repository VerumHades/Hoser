import React from "react";
import { colors, typography, spacing, borders } from "../theme";

interface SectionProps {
    title: string;
    description?: string;
    className?: string;
    children: React.ReactNode;
}

/**
 * Section is a reusable layout block for grouping related form elements or content.
 * Fully theme-driven with dark/light support, spacing, and typography consistency.
 */
export const Section: React.FC<SectionProps> = ({ title, description, className = "", children }) => {
    return (
        <div className={`mb-8 border-b ${colors.border.light} dark:${colors.border.dark} pb-6 ${className}`}>
            <h2 className={`${typography.title} mb-1 ${colors.text.light} dark:${colors.text.dark}`}>
                {title}
            </h2>
            {description && (
                <p className={`${typography.subtitle} mb-3`}>
                    {description}
                </p>
            )}
            {children}
        </div>
    );
};

export default Section;
