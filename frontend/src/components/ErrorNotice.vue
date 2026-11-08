<script lang="ts" setup>
import type { AppError } from '../types/errors'
defineProps<{ error: AppError }>()
const emit = defineEmits<{ dismiss: [] }>()
</script>
<template><section class="error-notice" role="alert" data-testid="error-notice"><div><strong>{{ error.code }}</strong><span>{{ error.message }}</span></div><button type="button" aria-label="关闭错误" @click="emit('dismiss')">×</button></section></template>
