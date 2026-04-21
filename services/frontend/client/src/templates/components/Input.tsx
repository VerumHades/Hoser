import React from "react";
import { typography, spacing } from "../theme";

type TextInputProps = React.InputHTMLAttributes<HTMLInputElement>

export const TextInput: React.FC<TextInputProps> = (props) => {
    return <input className={`${typography.input} ${spacing.sm}`} {...props} />;
};
