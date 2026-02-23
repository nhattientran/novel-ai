import { ref, onUnmounted } from 'vue'

export function useDebounceFn<T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  const timeoutId = ref<ReturnType<typeof setTimeout> | null>(null)

  onUnmounted(() => {
    if (timeoutId.value) {
      clearTimeout(timeoutId.value)
    }
  })

  return (...args: Parameters<T>) => {
    if (timeoutId.value) {
      clearTimeout(timeoutId.value)
    }

    timeoutId.value = setTimeout(() => {
      fn(...args)
    }, delay)
  }
}
