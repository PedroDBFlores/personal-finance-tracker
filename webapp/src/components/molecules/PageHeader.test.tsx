import { describe, it, expect } from "vitest";
import { render, screen } from "@solidjs/testing-library";
import { PageHeader } from "./PageHeader";

describe("PageHeader", () => {
  it("renders the title and the action link", () => {
    render(() => <PageHeader title="Accounts" action={{ to: "/accounts/new", label: "New account" }} />);
    expect(screen.getByRole("heading", { name: "Accounts" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "New account" })).toHaveAttribute("href", "/accounts/new");
  });

  it("renders without an action", () => {
    render(() => <PageHeader title="Details" />);
    expect(screen.getByRole("heading", { name: "Details" })).toBeInTheDocument();
    expect(screen.queryByRole("link")).toBeNull();
  });
});
