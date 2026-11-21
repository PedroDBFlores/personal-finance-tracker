import type { JSX } from "@solidjs/web";

export function EmptyState(props: { children?: JSX.Element }) {
  return <p class="text-gray-500 dark:text-gray-400">{props.children}</p>;
}
