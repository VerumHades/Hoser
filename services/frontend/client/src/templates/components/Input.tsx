import React from "react";
import { typography, borders, spacing } from "../theme";

interface TextInputProps extends React.InputHTMLAttributes<HTMLInputElement> {}

export const TextInput: React.FC<TextInputProps> = (props) => {
    return <input className={`${typography.input} ${spacing.sm}`} {...props} />;
};
