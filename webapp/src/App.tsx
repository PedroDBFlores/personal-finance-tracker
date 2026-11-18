import { BrowserRouter, Route, Routes } from "react-router";
import { Container, Heading, VStack } from "@chakra-ui/react";
import { RouterLink } from "./components/RouterLink";
import AccountsPage from "./pages/AccountsPage";
import NewAccountPage from "./pages/NewAccountPage";
import AccountDetailPage from "./pages/AccountDetailPage";
import NewMovementPage from "./pages/NewMovementPage";

// AppRoutes is exported separately so tests can render it inside a MemoryRouter.
export function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<AccountsPage />} />
      <Route path="/accounts/new" element={<NewAccountPage />} />
      <Route path="/accounts/:id" element={<AccountDetailPage />} />
      <Route path="/accounts/:id/movements/new" element={<NewMovementPage />} />
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
