export function Loading(props: { label?: string }) {
  return (
    <p role="status" class="text-gray-500 dark:text-gray-400">
      {props.label ?? "Loading..."}
    </p>
  );
}
