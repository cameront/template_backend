import {
  TwirpContext,
  TwirpServer,
  RouterEvents,
  TwirpError,
  TwirpErrorCode,
  Interceptor,
  TwirpContentType,
  chainInterceptors,
} from "twirp-ts";
import { CounterRequest, CounterValue, IncrementRequest } from "./countservice";

//==================================//
//          Client Code             //
//==================================//

interface Rpc {
  request(
    service: string,
    method: string,
    contentType: "application/json" | "application/protobuf",
    data: object | Uint8Array
  ): Promise<object | Uint8Array>;
}

export interface CounterClient {
  GetValue(request: CounterRequest): Promise<CounterValue>;
  Increment(request: IncrementRequest): Promise<CounterValue>;
}

export class CounterClientJSON implements CounterClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.GetValue.bind(this);
    this.Increment.bind(this);
  }
  GetValue(request: CounterRequest): Promise<CounterValue> {
    const data = CounterRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "counter.Counter",
      "GetValue",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      CounterValue.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }

  Increment(request: IncrementRequest): Promise<CounterValue> {
    const data = IncrementRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "counter.Counter",
      "Increment",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      CounterValue.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }
}

export class CounterClientProtobuf implements CounterClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.GetValue.bind(this);
    this.Increment.bind(this);
  }
  GetValue(request: CounterRequest): Promise<CounterValue> {
    const data = CounterRequest.toBinary(request);
    const promise = this.rpc.request(
      "counter.Counter",
      "GetValue",
      "application/protobuf",
      data
    );
    return promise.then((data) => CounterValue.fromBinary(data as Uint8Array));
  }

  Increment(request: IncrementRequest): Promise<CounterValue> {
    const data = IncrementRequest.toBinary(request);
    const promise = this.rpc.request(
      "counter.Counter",
      "Increment",
      "application/protobuf",
      data
    );
    return promise.then((data) => CounterValue.fromBinary(data as Uint8Array));
  }
}

//==================================//
//          Server Code             //
//==================================//

export interface CounterTwirp<T extends TwirpContext = TwirpContext> {
  GetValue(ctx: T, request: CounterRequest): Promise<CounterValue>;
  Increment(ctx: T, request: IncrementRequest): Promise<CounterValue>;
}

export enum CounterMethod {
  GetValue = "GetValue",
  Increment = "Increment",
}

export const CounterMethodList = [
  CounterMethod.GetValue,
  CounterMethod.Increment,
];

export function createCounterServer<T extends TwirpContext = TwirpContext>(
  service: CounterTwirp<T>
) {
  return new TwirpServer<CounterTwirp, T>({
    service,
    packageName: "counter",
    serviceName: "Counter",
    methodList: CounterMethodList,
    matchRoute: matchCounterRoute,
  });
}

function matchCounterRoute<T extends TwirpContext = TwirpContext>(
  method: string,
  events: RouterEvents<T>
) {
  switch (method) {
    case "GetValue":
      return async (
        ctx: T,
        service: CounterTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, CounterRequest, CounterValue>[]
      ) => {
        ctx = { ...ctx, methodName: "GetValue" };
        await events.onMatch(ctx);
        return handleCounterGetValueRequest(ctx, service, data, interceptors);
      };
    case "Increment":
      return async (
        ctx: T,
        service: CounterTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, IncrementRequest, CounterValue>[]
      ) => {
        ctx = { ...ctx, methodName: "Increment" };
        await events.onMatch(ctx);
        return handleCounterIncrementRequest(ctx, service, data, interceptors);
      };
    default:
      events.onNotFound();
      const msg = `no handler found`;
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleCounterGetValueRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: CounterTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, CounterRequest, CounterValue>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleCounterGetValueJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handleCounterGetValueProtobuf<T>(ctx, service, data, interceptors);
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleCounterIncrementRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: CounterTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, IncrementRequest, CounterValue>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleCounterIncrementJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handleCounterIncrementProtobuf<T>(
        ctx,
        service,
        data,
        interceptors
      );
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}
async function handleCounterGetValueJSON<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: CounterTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, CounterRequest, CounterValue>[]
) {
  let request: CounterRequest;
  let response: CounterValue;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = CounterRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      CounterRequest,
      CounterValue
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.GetValue(ctx, inputReq);
    });
  } else {
    response = await service.GetValue(ctx, request!);
  }

  return JSON.stringify(
    CounterValue.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}

async function handleCounterIncrementJSON<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: CounterTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, IncrementRequest, CounterValue>[]
) {
  let request: IncrementRequest;
  let response: CounterValue;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = IncrementRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      IncrementRequest,
      CounterValue
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Increment(ctx, inputReq);
    });
  } else {
    response = await service.Increment(ctx, request!);
  }

  return JSON.stringify(
    CounterValue.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}
async function handleCounterGetValueProtobuf<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: CounterTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, CounterRequest, CounterValue>[]
) {
  let request: CounterRequest;
  let response: CounterValue;

  try {
    request = CounterRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      CounterRequest,
      CounterValue
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.GetValue(ctx, inputReq);
    });
  } else {
    response = await service.GetValue(ctx, request!);
  }

  return Buffer.from(CounterValue.toBinary(response));
}

async function handleCounterIncrementProtobuf<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: CounterTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, IncrementRequest, CounterValue>[]
) {
  let request: IncrementRequest;
  let response: CounterValue;

  try {
    request = IncrementRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      IncrementRequest,
      CounterValue
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Increment(ctx, inputReq);
    });
  } else {
    response = await service.Increment(ctx, request!);
  }

  return Buffer.from(CounterValue.toBinary(response));
}
