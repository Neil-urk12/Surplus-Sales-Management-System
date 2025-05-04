export const defaultImages = {
  cab: [
    'https://loremflickr.com/600/400/car',
    'https://placehold.co/600x400/png?text=Car',
    'https://via.placeholder.com/600x400.png?text=Car',
    '/assets/images/default-car.png'
  ],
  accessory: [
    'https://loremflickr.com/600/400/accessory',
    'https://placehold.co/600x400/png?text=Accessory',
    'https://via.placeholder.com/600x400.png?text=Accessory',
    '/assets/images/default-accessory.png'
  ]
} as const;

type ImageType = keyof typeof defaultImages;
type FallbackArray<T extends ImageType> = (typeof defaultImages)[T];
type FallbackImage<T extends ImageType> = FallbackArray<T>[number];

function getDefaultImage<T extends ImageType>(type: T): FallbackImage<T> {
  return defaultImages[type][0];
}

function getNextFallbackImage<T extends ImageType>(currentImage: string, type: T): FallbackImage<T> {
  const fallbacks = defaultImages[type] as readonly string[];
  let currentIndex = -1;
  
  // Find the index of the current image
  for (let i = 0; i < fallbacks.length; i++) {
    if (fallbacks[i] === currentImage) {
      currentIndex = i;
      break;
    }
  }
  
  // If current image is not found or is the last one, return first fallback
  if (currentIndex === -1 || currentIndex >= fallbacks.length - 1) {
    return fallbacks[0] as FallbackImage<T>;
  }
  
  return fallbacks[currentIndex + 1] as FallbackImage<T>;
}

export { getDefaultImage, getNextFallbackImage };

// Note: getNextFallbackImage is not currently used in the codebase,
// so we're removing it to avoid type issues and simplify the code.
// If needed later, we can reimplement it with proper typing. 