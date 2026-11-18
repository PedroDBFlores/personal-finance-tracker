import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import App from "./App";
import * as api from "./api";

vi.mock("./api", async (importOriginal) => {
  const actual = await importOriginal<typeof api>();
  return {
    ...actual,
    listAccounts: vi.fn(),
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

beforeEach(() => {
  vi.mocked(api.listAccounts).mockResolvedValue(accounts);
  vi.mocked(api.createAccount).mockResolvedValue(accounts[0]);
  vi.mocked(api.listMovements).mockResolvedValue(movements);
  vi.mocked(api.createMovement).mockResolvedValue(movements[0]);
});

describe("App", () => {
  it("loads and lists accounts", async () => {
    render(<App />);
    expect(await screen.findByText("Bank A")).toBeInTheDocument();
    expect(screen.getByText("Broker B")).toBeInTheDocument();
    expect(api.listAccounts).toHaveBeenCalled();
  });

  it("creates an account from the form", async () => {
    render(<App />);
    await screen.findByText("Bank A");

    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: "New Bank" } });
    fireEvent.change(screen.getByLabelText(/currency/i), { target: { value: "USD" } });
    fireEvent.click(screen.getByRole("button", { name: /add account/i }));

    await waitFor(() => expect(api.createAccount).toHaveBeenCalled());
    const input = vi.mocked(api.createAccount).mock.calls[0][0];
    expect(input.name).toBe("New Bank");
    expect(input.currency).toBe("USD");
  });

  it("shows movements when an account is selected", async () => {
    render(<App />);
    fireEvent.click(await screen.findByText("Bank A"));

    expect(await screen.findByText(/Salary/)).toBeInTheDocument();
    expect(api.listMovements).toHaveBeenCalledWith("1");
    expect(screen.getByText(/1000\.00/)).toBeInTheDocument();
  });
});
