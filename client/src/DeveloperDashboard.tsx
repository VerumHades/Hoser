// src/components/DashboardApp.tsx
import Dashboard, { Page } from "./Dashboard";
import { BarChart3, Images } from "lucide-react";

export default function DashboardApp() {
  return (
    <Dashboard>
        <Page name="Listings" icon={<Images />}>
            <div>Listing content goes here</div>
        </Page>
        <Page name="Earnings" icon={<BarChart3 />}>
            <div>Earnings content goes here</div>
        </Page>
    </Dashboard>
  );
}
