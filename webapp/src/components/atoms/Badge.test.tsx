import { describe, it, expect } from "vitest";
import { render, screen } from "@solidjs/testing-library";
import { Badge } from "./Badge";

describe("Badge", () => {
  it("renders a pill with the base and tone classes", () => {
    render(() => <Badge tone="purple">USD</Badge>);
    const badge = screen.getByText("USD");
    expect(badge.className).toContain("badge");
    expect(badge.className).toContain("badge-purple");
  });

  it("defaults to the neutral tone", () => {
    render(() => <Badge>Bank</Badge>);
    const badge = screen.getByText("Bank");
    expect(badge.className).toContain("badge");
    expect(badge.className).not.toContain("badge-purple");
  });
});
