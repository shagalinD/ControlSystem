export const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

export const validatePassword = (password: string): boolean => {
  return password.length >= 6
}

export const validateRequired = (value: string): boolean => {
  return value.trim().length > 0
}

export const validateDeadline = (deadline: string): boolean => {
  if (!deadline) return true
  const deadlineDate = new Date(deadline)
  const today = new Date()
  return deadlineDate > today
}
