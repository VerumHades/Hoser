// src/components/AccountPage.tsx
import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getUser, type User } from "../user";
import backend_constants from "../backend_constants";
import EditableText from "../components/EditableText";


export interface Listing {
    ID: string | number;
    Title: string;
    Description: string;
    icon?: string;
    image?: string;
    Author: string;
}

interface ListingDisplayProps {
    listing: Listing
}

export default function DeveloperListingDisplay({listing}: ListingDisplayProps){
  return (
    <div className="w-full h-full p-6  shadow-md rounded-lg">
      <h1 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">
        <EditableText text={listing.Title} onChange={() => {}}></EditableText>
      </h1>
      <p className="mb-2 text-gray-700 dark:text-gray-300">
        <EditableText text={listing.Description} onChange={() => {}}></EditableText>
      </p>
    </div>
  );
};
