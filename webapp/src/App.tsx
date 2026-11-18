import { BrowserRouter, Route, Routes } from "react-router";
import { Container, Heading, VStack } from "@chakra-ui/react";
import { RouterLink } from "./components/RouterLink";
import AccountsPage from "./pages/AccountsPage";
import NewAccountPage from "./pages/NewAccountPage";
import AccountDetailPage from "./pages/AccountDetailPage";
import NewMovementPage from "./pages/NewMovementPage";
import NewInvestmentPage from "./pages/NewInvestmentPage";
import InvestmentDetailPage from "./pages/InvestmentDetailPage";
import NewInvestmentTransactionPage from "./pages/NewInvestmentTransactionPage";

// AppRoutes is exported separately so tests can render it inside a MemoryRouter.
export function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<AccountsPage />} />
      <Route path="/accounts/new" element={<NewAccountPage />} />
      <Route path="/accounts/:id" element={<AccountDetailPage />} />
      <Route path="/accounts/:id/movements/new" element={<NewMovementPage />} />
      <Route path="/accounts/:id/investments/new" element={<NewInvestmentPage />} />
      <Route path="/accounts/:id/investments/:investmentId" element={<InvestmentDetailPage />} />
      <Route
        path="/accounts/:id/investments/:investmentId/transactions/new"
        element={<NewInvestmentTransactionPage />}
      />
    </Routes>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <Container maxW="container.md" py={8}>
        <VStack align="stretch" gap={6}>
          <Heading size="2xl">
            <RouterLink to="/">Personal Finance Tracker</RouterLink>
          </Heading>
          <AppRoutes />
        </VStack>
      </Container>
    </BrowserRouter>
  );
}
