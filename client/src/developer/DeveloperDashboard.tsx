// src/components/DashboardApp.tsx
import Dashboard, { Page } from "../components/Dashboard";
import { BarChart3, Images } from "lucide-react";
import RequireLogin from "../components/RequireLogin";

export default function DashboardApp() {
  return (
    <RequireLogin>
      <RequireLogin
        condition={(user) => user.IsDeveloper}
        redirect="/developer/join"
      >
        <Dashboard>
          <Page name="Listings" icon={<Images />}>
            <div>Listing content goes here</div>
          </Page>
          <Page name="Earnings" icon={<BarChart3 />}>
            <div>Earnings content goes here</div>
          </Page>
        </Dashboard>
      </RequireLogin>
    </RequireLogin>
  );
}
