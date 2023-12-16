import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { CounterClient } from "../codegen/countservice.client";

export function getRpcHost() {
    let url = window.location.origin;
    if (url.includes(":5000")) {
        url = url.replace(":5000", ":5001"); // localhost dev ui server listens on 5000, rpc on 5001
    }
    return url;
}

let transport = new TwirpFetchTransport({
    baseUrl: getRpcHost() + "/rpc",
    // this is only necessary because in dev we run the UI and RPC servers on
    // different ports.
    fetchInit: { credentials: 'include' },
});

export const client = new CounterClient(transport);
