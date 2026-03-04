import { FetchRPC } from 'twirp-ts';
import { PublicClientJSON } from '../codegen/public.twirp';
import { UserClientJSON } from '../codegen/user.twirp';

function isDevMode() {
  // We assume we're only on a port in dev mode.
  return window.location.port.length > 0;
}

// NOTE: if you want to serve the frontend from the go sever in dev mode, FIRST
// switch this variable to true and THEN run `npm run build` from the /_ui 
// directory. Then load the RPC port in your browser, e.g. localhost:5001/
const serveUIFromGoInDevMode = false;

export function getRpcHost() {
  let url = window.location.origin;
  if (isDevMode() && !serveUIFromGoInDevMode) {
    const portNum = parseInt(window.location.port);
    // By convention, we always configure our RPC host to listen on one port
    // number higher than our UI server in dev mode.
    url = url.replace(`:${portNum}`, `:${portNum + 1}`);
  }
  return url;
}

export const publicClient = new PublicClientJSON(FetchRPC({
  baseUrl: getRpcHost() + "/rpc/public",
  credentials: "include",
}));

export const userClient = new UserClientJSON(FetchRPC({
  baseUrl: getRpcHost() + "/rpc/user",
  credentials: "include",
}));
// other clients can be added with different prefixes...
