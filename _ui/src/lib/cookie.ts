const cookieStore = { set: false };

export function setCookie() {
  cookieStore.set = true;
}

export function deleteCookie(authCookieName = "ac", cb: () => void) {
  document.cookie = authCookieName + "=" + ";expires=Thu, 01 Jan 1970 00:00:01 GMT";
  cb();
}