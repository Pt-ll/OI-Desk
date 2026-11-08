export type AppErrorCode = 'UNKNOWN' | 'VALIDATION_FAILED' | 'PERMISSION_DENIED' | 'NOT_FOUND' | 'NETWORK_ERROR' | 'DATABASE_ERROR' | 'PROVIDER_ERROR'

export interface AppError {
  code: AppErrorCode | string
  message: string
  retryable?: boolean
  details?: string
}

export function normalizeError(error: unknown, fallback = '操作失败'): AppError {
  if (typeof error === 'object' && error !== null) {
    const candidate = error as Partial<AppError> & { message?: unknown }
    if (typeof candidate.message === 'string') return { code: typeof candidate.code === 'string' ? candidate.code : 'UNKNOWN', message: candidate.message, retryable: candidate.retryable, details: typeof candidate.details === 'string' ? candidate.details : undefined }
  }
  if (error instanceof Error && error.message) return { code: 'UNKNOWN', message: error.message }
  return { code: 'UNKNOWN', message: fallback }
}
