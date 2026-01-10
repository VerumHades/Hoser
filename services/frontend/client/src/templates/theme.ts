/**
 * Central theme configuration for consistent styling.
 * Includes colors, spacing, borders, typography, and variants for dark/light mode support.
 */

export const colors = {
    primary: {
        light: "bg-indigo-500 text-white hover:bg-indigo-600",
        dark: "bg-indigo-400 text-black hover:bg-indigo-500",
    },
    secondary: {
        light: "bg-slate-300 text-slate-900 hover:bg-slate-400",
        dark: "bg-slate-700 text-white hover:bg-slate-600",
    },
    success: {
        light: "bg-green-500 text-white hover:bg-green-600",
        dark: "bg-green-400 text-black hover:bg-green-500",
    },
    danger: {
        light: "bg-red-500 text-white hover:bg-red-600",
        dark: "bg-red-400 text-black hover:bg-red-500",
    },
    info: {
        light: "bg-indigo-100 text-indigo-900",
        dark: "bg-indigo-900 text-indigo-100",
    },
    border: {
        light: "border-slate-300",
        dark: "border-slate-700",
    },
    background: {
        light: "bg-white",
        dark: "bg-gray-900",
    },
    text: {
        light: "text-slate-900",
        dark: "text-slate-100",
    },
};

export const spacing = {
    xs: "p-1",
    sm: "p-2",
    md: "p-4",
    lg: "p-6",
    xl: "p-8",
};

export const borders = {
    rounded: "rounded",
    roundedLg: "rounded-lg",
    border: "border",
};

export const typography = {
    title: "text-lg font-semibold",
    subtitle: "text-sm text-slate-600 dark:text-slate-400",
    label: "font-medium",
    input: "px-3 py-1 border rounded w-full dark:bg-gray-800 dark:border-gray-600 dark:text-white",
};
