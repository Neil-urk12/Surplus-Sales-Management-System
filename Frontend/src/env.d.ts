/// <reference types="vite/client" />

declare namespace NodeJS {
  interface ProcessEnv {
    NODE_ENV: string;
    VUE_ROUTER_MODE: 'hash' | 'history' | 'abstract' | undefined;
    VUE_ROUTER_BASE: string | undefined;
  }
}

declare interface ImportMetaEnv {
  readonly VITE_TURNSTILE_SITE_KEY: string
  // add other env variables here if needed
}

declare interface ImportMeta {
  readonly env: ImportMetaEnv
}
