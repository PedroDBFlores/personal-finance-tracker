import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor, within } from "@solidjs/testing-library";
import { MovementMonthBrowser } from "./MovementMonthBrowser";
import * as api from "../../api";

vi.mock("../../api", async (importOriginal) => {
  const actual = await importOriginal<typeof api>();
  return {
    ...actual,
    listMovementMonths: vi.fn(),
    listMovementsByMonth: vi.fn(),
  };
});

const months: api.MovementMonth[] = [
  { year: 2025, month: 2 },
  { year: 2025, month: 1 },
];

const tick = () => Promise.resolve();

beforeEach(() => {
  vi.mocked(api.listMovementMonths).mockResolvedValue(months);
  vi.mocked(api.listMovementsByMonth).mockImplementation((accountId: string, year: number, month: number) =>
    Promise.resolve(
      year === 2025 && month === 2
        ? [{ id: "9", account_id: accountId, description: "Feb Salary", value: "1100.00", type: "Credit", date: "2025-02-02T00:00:00Z" }]
        : [{ id: "1", account_id: accountId, description: "Jan Salary", value: "1000.00", type: "Credit", date: "2025-01-02T00:00:00Z" }],
    ),
  );
});

describe("MovementMonthBrowser", () => {
  it("defaults to the most recent month and shows its movements", async () => {
    render(() => <MovementMonthBrowser accountId="1" />);

    expect(await screen.findByText("Feb Salary")).toBeInTheDocument();
    expect(api.listMovementMonths).toHaveBeenCalledWith("1");
    expect(api.listMovementsByMonth).toHaveBeenCalledWith("1", 2025, 2);
    expect(screen.getByRole("combobox", { name: /month/i })).toHaveValue("2025-2");
  });

  it("loads the movements of the picked month", async () => {
    render(() => <MovementMonthBrowser accountId="1" />);
    await screen.findByText("Feb Salary");

    fireEvent.change(screen.getByRole("combobox", { name: /month/i }), { target: { value: "2025-1" } });
    await tick();

    expect(await screen.findByText("Jan Salary")).toBeInTheDocument();
    expect(api.listMovementsByMonth).toHaveBeenCalledWith("1", 2025, 1);
  });

  it("shows each movement as a row with its type pill and value", async () => {
    render(() => <MovementMonthBrowser accountId="1" />);
    const row = (await screen.findByText("Feb Salary")).closest("li")!;
    expect(within(row).getByText("Credit")).toBeInTheDocument();
    expect(within(row).getByText(/1100\.00/)).toBeInTheDocument();
  });

  it("shows an empty state when the account has no movements", async () => {
    vi.mocked(api.listMovementMonths).mockResolvedValue([]);
    render(() => <MovementMonthBrowser accountId="1" />);

    expect(await screen.findByText(/no movements yet/i)).toBeInTheDocument();
    expect(api.listMovementsByMonth).not.toHaveBeenCalled();
  });

  it("shows an error when loading months fails", async () => {
    vi.mocked(api.listMovementMonths).mockRejectedValue(new Error("boom"));
    render(() => <MovementMonthBrowser accountId="1" />);

    await waitFor(() => expect(screen.getByRole("alert")).toHaveTextContent("boom"));
  });
});
