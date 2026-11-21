import { Show } from "solid-js";
import type { JSX } from "@solidjs/web";

type ListRowProps = {
  left: JSX.Element;
  right?: JSX.Element;
};

export function ListRow(props: ListRowProps) {
  return (
    <li>
      <div class="card flex justify-between items-center gap-2">
        <span>{props.left}</span>
        <Show when={props.right}>
          {(right) => <span class="flex gap-2 items-center">{right()}</span>}
        </Show>
      </div>
    </li>
  );
}
