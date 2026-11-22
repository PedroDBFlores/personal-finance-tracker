import { RouterLink } from "../atoms/RouterLink";

export function NavBar() {
  return (
    <header class="bg-brand-500 text-white">
      <div class="max-w mx-auto px-4 py-3 flex items-center justify-between">
        <RouterLink to="/" class="text-lg font-semibold text-white hover:underline">
          Personal Finance Tracker
        </RouterLink>
      </div>
    </header>
  );
}
