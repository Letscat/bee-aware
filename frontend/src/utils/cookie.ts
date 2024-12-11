export function getCookie(name:string):string|null {
    const cookies = document.cookie.split(';');
    if (!cookies) return null;
    for (let cookie of cookies) {
      const [key, value] = cookie.trim().split('=');
      if (key === name) return value;
    }
    return null;
  }
  
  export function setCookie(name:string, value:string, days:number) {
    const date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    const expires = 'expires=' + date.toUTCString();
    document.cookie = name + '=' + value + ';' + expires + ';path=/';
  }