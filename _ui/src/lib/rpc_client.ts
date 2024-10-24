import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { CounterClient } from "../codegen/countservice.client";

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

function getTransport(path: string) {
  return new TwirpFetchTransport({
    baseUrl: getRpcHost() + path,
    // this is only necessary because in dev we run the UI and RPC servers on
    // different ports.
    fetchInit: { credentials: 'include' },
  });
}

export const client = new CounterClient(getTransport("/rpc/");
// other clients can be added with different prefixes like
// export const counterClient = new CounterClient(getTransport("/rpc/counter/");
// export const otherClient = new OtherClient(getTransport("/rpc/other/");
