import { writable, type Writable } from "svelte/store";

export class Loader {
  loading: Writable<boolean>;
  currentlyLoading: Set<string>;

  constructor() {
    this.loading = writable(false);
    this.currentlyLoading = new Set<string>();
  }

  store(): Writable<boolean> {
    return this.loading;
  }

  start(id?: string) {
    if (!id) id = `${Math.random() * 10000000}`;
    this.currentlyLoading.add(id);
    this.loading.set(true);

    return () => {
      this.currentlyLoading.delete(id);
      if (this.currentlyLoading.size == 0) this.loading.set(false);
    }
  }
}
