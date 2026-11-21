import type { JSX } from "@solidjs/web";

type FormFieldProps = {
  id: string;
  label: string;
  class?: string;
  children?: JSX.Element;
};

export function FormField(props: FormFieldProps) {
  return (
    <div class={props.class ?? ""}>
      <label class="block text-sm font-medium mb-1" for={props.id}>
        {props.label}
      </label>
      {props.children}
    </div>
  );
}
