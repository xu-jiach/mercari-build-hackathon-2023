export const ItemStatus = {
    ItemStatusInitial: 1,
    ItemStatusOnSale: 2,
    ItemStatusSoldOut: 3,
} as const;

export type ItemStatus = (typeof ItemStatus)[keyof typeof ItemStatus];

export interface Item {
    id: number;
    name: string;
    category_id?: number;
    category_name?: string;
    user_id?: number;
    price: number;
    status?: ItemStatus;
    description?: string;
}