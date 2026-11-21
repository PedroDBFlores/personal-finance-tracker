import type { JSX } from "@solidjs/web";

// The Solid router intercepts clicks on plain anchors, so a router link is an
// <a href>; `to` is mapped to `href`.
type RouterLinkProps = {
  to: string;
  class?: string;
  children?: JSX.Element;
};

export function RouterLink(props: RouterLinkProps) {
  return (
    <a href={props.to} class={props.class}>
      {props.children}
    </a>
  );
}
