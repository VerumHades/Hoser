// src/components/DashboardApp.tsx
import Dashboard, { Page } from "../components/Dashboard";
import { BarChart3, Images } from "lucide-react";
import RequireLogin from "../components/RequireLogin";
import DeveloperListings from "./DeveloperListings";

export default function DashboardApp() {
  return (
    <RequireLogin>
      <RequireLogin
        condition={(user) => user.isDeveloper}
        redirect="/developer/join"
      >
        <div className="w-full h-full">
        <Dashboard>
          <Page name="Listings" icon={<Images />}>
            <DeveloperListings></DeveloperListings>
          </Page>
          <Page name="Earnings" icon={<BarChart3 />}>
            <div>Earnings content goes here</div>
          </Page>
        </Dashboard>
        </div>
      </RequireLogin>
    </RequireLogin>
  );
}
