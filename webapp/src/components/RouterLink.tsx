import { chakra } from "@chakra-ui/react";
import { Link as ReactRouterLink } from "react-router";

// chakra() wraps react-router's Link so it accepts Chakra style props while
// `to` and friends are consumed by the router.
export const RouterLink = chakra(ReactRouterLink);
