import { describe, it, expect } from "vitest";
import { render, screen, within } from "@solidjs/testing-library";
import { ListRow } from "./ListRow";
import { Badge } from "../atoms/Badge";

describe("ListRow", () => {
  it("renders left and right content in a flex card row", () => {
    render(() => (
      <ul>
        <ListRow left={<span>Left</span>} right={<Badge>Right</Badge>} />
      </ul>
    ));
    const item = screen.getByRole("listitem");
    const card = item.firstElementChild!;
    expect(card.className).toContain("card");
    expect(card.className).toContain("flex");
    expect(card.className).toContain("justify-between");
    expect(within(item).getByText("Left")).toBeInTheDocument();
    expect(within(item).getByText("Right")).toBeInTheDocument();
  });
});
