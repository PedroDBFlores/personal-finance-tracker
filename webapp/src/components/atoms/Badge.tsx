import type { JSX } from "@solidjs/web";

type BadgeProps = {
  tone?: "default" | "green" | "red" | "purple";
  children?: JSX.Element;
};

export function Badge(props: BadgeProps) {
  const toneClass = () => (props.tone && props.tone !== "default" ? ` badge-${props.tone}` : "");
  return <span class={`badge${toneClass()}`}>{props.children}</span>;
}
