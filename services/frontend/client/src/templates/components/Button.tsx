import React from "react";
import { colors, borders, spacing } from "../theme";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: "primary" | "secondary" | "success" | "danger" | "info";
    size?: "sm" | "md" | "lg";
}

export const Button: React.FC<ButtonProps> = ({ variant = "primary", size = "md", className = "", ...props }) => {
    const sizeClass = {
        sm: spacing.sm,
        md: spacing.md,
        lg: spacing.lg,
    }[size];

    const colorClass = colors[variant].light + " dark:" + colors[variant].dark;

    return <button className={`${borders.rounded} ${sizeClass} ${colorClass} font-semibold ${className}`} {...props} />;
};
