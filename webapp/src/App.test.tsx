import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@solidjs/testing-library";
import { createRouter, memoryHistory } from "@solidjs/router";
import { routes } from "./App";
import * as api from "./api";

vi.mock("./api", async (importOriginal) => {
  const actual = await importOriginal<typeof api>();
  return {
    ...actual,
    listAccounts: vi.fn(),
    getAccount: vi.fn(),
    getAccountBalance: vi.fn(),
    createAccount: vi.fn(),
    listMovements: vi.fn(),
    createMovement: vi.fn(),
    listMovementMonths: vi.fn(),
    listMovementsByMonth: vi.fn(),
    listInvestments: vi.fn(),
    getInvestment: vi.fn(),
    createInvestment: vi.fn(),
    listInvestmentTransactions: vi.fn(),
    createInvestmentTransaction: vi.fn(),
  };
});

const accounts: api.Account[] = [
  { id: "1", name: "Bank A", type: "Bank", currency: "EUR", created_at: "2024-01-02T03:04:05Z" },
  { id: "2", name: "Broker B", type: "Brokerage", currency: "USD", created_at: "2024-01-02T03:04:05Z" },
];

const movements: api.Movement[] = [
  { id: "1", account_id: "1", description: "Salary", value: "1000.00", type: "Credit", date: "2024-01-02T03:04:05Z" },
];

const investment: api.Investment = {
  id: "inv-1",
  account_id: "1",
  name: "Vanguard ETF",
  type: "ETF",
  current_value: "760.00",
  currency: "USD",
  purchase_date: "2024-01-02T03:04:05Z",
  maturity_date: null,
};

const investmentTransactions: api.InvestmentTransaction[] = [
  { id: "tx-1", investment_id: "inv-1", type: "Buy", amount: "0.00", quantity: "10.00", price: "100.00", date: "2024-01-02T03:04:05Z", description: "Initial Buy" },
];

function renderAt(path: string) {
  const Router = createRouter({ routes, history: memoryHistory(path) });
  return render(() => <Router />);
}

// Solid 2.0 batches signal writes to a microtask; await one so form fields are
// committed before the next interaction.
const tick = () => Promise.resolve();

beforeEach(() => {
  vi.mocked(api.listAccounts).mockResolvedValue(accounts);
  vi.mocked(api.getAccount).mockResolvedValue(accounts[0]);
  vi.mocked(api.getAccountBalance).mockImplementation((id: string) =>
    Promise.resolve({
      account_id: id,
      currency: "EUR",
      movements_balance: "974.01",
      investments_balance: "1000.00",
      total_balance: "1974.01",
    }),
  );
  vi.mocked(api.createAccount).mockResolvedValue(accounts[0]);
  vi.mocked(api.listMovements).mockResolvedValue(movements);
  vi.mocked(api.createMovement).mockResolvedValue(movements[0]);
  vi.mocked(api.listMovementMonths).mockResolvedValue([{ year: 2024, month: 1 }]);
  vi.mocked(api.listMovementsByMonth).mockResolvedValue(movements);
  vi.mocked(api.listInvestments).mockResolvedValue([investment]);
  vi.mocked(api.getInvestment).mockResolvedValue(investment);
  vi.mocked(api.createInvestment).mockResolvedValue(investment);
  vi.mocked(api.listInvestmentTransactions).mockResolvedValue(investmentTransactions);
  vi.mocked(api.createInvestmentTransaction).mockResolvedValue(investmentTransactions[0]);
});

describe("routing", () => {
  it("main page lists accounts", async () => {
    renderAt("/");
    expect(await screen.findByText("Bank A")).toBeInTheDocument();
    expect(screen.getByText("Broker B")).toBeInTheDocument();
    expect(api.listAccounts).toHaveBeenCalled();
  });

  it("main page shows each account's total balance", async () => {
    renderAt("/");
    expect((await screen.findAllByText(/1974\.01/)).length).toBe(2);
    expect(api.getAccountBalance).toHaveBeenCalledWith("1");
    expect(api.getAccountBalance).toHaveBeenCalledWith("2");
  });

  it("account detail page shows the balance breakdown", async () => {
    renderAt("/accounts/1");
    expect(await screen.findByText(/1974\.01 EUR/)).toBeInTheDocument();
    expect(screen.getByText(/Movements 974\.01/)).toBeInTheDocument();
    expect(screen.getByText(/Investments 1000\.00/)).toBeInTheDocument();
    expect(api.getAccountBalance).toHaveBeenCalledWith("1");
  });

  it("navigates to the new-account page, creates, and returns to the list", async () => {
    renderAt("/");
    fireEvent.click(await screen.findByRole("link", { name: /new account/i }));

    fireEvent.input(await screen.findByLabelText(/name/i), { target: { value: "New Bank" } });
    await tick();
    fireEvent.change(screen.getByLabelText(/currency/i), { target: { value: "USD" } });
    await tick();
    fireEvent.click(screen.getByRole("button", { name: /add account/i }));

    await waitFor(() => expect(api.createAccount).toHaveBeenCalled());
    expect(vi.mocked(api.createAccount).mock.calls[0][0].name).toBe("New Bank");
    // back on the list, which reloads
    await waitFor(() => expect(api.listAccounts).toHaveBeenCalledTimes(2));
  });

  it("account detail page shows movements via the month browser", async () => {
    renderAt("/accounts/1");

    expect(await screen.findByText(/Salary/)).toBeInTheDocument();
    expect(api.getAccount).toHaveBeenCalledWith("1");
    expect(api.listMovementMonths).toHaveBeenCalledWith("1");
    expect(api.listMovementsByMonth).toHaveBeenCalledWith("1", 2024, 1);
  });

  it("navigates from detail to the new-movement page, creates, and returns", async () => {
    renderAt("/accounts/1");
    fireEvent.click(await screen.findByRole("link", { name: /add movement/i }));

    fireEvent.input(await screen.findByLabelText(/description/i), { target: { value: "Groceries" } });
    await tick();
    fireEvent.input(screen.getByLabelText(/amount/i), { target: { value: "25.99" } });
    await tick();
    fireEvent.click(screen.getByRole("button", { name: /add movement/i }));

    await waitFor(() => expect(api.createMovement).toHaveBeenCalledWith("1", { description: "Groceries", value: "25.99", type: "Credit" }));
    // back on the detail page, the month browser reloads
    await waitFor(() => expect(api.listMovementMonths).toHaveBeenCalledTimes(2));
  });

  it("account detail page also lists investments", async () => {
    renderAt("/accounts/1");

    expect(await screen.findByText("Vanguard ETF")).toBeInTheDocument();
    expect(api.listInvestments).toHaveBeenCalledWith("1");
  });

  it("navigates from detail to the new-investment page, creates, and returns", async () => {
    renderAt("/accounts/1");
    fireEvent.click(await screen.findByRole("link", { name: /new investment/i }));

    fireEvent.input(await screen.findByLabelText(/name/i), { target: { value: "New ETF" } });
    await tick();
    fireEvent.click(screen.getByRole("button", { name: /add investment/i }));

    await waitFor(() => expect(api.createInvestment).toHaveBeenCalledWith("1", { name: "New ETF", type: "ETF", currency: "USD" }));
    // back on the detail page, investments reloaded
    await waitFor(() => expect(api.listInvestments).toHaveBeenCalledTimes(2));
  });

  it("investment detail page shows its transactions", async () => {
    renderAt("/accounts/1/investments/inv-1");

    expect(await screen.findByText(/Initial Buy/)).toBeInTheDocument();
    expect(api.getInvestment).toHaveBeenCalledWith("1", "inv-1");
    expect(api.listInvestmentTransactions).toHaveBeenCalledWith("1", "inv-1");
  });

  it("navigates from investment detail to the new-transaction page, creates, and returns", async () => {
    renderAt("/accounts/1/investments/inv-1");
    fireEvent.click(await screen.findByRole("link", { name: /add transaction/i }));

    fireEvent.change(await screen.findByLabelText(/transaction type/i), { target: { value: "Deposit" } });
    await tick();
    fireEvent.input(screen.getByLabelText(/amount/i), { target: { value: "500.00" } });
    await tick();
    fireEvent.input(screen.getByLabelText(/description/i), { target: { value: "Deposit" } });
    await tick();
    fireEvent.click(screen.getByRole("button", { name: /add transaction/i }));

    await waitFor(() =>
      expect(api.createInvestmentTransaction).toHaveBeenCalledWith("1", "inv-1", {
        type: "Deposit",
        amount: "500.00",
        quantity: "",
        price: "",
        description: "Deposit",
      }),
    );
    // back on the investment detail page, transactions reloaded
    await waitFor(() => expect(api.listInvestmentTransactions).toHaveBeenCalledTimes(2));
  });

  it("computes the amount from quantity and price for buy/sell transactions", async () => {
    renderAt("/accounts/1/investments/inv-1");
    fireEvent.click(await screen.findByRole("link", { name: /add transaction/i }));

    fireEvent.change(await screen.findByLabelText(/transaction type/i), { target: { value: "Buy" } });
    await tick();
    fireEvent.input(screen.getByLabelText(/quantity/i), { target: { value: "10" } });
    await tick();
    fireEvent.input(screen.getByLabelText(/price/i), { target: { value: "100.00" } });
    await tick();
    fireEvent.click(screen.getByRole("button", { name: /add transaction/i }));

    await waitFor(() =>
      expect(api.createInvestmentTransaction).toHaveBeenCalledWith("1", "inv-1", {
        type: "Buy",
        amount: "1000.00",
        quantity: "10",
        price: "100.00",
        description: "",
      }),
    );
  });
});
