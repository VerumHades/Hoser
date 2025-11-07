// src/components/AccountPage.tsx
import React, {  useState } from "react";
import EditableText from "../components/EditableText";
import { backend_request, ListingAccessModes } from "../backend";
import { Delete } from "lucide-react";
import PromptButton from "../components/PromptButton";
import SelectBox from "../components/SelectBox";

export interface Listing {
  ID: string;
  Title: string;
  Description: string;
  Author: string;
  AccessMode: number
}

interface ListingDisplayProps {
  listing: Listing,
  onShouldClose?: () => void
}

export default function DeveloperListingDisplay({ listing, onShouldClose }: ListingDisplayProps) {
  const [currentText, setText] = useState<string>("");
  const deleteKeyword = "Delete"
  const isDeleteLocked = currentText != deleteKeyword

  const handleChange = (event: React.KeyboardEvent<HTMLInputElement>) => {
      setText((event.target as HTMLInputElement).value);
  };

  const deleteThisListing = async () => {
    if((await backend_request("/developer/listing", "DELETE", {id: listing.ID})).status == 200)
      onShouldClose?.()
  }

  const onSelectedAccessMode = (mode: string, oldMode: string) => {
    backend_request("/developer/listing", "PUT", { id: listing.ID, accessMode: parseInt(mode) })
  }   

  return (
    <div className="w-full h-full p-6  shadow-md rounded-lg">
      <div className="flex flex-row justify-between mb-10">
        <SelectBox options={ListingAccessModes} onSelected={onSelectedAccessMode} defaultValue={""+listing.AccessMode}></SelectBox>
        <PromptButton
          className="bg-red-500 hover:bg-red-600 text-white font-semibold py-2 px-4 rounded transition-colors dark:bg-red-600 dark:hover:bg-red-700 flex flex-row gap-2"
          submitClassName={`
            flex flex-row gap-2 font-semibold py-2 px-4 rounded transition-colors
            ${isDeleteLocked
              ? "bg-red-300 text-gray-200 cursor-not-allowed dark:bg-red-700 dark:text-gray-400"
              : "bg-red-500 hover:bg-red-600 text-white dark:bg-red-600 dark:hover:bg-red-700"
            }
          `}
          promptContents={
            <>
            <p>This operation will delete this listing. If you are sure you want to do this type "Delete" as confirmation.</p>
            <input onKeyUp={handleChange} type="text"
            className="w-full p-3 my-10 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500
          dark:bg-gray-700 dark:border-gray-600 dark:text-gray-100 dark:placeholder-gray-400 dark:focus:ring-red-400 dark:focus:border-red-400"></input>
            </>
          }
          isSubmitLocked={() => isDeleteLocked}
          onCancel={() => setText("")}
          onSubmit={() => deleteThisListing()}
        >
          Delete <Delete></Delete>
        </PromptButton>
      </div>
      <h1 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">
        <EditableText 
          text={listing.Title} 
          onChange={(name: string) => { backend_request("/developer/listing", "PUT", { id: listing.ID, title: name }) }}>

        </EditableText>
      </h1>
      <p className="mb-2 text-gray-700 dark:text-gray-300">
        <EditableText 
          text={listing.Description} 
          onChange={(name: string) => { backend_request("/developer/listing", "PUT", { id: listing.ID, description: name }) }}>

        </EditableText>
      </p>
    </div>
  );
};
