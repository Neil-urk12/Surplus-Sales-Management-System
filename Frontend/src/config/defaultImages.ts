export const defaultImages = {
  cab: [
    'https://loremflickr.com/600/400/car',
    'https://placehold.co/600x400/png?text=Car',
    'https://via.placeholder.com/600x400.png?text=Car',
    '/assets/images/default-car.png'
<<<<<<< HEAD
  ],
  accessory: [
    'https://loremflickr.com/600/400/accessory',
    'https://placehold.co/600x400/png?text=Accessory',
    'https://via.placeholder.com/600x400.png?text=Accessory',
    '/assets/images/default-accessory.png'
  ]
} as const;

// More descriptive type names that better explain their purpose
type DefaultImageTypeArray<T extends keyof typeof defaultImages> = (typeof defaultImages)[T];
type DefaultImageType<T extends keyof typeof defaultImages> = DefaultImageTypeArray<T>[number];

/**
 * Returns the first fallback image from the fallback chain for the specified type.
 * This function always returns the first image in the fallback sequence, regardless
 * of any previous fallback attempts.
 * 
 * @param type - The type of image fallback chain to use (e.g., 'cab' or 'accessory')
 * @returns The first fallback image URL in the chain
 */
function getFirstFallbackImage<T extends keyof typeof defaultImages>(type: T): DefaultImageType<T> {
  return defaultImages[type][0];
}

/**
 * Returns the next fallback image in the sequence for the specified type.
 * If the current image is the last in the sequence, it wraps around to the first image.
 * If the current image is not found in the sequence, returns the first fallback image.
 * 
 * @param currentImage - The current image URL in use
 * @param type - The type of image fallback chain to use (e.g., 'cab' or 'accessory')
 * @returns The next fallback image URL in the chain
 * @throws Error if there are no fallback images defined for the given type
 */
function getNextFallbackImage<T extends keyof typeof defaultImages>(currentImage: string, type: T): DefaultImageType<T> {
  const fallbacks = defaultImages[type] as readonly string[];
  
  if (fallbacks.length === 0) {
    throw new Error(`No fallback images defined for type: ${type}`);
  }

  const currentIndex = fallbacks.indexOf(currentImage);
  const nextIndex = (currentIndex + 1) % fallbacks.length;
  return fallbacks[nextIndex] as DefaultImageType<T>;
}

export { getFirstFallbackImage, getNextFallbackImage }; 
=======
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
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
