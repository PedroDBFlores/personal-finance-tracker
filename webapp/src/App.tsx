import { createRouter } from "@solidjs/router";
import type { RouteDefinition } from "@solidjs/router";
import { RouterLink } from "./components/RouterLink";
import AccountsPage from "./pages/AccountsPage";
import NewAccountPage from "./pages/NewAccountPage";
import AccountDetailPage from "./pages/AccountDetailPage";
import NewMovementPage from "./pages/NewMovementPage";
import NewInvestmentPage from "./pages/NewInvestmentPage";
import InvestmentDetailPage from "./pages/InvestmentDetailPage";
import NewInvestmentTransactionPage from "./pages/NewInvestmentTransactionPage";

// The route table is exported separately so tests can mount it with a
// memory history via createRouter({ routes, history: memoryHistory(path) }).
export const routes: RouteDefinition[] = [
  { path: "/", component: AccountsPage },
  { path: "/accounts/new", component: NewAccountPage },
  { path: "/accounts/:id", component: AccountDetailPage },
  { path: "/accounts/:id/movements/new", component: NewMovementPage },
  { path: "/accounts/:id/investments/new", component: NewInvestmentPage },
  { path: "/accounts/:id/investments/:investmentId", component: InvestmentDetailPage },
  { path: "/accounts/:id/investments/:investmentId/transactions/new", component: NewInvestmentTransactionPage },
];

export function AppRoutes() {
  const Router = createRouter({ routes });
  return <Router />;
}

export default function App() {
  return (
    <div class="max-w-3xl mx-auto px-4 py-8 flex flex-col gap-6">
      <h1 class="text-3xl font-semibold">
        <RouterLink to="/" class="hover:underline">
          Personal Finance Tracker
        </RouterLink>
      </h1>
      <AppRoutes />
    </div>
  );
}
