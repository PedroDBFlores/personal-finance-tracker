import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router";
import { AppRoutes } from "./App";
import { Provider } from "./components/ui/provider";
import * as api from "./api";

vi.mock("./api", async (importOriginal) => {
  const actual = await importOriginal<typeof api>();
  return {
    ...actual,
    listAccounts: vi.fn(),
    getAccount: vi.fn(),
    createAccount: vi.fn(),
    listMovements: vi.fn(),
    createMovement: vi.fn(),
  };
});

const accounts: api.Account[] = [
  { id: "1", name: "Bank A", type: "Bank", currency: "EUR", created_at: "2024-01-02T03:04:05Z" },
  { id: "2", name: "Broker B", type: "Brokerage", currency: "USD", created_at: "2024-01-02T03:04:05Z" },
];

const movements: api.Movement[] = [
  { id: "1", account_id: "1", description: "Salary", value: "1000.00", type: "Credit", date: "2024-01-02T03:04:05Z" },
];

function renderAt(path: string) {
  return render(
    <Provider>
      <MemoryRouter initialEntries={[path]}>
        <AppRoutes />
      </MemoryRouter>
    </Provider>,
  );
}

beforeEach(() => {
  vi.mocked(api.listAccounts).mockResolvedValue(accounts);
  vi.mocked(api.getAccount).mockResolvedValue(accounts[0]);
  vi.mocked(api.createAccount).mockResolvedValue(accounts[0]);
  vi.mocked(api.listMovements).mockResolvedValue(movements);
  vi.mocked(api.createMovement).mockResolvedValue(movements[0]);
});

describe("routing", () => {
  it("main page lists accounts", async () => {
    renderAt("/");
    expect(await screen.findByText("Bank A")).toBeInTheDocument();
    expect(screen.getByText("Broker B")).toBeInTheDocument();
    expect(api.listAccounts).toHaveBeenCalled();
  });

  it("navigates to the new-account page, creates, and returns to the list", async () => {
    renderAt("/");
    fireEvent.click(await screen.findByRole("link", { name: /new account/i }));

    fireEvent.change(await screen.findByLabelText(/name/i), { target: { value: "New Bank" } });
    fireEvent.change(screen.getByLabelText(/currency/i), { target: { value: "USD" } });
    fireEvent.click(screen.getByRole("button", { name: /add account/i }));

    await waitFor(() => expect(api.createAccount).toHaveBeenCalled());
    expect(vi.mocked(api.createAccount).mock.calls[0][0].name).toBe("New Bank");
    // back on the list, which reloads
    await waitFor(() => expect(api.listAccounts).toHaveBeenCalledTimes(2));
  });

  it("account detail page shows movements", async () => {
    renderAt("/accounts/1");

    expect(await screen.findByText(/Salary/)).toBeInTheDocument();
    expect(api.getAccount).toHaveBeenCalledWith("1");
    expect(api.listMovements).toHaveBeenCalledWith("1");
  });

  it("navigates from detail to the new-movement page, creates, and returns", async () => {
    renderAt("/accounts/1");
    fireEvent.click(await screen.findByRole("link", { name: /add movement/i }));

    fireEvent.change(await screen.findByLabelText(/description/i), { target: { value: "Groceries" } });
    fireEvent.change(screen.getByLabelText(/amount/i), { target: { value: "25.99" } });
    fireEvent.click(screen.getByRole("button", { name: /add movement/i }));

    await waitFor(() => expect(api.createMovement).toHaveBeenCalledWith("1", { description: "Groceries", value: "25.99", type: "Credit" }));
    // back on the detail page, movements reloaded
    await waitFor(() => expect(api.listMovements).toHaveBeenCalledTimes(2));
  });
});
