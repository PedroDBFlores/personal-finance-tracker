import { describe, it, expect } from "vitest";
import { render, screen } from "@solidjs/testing-library";
import { NavBar } from "./NavBar";

describe("NavBar", () => {
  it("shows the app name on the left linking home, on an orange bar", () => {
    const { container } = render(() => <NavBar />);
    const header = container.querySelector("header");
    expect(header).not.toBeNull();
    expect(header!.className).toContain("bg-brand-500");
    const brand = screen.getByRole("link", { name: /personal finance tracker/i });
    expect(brand).toHaveAttribute("href", "/");
    expect(brand.className).toContain("text-white");
  });
});
