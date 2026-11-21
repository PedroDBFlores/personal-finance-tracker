import { createRouter } from "@solidjs/router";
import type { RouteDefinition } from "@solidjs/router";
import { NavBar } from "./components/molecules/NavBar";
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
    <div>
      <NavBar />
      <main class="max-w-3xl mx-auto px-4 py-8">
        <AppRoutes />
      </main>
    </div>
  );
}
