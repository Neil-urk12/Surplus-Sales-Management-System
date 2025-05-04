export const materialImages = {
  steel: 'https://loremflickr.com/600/400/steel',
  concrete: 'https://loremflickr.com/600/400/concrete',
  lumber: 'https://loremflickr.com/600/400/lumber'
}

// Helper function to get image URL by material type
export const getMaterialImage = (type: string): string => {
  const key = type.toLowerCase() as keyof typeof materialImages
  return materialImages[key] || 'https://loremflickr.com/600/400/material' // fallback image
} 