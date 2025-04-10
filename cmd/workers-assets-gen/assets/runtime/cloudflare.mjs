import { connect } from "cloudflare:sockets";

export function createRuntimeContext({ env, ctx, binding }) {
  return {
    env,
    ctx,
    connect,
    binding,
  };
}
