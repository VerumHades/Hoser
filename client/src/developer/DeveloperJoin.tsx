import React, { useState } from "react";
import { MultiForm, FormPart, FormInput, useMultiformStateHandler } from "../components/MultiForm";
import RequireLogin from "../components/RequireLogin";

const DeveloperJoin: React.FC = () => {
  const [submitted, setSubmitted] = useState(false);
  const handler = useMultiformStateHandler()

  return (
    <div className="flex flex-col items-center justify-center min-h-screen px-4">
      {!submitted ? (
        <>
          <h1 className="text-2xl font-semibold text-gray-900 dark:text-white text-center mb-6">
            Join as a Developer
          </h1>

          <MultiForm onSubmit={() => {
            console.log(handler.data)
            setSubmitted(true);
          }}>
            <FormPart>
              <FormInput handler={handler} label="Full Name" name="fullName" required />
              <FormInput handler={handler} label="Email" name="email" type="email" required />
            </FormPart>

            <FormPart>
              <FormInput
                handler={handler}
                label="GitHub Profile"
                name="github"
                type="url"
                placeholder="https://github.com/yourname"
              />
            </FormPart>

            <FormPart>
              <FormInput handler={handler} label="Short Bio" name="bio" type="textarea" />
              <FormInput handler={handler} label="Experience (years)" name="experience" type="number" />
            </FormPart>
          </MultiForm>
        </>
      ) : (
        <div className="text-center">
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-3">
            Welcome aboard!
          </h2>
          <p className="text-gray-600 dark:text-gray-400">
            Your developer application has been received.
          </p>
        </div>
      )}
    </div>
  );
};

export default DeveloperJoin;
