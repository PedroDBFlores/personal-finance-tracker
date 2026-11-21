import { describe, it, expect } from "vitest";
import { render, screen, within } from "@solidjs/testing-library";
import { AccountList } from "./AccountList";
import type { Account } from "../../api";

const accounts: Account[] = [
  { id: "1", name: "Bank A", type: "Bank", currency: "EUR", created_at: "2024-01-02T03:04:05Z" },
  { id: "2", name: "Broker B", type: "Brokerage", currency: "USD", created_at: "2024-01-02T03:04:05Z" },
];

describe("AccountList", () => {
  it("lays out each row as flex with the name on the left and the pills on the right", () => {
    render(() => <AccountList accounts={accounts} />);
    const rows = screen.getAllByRole("listitem");
    expect(rows).toHaveLength(2);

    const first = rows[0];
    const card = first.firstElementChild!;
    expect(card.className).toContain("flex");
    expect(card.className).toContain("justify-between");

    const link = within(first).getByRole("link", { name: "Bank A" });
    expect(link).toHaveAttribute("href", "/accounts/1");
    expect(within(first).getByText("Bank")).toBeInTheDocument();
    expect(within(first).getByText("EUR")).toBeInTheDocument();
  });

  it("shows an empty state when there are no accounts", () => {
    render(() => <AccountList accounts={[]} />);
    expect(screen.getByText(/no accounts yet/i)).toBeInTheDocument();
  });
});
