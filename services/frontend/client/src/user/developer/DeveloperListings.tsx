import { useState } from "react"
import DeveloperListingDisplay from "./DeveloperListing"
import DeveloperListingsCollection from "./DeveloperListingsCollection"

import { type DeveloperListing } from "../../backend"
import { useNavigate } from "react-router-dom"

export default function DeveloperListings() {
    const navigate = useNavigate()
    return <DeveloperListingsCollection onSelect={(listing) => navigate("/dashboard/developer/listing/" + listing.id)}/>
}
