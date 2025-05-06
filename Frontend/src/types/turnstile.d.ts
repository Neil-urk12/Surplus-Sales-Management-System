declare namespace Turnstile {
  interface RenderOptions {
    sitekey: string
    callback?: (token: string) => void
    'error-callback'?: () => void
    'expired-callback'?: () => void
    theme?: 'light' | 'dark'
    tabindex?: number
  }
  function render(el: HTMLElement, opts: RenderOptions): string
}

declare global {
  interface Window {
    turnstile: typeof Turnstile
  }
}