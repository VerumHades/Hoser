import React, { useRef, useState, type ReactElement } from "react";
import { motion, AnimatePresence } from "framer-motion";

interface MultiFormProps {
  children: React.ReactNode;
  onSubmit?: () => void;
}

interface FormPartProps {
  children: React.ReactNode;
}

export interface MultiformStateHandler {
  data: Record<string, any>
}

export function useMultiformStateHandler(): MultiformStateHandler{
  return {data: {}}
}

interface FormInputProps extends React.InputHTMLAttributes<HTMLInputElement | HTMLTextAreaElement> {
  label: string;
  name: string;
  handler: MultiformStateHandler
  type?: string;
}

// --- FORM INPUT ---
export const FormInput: React.FC<FormInputProps> = ({ label, name, handler, type = "text", ...rest }) => {
  const sharedClass =
    "w-full px-3 py-2 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100 " +
    "border border-gray-300 dark:border-gray-700 focus:ring-2 focus:ring-blue-500 focus:outline-none";

  const set_state_value = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    let input = e.target as HTMLInputElement
    handler.data[input.name] = input.value
  }
  
  return (
    <div className="mb-4">
      <label htmlFor={name} className="block mb-2 text-gray-800 dark:text-gray-200 font-medium">
        {label}
      </label>
      {type === "textarea" ? (
        <textarea id={name} name={name} {...(rest as any)} rows={4} onChange={set_state_value} className={sharedClass} />
      ) : (
        <input id={name} name={name} type={type} {...rest} onChange={set_state_value} className={sharedClass} />
      )}
    </div>
  );
};

// --- FORM PART ---
export const FormPart: React.FC<FormPartProps> = ({ children }) => <>{children}</>;
// --- MULTI FORM ---
export const MultiForm: React.FC<MultiFormProps> = ({ children, onSubmit }) => {
  const steps = React.Children.toArray(children);
  const [step, setStep] = useState(0);

  const next = () => setStep((s) => Math.min(s + 1, steps.length - 1));
  const prev = () => setStep((s) => Math.max(s - 1, 0));

  const handleSubmit = (formData: Record<string, any>) => {
    if (step === steps.length - 1) {
      onSubmit?.();
    } else {
      next();
    }
  };

  return (
    <form
      action={handleSubmit}
      className="w-full space-y-4 max-w-lg"
    >
      {/* Progress indicator */}
      <div className="flex justify-center mb-4">
        <div className="flex space-x-2">
          {steps.map((_, i) => (
            <div
              key={i}
              className={`w-3 h-3 rounded-full transition-all duration-300 ${
                i === step
                  ? "bg-blue-600 dark:bg-blue-400 scale-125"
                  : "bg-gray-300 dark:bg-gray-700"
              }`}
            />
          ))}
        </div>
      </div>

      {/* Active step */}
      <AnimatePresence mode="wait">
        <motion.div
          key={step}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          transition={{ duration: 0.3 }}
          className="space-y-4"
        >
          {steps[step]}
        </motion.div>
      </AnimatePresence>

      {/* Buttons */}
      <div className="flex justify-between pt-4">
        {step > 0 && (
          <button
            type="button"
            onClick={prev}
            className="px-4 py-2 text-gray-600 dark:text-gray-300 hover:underline"
          >
            Back
          </button>
        )}
        <button
          type="submit"
          className="ml-auto bg-blue-600 hover:bg-blue-700 text-white font-medium px-5 py-2 rounded-lg transition-colors"
        >
          {step === steps.length - 1 ? "Submit" : "Next"}
        </button>
      </div>
    </form>
  );
};
