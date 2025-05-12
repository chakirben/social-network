export function timePassed(date) {
    
    const now = new Date()
    const pastDate = new Date(date)
    const diff = now - pastDate 
  
    const seconds = Math.floor(diff / 1000)
    const minutes = Math.floor(seconds / 60)
    const hours = Math.floor(minutes / 60)
    const days = Math.floor(hours / 24)
    const weeks = Math.floor(days / 7)
    const months = Math.floor(weeks / 4)
    const years = Math.floor(months / 12)
  
    if (seconds < 60) {
      return `${seconds} sec ago`
    } else if (minutes < 60) {
      return `${minutes} mins ago`
    } else if (hours < 24) {
      return `${hours} hours ago`
    } else if (days < 7) {
      return `${days} days ago`
    } else if (weeks < 4) {
      return `${weeks} weeks ago`
    } else if (months < 12) {
      return `${months} monts ago`
    } else {
      return `${years} years ago`
    }
  }