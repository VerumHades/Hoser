// types/listing.ts
export const ListingAccessModes = { Private: 0, Public: 1 } as const;
export type ListingAccessMode = (typeof ListingAccessModes)[keyof typeof ListingAccessModes];

export interface Money { amount: number; code: string; }

export interface HardwareSpecification { cpu?: number; ramBytes?: number; diskBytes?: number; }