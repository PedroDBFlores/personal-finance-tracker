import { Show } from "solid-js";
import { RouterLink } from "../atoms/RouterLink";

type PageHeaderProps = {
  title: string;
  action?: { to: string; label: string };
};

export function PageHeader(props: PageHeaderProps) {
  return (
    <div class="flex justify-between items-center mb-3">
      <h2 class="text-xl font-semibold">{props.title}</h2>
      <Show when={props.action}>
        {(action) => (
          <RouterLink to={action().to} class="btn-primary">
            {action().label}
          </RouterLink>
        )}
      </Show>
    </div>
  );
}
