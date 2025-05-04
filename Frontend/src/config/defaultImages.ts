export const defaultImages = {
  cab: [
    'https://loremflickr.com/600/400/car',
    'https://placehold.co/600x400/png?text=Car',
    'https://via.placeholder.com/600x400.png?text=Car',
    '/assets/images/default-car.png'
  ]
} as const;

type ImageType = keyof typeof defaultImages;
type FallbackArray = typeof defaultImages[ImageType];
type FallbackImage = FallbackArray[number];

// Helper function to assert a value is not undefined
function assertDefined<T>(value: T | undefined): asserts value is T {
  if (value === undefined) {
    throw new Error('Expected value to be defined');
  }
}

// Ensure we always have at least one fallback image
const ensureValidFallback = <T extends readonly string[]>(fallbacks: T): NonNullable<T[number]> => {
  if (fallbacks.length === 0) {
    throw new Error('Fallback array cannot be empty');
  }
  const firstImage = fallbacks[0];
  assertDefined(firstImage);
  return firstImage;
};

export const getNextFallbackImage = (currentImage: string, type: ImageType): FallbackImage => {
  const fallbacks = defaultImages[type];
  const currentIndex = fallbacks.indexOf(currentImage as FallbackImage);
  
  // If current image is not found or is the last one, return first fallback
  if (currentIndex === -1 || currentIndex >= fallbacks.length - 1) {
    return ensureValidFallback(fallbacks);
  }
  
  const nextImage = fallbacks[currentIndex + 1];
  assertDefined(nextImage);
  return nextImage;
};

export const getDefaultImage = (type: ImageType): FallbackImage => {
  return ensureValidFallback(defaultImages[type]);
}; 